package service

import (
	"encoding/json"
	"github.com/duxphp/duxgo/core/exception"
	"github.com/duxphp/duxgo/core/response"
	"duxgopkg"
	"io/ioutil"
)

func Upload(ctx echo.Context, extend ...map[string]any) error {
	// 获取表单文件
	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}
	// 读取上传文件
	openFile, err := file.Open()
	if err != nil {
		return exception.Internal(err)
	}
	fileByte, err := ioutil.ReadAll(openFile)
	if err != nil {
		return exception.Internal(err)
	}
	data, err := pkg.NewUpload().Upload(fileByte, file.Filename)
	if err != nil {
		return err
	}

	var extData []byte
	if len(extend) > 0 {
		extData, _ = json.Marshal(extend[0])
	}

	// event 上传前接口

	//core.Db.Create(&model.ToolFile{
	//	DirId:   dirId,
	//	HasType: HasType,
	//	Driver:  data.Driver,
	//	Url:     data.Url,
	//	Path:    data.Path,
	//	Title:   data.Filename,
	//	Ext:     data.Ext,
	//	Size:    cast.ToInt(data.Size),
	//	Extend:  extData,
	//})

	var files []map[string]any
	files = append(files, map[string]any{
		"url":   data.Url,
		"ext":   data.Ext,
		"title": data.Filename,
		"size":  data.Size,
	})
	return response.New(ctx).Send("ok", files)
}
