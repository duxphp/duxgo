package upload

import (
	"bytes"
	"fmt"
	"github.com/duxphp/duxgo/core"
	"github.com/duxphp/duxgo/exception"
	"github.com/duxphp/duxgo/util/function"
	"github.com/h2non/filetype"
	_ "go.beyondstorage.io/services/fs/v4"
	_ "go.beyondstorage.io/services/kodo/v3"
	_ "go.beyondstorage.io/services/oss/v3"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
	"path/filepath"
	"strings"
)

type Upload struct {
	File   *[]byte
	Driver string
	Store  types.Storager
	Url    string
}

type ConfigLocal struct {
	Path    string
	UrlPath string
}

type ConfigQiniu struct {
	AccountName string
	AccountKey  string
	Domain      string
	Bucket      string
}

type ConfigOss struct {
	AccountName string
	AccountKey  string
	Location    string
	Bucket      string
	Domain      string
}

// New 上传对象
func New(driver string, driverConfig any) (*Upload, error) {
	driverStr := ""
	url := ""
	switch driver {
	case "local":
		config := driverConfig.(ConfigLocal)
		abs, err := filepath.Abs("./" + config.Path)
		if err != nil {
			return nil, exception.Internal(err)
		}
		driverStr = fmt.Sprintf("fs://%v", abs+"/")
		url = config.UrlPath
	case "qiniu":
		config := driverConfig.(ConfigQiniu)
		driverStr = fmt.Sprintf("kodo://%v/uploads/?credential=hmac:%v:%v&endpoint=%v", config.Bucket, config.AccountName, config.AccountKey, config.Domain)
		url = config.Domain + "/uploads"
	case "oss":
		config := driverConfig.(ConfigOss)
		driverStr = fmt.Sprintf("oss://%v/uploads/?credential=hmac:%v:%v&endpoint=https:%v.aliyuncs.com", config.Bucket, config.AccountName, config.AccountKey, config.Location)
		url = config.Domain + "/uploads"
	}

	store, err := services.NewStoragerFromString(driverStr)
	if err != nil {
		return nil, exception.Internal(err)
	}

	return &Upload{
		Driver: driver,
		Store:  store,
		Url:    url,
	}, nil
}

type File struct {
	Url      string
	Filename string
	Path     string
	Size     int
	Ext      string
	Driver   string
}

// Save 保存文件
func (t *Upload) Save(file []byte, name string, dir string) (*File, error) {
	length := len(file)
	kind, _ := filetype.Match(file)
	realExt := kind.Extension
	ext := filepath.Ext(name)
	if ext == "" {
		ext = realExt
	}
	core.Logger.Debug().Msg("upload save")
	ext = strings.Trim(ext, ".")
	reader := bytes.NewReader(file)
	filename := dir + "/" + function.Md5(string(file)) + "." + ext
	core.Logger.Debug().Interface("filename", filename).Interface("reader", reader).Msg("upload save2")
	_, err := t.Store.Write(filename, reader, int64(length))
	if err != nil {
		return nil, exception.Internal(err)
	}
	url := t.Url + "/" + filename
	core.Logger.Debug().Msg("upload save3")

	return &File{
		Url:      url,
		Filename: name,
		Path:     filename,
		Size:     length,
		Ext:      ext,
		Driver:   t.Driver,
	}, nil
}

// Del 删除文件
func (t *Upload) Del(path string) error {
	err := t.Store.Delete(path)
	if err != nil {
		return exception.Internal(err)
	}
	return nil
}
