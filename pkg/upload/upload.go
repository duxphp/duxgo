package upload

import (
	"bytes"
	"github.com/duxphp/duxgo/exception"
	"github.com/duxphp/duxgo/util/function"
	"github.com/h2non/filetype"
	"github.com/taokunTeam/go-storage/kodo"
	"github.com/taokunTeam/go-storage/local"
	"github.com/taokunTeam/go-storage/storage"
	"path/filepath"
	"strings"
)

type Upload struct {
	File   *[]byte
	Driver string
	Store  storage.Storage
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
	url := ""
	var store storage.Storage
	var err error
	switch driver {
	case "local":
		config := driverConfig.(ConfigLocal)
		abs, err := filepath.Abs("./" + config.Path)
		if err != nil {
			return nil, exception.Internal(err)
		}

		store, err = local.Init(local.Config{
			RootDir: abs,
			AppUrl:  config.UrlPath,
		})
		if err != nil {
			return nil, err
		}
		url = config.UrlPath
	case "qiniu":
		config := driverConfig.(ConfigQiniu)

		store, err = kodo.Init(kodo.Config{
			AccessKey: config.AccountName,
			Bucket:    config.Bucket,
			Domain:    config.Domain,
			SecretKey: config.AccountKey,
		})
		if err != nil {
			return nil, err
		}
		//driverStr = fmt.Sprintf("kodo://%v/uploads/?credential=hmac:%v:%v&endpoint=%v", config.Bucket, config.AccountName, config.AccountKey, config.Domain)
		url = config.Domain
	}

	//store, err := services.NewStoragerFromString(driverStr)
	//if err != nil {
	//	return nil, exception.Internal(err)
	//}

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
	ext = strings.Trim(ext, ".")
	reader := bytes.NewReader(file)

	if reader.Len() == 0 {
		return nil, exception.BusinessError("上传格式错误")
	}

	filename := "uploads/" + dir + "/" + function.Md5(string(file)) + "." + ext
	err := t.Store.Put(filename, reader, int64(length), kind.MIME.Value)
	if err != nil {
		return nil, exception.Internal(err)
	}
	url := t.Url + "/" + filename

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
