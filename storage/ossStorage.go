package storage

import (
	"context"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"net/url"
	"strings"
)

type OssStorage struct {
	Client     *oss.Client
	Domain     string
	BucketName string
}

func NewOssStorage(accessKeyId, accessKeySecret, endpoint, bucketName, domain string) *OssStorage {
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return nil
	}
	return &OssStorage{
		Client:     client,
		Domain:     domain,
		BucketName: bucketName,
	}
}

func (ofs *OssStorage) write(ctx context.Context, path string, contents string, config map[string]interface{}) error {
	return ofs.writeStream(ctx, path, strings.NewReader(contents), config)
}

func (ofs *OssStorage) writeStream(ctx context.Context, path string, stream io.Reader, config map[string]interface{}) error {
	bucket, err := ofs.Client.Bucket(ofs.BucketName)
	if err != nil {
		return err
	}
	options := []oss.Option{}
	if val, ok := config["Content-Type"]; ok {
		options = append(options, oss.ContentType(val.(string)))
	}
	err = bucket.PutObject(path, stream, options...)
	if err != nil {
		return err
	}
	return nil
}

func (ofs *OssStorage) read(ctx context.Context, path string) (string, error) {
	stream, err := ofs.readStream(ctx, path)
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

func (ofs *OssStorage) readStream(ctx context.Context, path string) (io.Reader, error) {
	bucket, err := ofs.Client.Bucket(ofs.BucketName)
	if err != nil {
		return nil, err
	}
	body, err := bucket.GetObject(path)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (ofs *OssStorage) delete(ctx context.Context, path string) error {
	bucket, err := ofs.Client.Bucket(ofs.BucketName)
	if err != nil {
		return err
	}
	err = bucket.DeleteObject(path)
	if err != nil {
		return err
	}
	return nil
}

func (ofs *OssStorage) publicUrl(ctx context.Context, path string) (string, error) {
	srcUrl := fmt.Sprintf("%s/%s", strings.TrimRight(ofs.Domain, "/"), path)
	srcUri, _ := url.Parse(srcUrl)
	finalUrl := srcUri.String()
	return finalUrl, nil
}

func (ofs *OssStorage) privateUrl(ctx context.Context, path string) (string, error) {
	bucket, err := ofs.Client.Bucket(ofs.BucketName)
	if err != nil {
		return "", err
	}
	finalUrl, err := bucket.SignURL(path, oss.HTTPGet, 3600)
	if err != nil {
		return "", err
	}
	return finalUrl, nil
}
