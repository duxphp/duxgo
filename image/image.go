package image

import (
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/duxphp/duxgo/v2/helper"
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

// New  image processing
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

// Resize image resizing
func (t *Image) Resize(width int, height int) error {
	if !t.Status {
		return nil
	}
	t.ImgBuffer = imaging.Resize(t.ImgBuffer, width, 0, imaging.Lanczos)
	t.ImgBuffer = imaging.Resize(t.ImgBuffer, 0, height, imaging.Lanczos)
	return nil
}

// WaterPos Watermark position
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

// Watermark image watermarking
func (t *Image) Watermark(file string, pos WaterPos, quality float64, imgMargin int) error {
	if !t.Status {
		return nil
	}
	if !helper.IsExist(file) {
		return nil
	}

	waterBuffer, err := imaging.Open(file)
	if err != nil {
		return err
	}

	imgWidth := t.ImgBuffer.Bounds().Dx()
	imgHeight := t.ImgBuffer.Bounds().Dy()
	waterWidth := waterBuffer.Bounds().Dx()
	waterHeight := waterBuffer.Bounds().Dy()
	// If the watermark image is larger than the original image, the watermark will not be removed.
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

// Save image
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
