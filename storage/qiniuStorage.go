package storage

import (
	"context"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
	"strings"
)

type QiniuFileStorage struct {
	Bucket string
	Domain string
	Mac    *qbox.Mac
}

func NewQiniuStorage(Bucket, AccessKey, SecretKey, Domain string) *QiniuFileStorage {
	return &QiniuFileStorage{
		Mac:    qbox.NewMac(AccessKey, SecretKey),
		Domain: Domain,
		Bucket: Bucket,
	}
}

func (qfs *QiniuFileStorage) write(ctx context.Context, path string, contents string, config map[string]any) error {
	return qfs.writeStream(ctx, path, strings.NewReader(contents), config)
}

func (qfs *QiniuFileStorage) writeStream(ctx context.Context, path string, stream io.Reader, config map[string]any) error {
	putPolicy := storage.PutPolicy{
		Scope: qfs.Bucket + ":" + path,
	}
	upToken := putPolicy.UploadToken(qfs.Mac)
	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	err := formUploader.Put(ctx, &ret, upToken, path, stream, -1, nil)
	if err != nil {
		return err
	}
	return nil
}

func (qfs *QiniuFileStorage) read(ctx context.Context, path string) (string, error) {
	stream, err := qfs.readStream(ctx, path)
	if err != nil {
		return "", err
	}
	buf := new(strings.Builder)
	_, err = io.Copy(buf, stream)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (qfs *QiniuFileStorage) readStream(ctx context.Context, path string) (io.Reader, error) {
	url, err := qfs.publicUrl(ctx, path)
	if err != nil {
		return nil, err
	}
	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New("failed to read file " + path)
	}
	return resp.RawBody(), nil
}

func (qfs *QiniuFileStorage) delete(ctx context.Context, path string) error {
	cfg := storage.Config{}
	bucketManager := storage.NewBucketManager(qfs.Mac, &cfg)
	return bucketManager.Delete(qfs.Bucket, path)
}

func (qfs *QiniuFileStorage) publicUrl(ctx context.Context, path string) (string, error) {
	return storage.MakePublicURL(qfs.Domain, path), nil
}
