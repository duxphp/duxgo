package pkg

import (
	"github.com/duxphp/duxgo/core"
	"github.com/duxphp/duxgo/exception"
	"github.com/duxphp/duxgo/pkg/image"
	"github.com/duxphp/duxgo/pkg/upload"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
	"io/ioutil"
	"strings"
	"time"
)

type Upload struct {
}

// NewUpload 存储库对象
func NewUpload() *Upload {
	return &Upload{}
}

// Upload 上传文件
func (s *Upload) Upload(file any, name string) (*upload.File, error) {
	var err error
	var fileByte []byte

	switch file.(type) {
	case string:
		if strings.Index(file.(string), "http://") != -1 || strings.Index(file.(string), "https://") != -1 {
			// 抓取远程文件
			resp, err := resty.New().R().Get(file.(string))
			if err != nil {
				return nil, exception.Internal(err)
			}
			if resp.StatusCode() != 200 {
				return nil, exception.Internal(resp.String())
			}
			fileByte = resp.Body()
		} else {
			// 读取本地文件
			buf, err := ioutil.ReadFile(file.(string))
			if err != nil {
				return nil, exception.Internal(err)
			}
			fileByte = buf
		}
	case []byte:
		fileByte = file.([]byte)
	default:
		return nil, exception.BusinessError("不支持的上传文件")
	}

	maxSize := core.Config["storage"].GetInt("driver.maxSize")

	size := len(fileByte) / 1024
	if size > maxSize {
		return nil, exception.BusinessErrorf("上传文件超出大小:%v", maxSize)
	}

	// 处理图片
	imgResize := core.Config["storage"].GetStringMap("imageResize")
	imgWater := core.Config["storage"].GetStringMap("imageWater")

	if cast.ToBool(imgResize["status"]) || cast.ToBool(imgWater["status"]) {
		img, err := image.New(fileByte)
		if err != nil {
			return nil, err
		}
		if img != nil {
			if cast.ToBool(imgResize["status"]) {
				err := img.Resize(cast.ToInt(imgResize["width"]), cast.ToInt(imgResize["height"]))
				if err != nil {
					return nil, err
				}
			}
			if cast.ToBool(imgWater["status"]) {
				err := img.Watermark(cast.ToString(imgWater["file"]), image.WaterPos(cast.ToInt(imgWater["position"])), cast.ToFloat64(imgWater["opacity"]), cast.ToInt(imgWater["margin"]))
				if err != nil {
					return nil, err
				}
			}
			fileByte, err = img.Save(90)
			if err != nil {
				return nil, err
			}
		}
	}
	var Type = core.Config["storage"].GetString("driver.type")
	core.Logger.Debug().Interface("type", Type).Msg("upload test")
	up, err := upload.New(Type, s.getConfig())
	if err != nil {
		return nil, err
	}

	fileInfo, err := up.Save(fileByte, name, time.Now().Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	return fileInfo, nil
}

// Remove 删除文件
func (s *Upload) Remove(path string, driver ...string) error {
	var err error
	var Type = core.Config["storage"].GetString("driver.type")
	if len(driver) > 0 {
		Type = driver[0]
	}
	up, err := upload.New(Type, s.getConfig())
	if err != nil {
		return err
	}
	err = up.Del(path)
	if err != nil {
		return err
	}
	return nil
}

func (s *Upload) getConfig() any {
	var driverConfig any
	var Type = core.Config["storage"].GetString("driver.type")
	switch Type {
	case "qiniu":
		var qiniuConfig upload.ConfigQiniu
		core.Config["storage"].UnmarshalKey("driver.qiniu", &qiniuConfig)
		driverConfig = qiniuConfig
	case "oss":
		var ossConfig upload.ConfigOss
		core.Config["storage"].UnmarshalKey("driver.oss", &ossConfig)
		driverConfig = ossConfig
	case "local":
		var localConfig upload.ConfigLocal
		core.Config["storage"].UnmarshalKey("driver.local", &localConfig)
		localConfig.UrlPath = core.Config["app"].GetString("app.baseurl") + "/" + localConfig.UrlPath
		driverConfig = localConfig
	}

	return driverConfig
}
