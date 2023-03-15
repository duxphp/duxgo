package storage

import (
	"context"
	"github.com/duxphp/duxgo/v2/config"
	"io"
)

type Storage struct {
	driver FileStorage
}

type FileStorage interface {
	write(ctx context.Context, path string, contents string, config map[string]any) error
	writeStream(ctx context.Context, path string, stream io.Reader, config map[string]any) error
	read(ctx context.Context, path string) (string, error)
	readStream(ctx context.Context, path string) (io.Reader, error)
	delete(ctx context.Context, path string) error
	publicUrl(ctx context.Context, path string) (string, error)
	privateUrl(ctx context.Context, path string) (string, error)
}

func New(types ...string) FileStorage {
	var Type string
	if len(types) <= 0 {
		Type = config.Get("storage").GetString("driver.type")
	} else {
		Type = types[0]
	}
	storeConfig := config.Get("storage").GetStringMapString("driver." + Type)

	var driver FileStorage
	switch Type {
	case "local":
		driver = NewLocalStorage(storeConfig["path"], storeConfig["domain"])
		break
	case "qiniu":
		driver = NewQiniuStorage(storeConfig["bucket"], storeConfig["accessKey"], storeConfig["secretKey"], storeConfig["domain"])
		break
	case "cos":
		driver = NewCoStorage(storeConfig["secretId"], storeConfig["secretKey"], storeConfig["region"], storeConfig["bucket"], storeConfig["domain"])
		break
	case "oss":
		driver = NewOssStorage(storeConfig["accessId"], storeConfig["accessSecret"], storeConfig["endpoint"], storeConfig["bucketName"], storeConfig["domain"])
		break
	}
	return driver
}
