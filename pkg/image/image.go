package image

import (
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/duxphp/duxgo/v2/util/function"
	"github.com/h2non/filetype"
	"github.com/samber/lo"
	"image"
)

type Image struct {
	Status    bool
	Ext       string
	Size      int
	ImgBuffer image.Image
}

// New 图片处理
func New(file []byte) (*Image, error) {
	kind, _ := filetype.Match(file)
	ext := kind.Extension
	// 过滤图片格式
	status := false
	imageTypes := []string{"jpg", "jpeg", "png", "gif", "tif", "tiff", "bmp"}
	_, ok := lo.Find[string](imageTypes, func(i string) bool {
		return i == ext
	})
	if ok {
		status = true
	}
	if !status {
		return nil, nil
	}
	// 初始化对象
	reader := bytes.NewReader(file)
	imgBuffer, err := imaging.Decode(reader)
	if err != nil {
		return nil, err
	}
	return &Image{
		Ext:       ext,
		Size:      len(file),
		Status:    status,
		ImgBuffer: imgBuffer,
	}, nil
}

// Resize 图片缩放
func (t *Image) Resize(width int, height int) error {
	if !t.Status {
		return nil
	}
	t.ImgBuffer = imaging.Resize(t.ImgBuffer, width, 0, imaging.Lanczos)
	t.ImgBuffer = imaging.Resize(t.ImgBuffer, 0, height, imaging.Lanczos)
	return nil
}

// WaterPos 水印未知
type WaterPos int

const (
	PosTop WaterPos = iota
	PostTopLeft
	PostTopRight
	PosLeft
	PosCenter
	PosRight
	PosBottom
	PosBottomLeft
	PosBottomRight
)

// Watermark 图片水印
func (t *Image) Watermark(file string, pos WaterPos, quality float64, imgMargin int) error {
	if !t.Status {
		return nil
	}
	if !function.IsExist(file) {
		return nil
	}
	// 载入水印图片
	waterBuffer, err := imaging.Open(file)
	if err != nil {
		return err
	}
	// 获取图片信息
	imgWidth := t.ImgBuffer.Bounds().Dx()
	imgHeight := t.ImgBuffer.Bounds().Dy()
	waterWidth := waterBuffer.Bounds().Dx()
	waterHeight := waterBuffer.Bounds().Dy()
	// 水印图片大于原始图片不处理水印
	margin := imgMargin + 50
	if imgWidth <= waterWidth+margin || imgHeight <= waterHeight+margin {
		return nil
	}

	left := 0
	top := 0
	iw := imgWidth / 2
	ww := waterWidth / 2
	ih := imgHeight / 2
	wh := waterHeight / 2
	switch pos {
	case 0:
		left = iw - ww
		top = margin
	case 1:
		top = margin
		left = margin
	case 2:
		top = margin
		left = imgWidth - waterWidth - margin
	case 3:
		top = ih - wh
		left = margin
	case 4:
		top = ih - wh
		left = iw - ww
	case 5:
		top = ih - wh
		left = imgWidth - waterWidth - margin
	case 6:
		top = imgHeight - waterHeight - margin
		left = iw - ww
	case 7:
		top = imgHeight - waterHeight - margin
		left = margin
	case 8:
		top = imgHeight - waterHeight - margin
		left = imgWidth - waterWidth - margin
	}
	t.ImgBuffer = imaging.Overlay(t.ImgBuffer, waterBuffer, image.Pt(left, top), quality)
	return nil
}

// Save 保存图片
func (t *Image) Save(quality int) ([]byte, error) {
	if !t.Status {
		return nil, nil
	}
	f, err := imaging.FormatFromFilename("dux." + t.Ext)
	if err != nil {
		return nil, err
	}
	reader := new(bytes.Buffer)
	err = imaging.Encode(reader, t.ImgBuffer, f, imaging.JPEGQuality(quality))
	if err != nil {
		return nil, err
	}
	return reader.Bytes(), nil
}
