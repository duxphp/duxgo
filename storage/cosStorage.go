package storage

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type CoStorage struct {
	Client     *cos.Client
	BucketName string
	Domain     string
}

func NewCoStorage(secretId, secretKey, region, bucketName, domain string) *CoStorage {
	u, _ := url.Parse("https://" + bucketName + ".cos." + region + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	})
	return &CoStorage{
		Client:     c,
		BucketName: bucketName,
		Domain:     domain,
	}
}

func (tfs *CoStorage) write(ctx context.Context, path string, contents string, config map[string]any) error {
	return tfs.writeStream(ctx, path, strings.NewReader(contents), config)
}

func (tfs *CoStorage) writeStream(ctx context.Context, path string, stream io.Reader, config map[string]any) error {
	_, err := tfs.Client.Object.Put(ctx, path, stream, nil)
	if err != nil {
		return err
	}
	return nil
}

func (tfs *CoStorage) read(ctx context.Context, path string) (string, error) {
	stream, err := tfs.readStream(ctx, path)
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

func (tfs *CoStorage) readStream(ctx context.Context, path string) (io.Reader, error) {
	resp, err := tfs.Client.Object.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (tfs *CoStorage) delete(ctx context.Context, path string) error {
	_, err := tfs.Client.Object.Delete(ctx, path)
	if err != nil {
		return err
	}
	return nil
}

func (tfs *CoStorage) publicUrl(ctx context.Context, path string) (string, error) {
	srcUrl := fmt.Sprintf("%s/%s", strings.TrimRight(tfs.Domain, "/"), path)
	srcUri, _ := url.Parse(srcUrl)
	finalUrl := srcUri.String()
	return finalUrl, nil
}

func (tfs *CoStorage) privateUrl(ctx context.Context, path string) (string, error) {
	u := tfs.Client.Object.GetObjectURL(path)
	return u.String(), nil
}
