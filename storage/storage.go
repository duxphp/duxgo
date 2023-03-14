package storage

import (
	"context"
	"github.com/duxphp/duxgo/v2/global"
	"io"
)

type Storage struct {
	driver FileStorage
}

// FileStorage 存储接口
type FileStorage interface {
	// 写入字符串到文件
	write(ctx context.Context, path string, contents string, config map[string]any) error
	// 写入文件流到文件
	writeStream(ctx context.Context, path string, stream io.Reader, config map[string]any) error
	// 读取文件到字符串
	read(ctx context.Context, path string) (string, error)
	// 读取文件到文件流
	readStream(ctx context.Context, path string) (io.Reader, error)
	// 删除文件
	delete(ctx context.Context, path string) error
	// 获取文件公开url
	publicUrl(ctx context.Context, path string) (string, error)
	// 获取文件私有签名url
	privateUrl(ctx context.Context, path string) (string, error)
}

// New 存储库对象
func New(types ...string) FileStorage {
	var Type string
	if len(types) <= 0 {
		Type = global.Config["storage"].GetString("driver.type")
	} else {
		Type = types[0]
	}
	config := global.Config["storage"].GetStringMapString("driver." + Type)

	var driver FileStorage
	switch Type {
	case "local":
		driver = NewLocalStorage(config["path"], config["domain"])
		break
	case "qiniu":
		driver = NewQiniuStorage(config["bucket"], config["accessKey"], config["secretKey"], config["domain"])
		break
	case "cos":
		driver = NewCoStorage(config["secretId"], config["secretKey"], config["region"], config["bucket"], config["domain"])
		break
	case "oss":
		driver = NewOssStorage(config["accessId"], config["accessSecret"], config["endpoint"], config["bucketName"], config["domain"])
		break
	}
	return driver
}
