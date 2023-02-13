package storage

import (
	"context"
	"io"
)

type FileInfo struct {
}

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
	// 获取文件url
	publicUrl(ctx context.Context, path string) (string, error)
}

//
//type Storage struct {
//}
//
//// New 存储库对象
//func New() *Storage {
//	var Type = registry.Config["storage"].GetString("driver.type")
//	var config Config
//	registry.Config["storage"].UnmarshalKey("driver."+Type, &config)
//	return &Storage{
//		config: config,
//	}
//}
