/**
 * @Author: lzw5399
 * @Date: 2020/9/30 23:45
 * @Desc: image processing service
 */
package service

import (
	"bytes"
	"image"
	"io"
	"mime/multipart"
	"net/http"

	"bank/distributedquery/src/model/request"

	"github.com/disintegration/imaging"
)

var supportImgType = [4]string{
	"image/png",
	"image/jpeg",
	"image/gif",
	"image/tiff",
}

func EnsureFileType(f multipart.File) (bool, string, error) {
	buff := make([]byte, 512) // docs tell that it take only first 512 bytes into consideration
	if _, err := f.Read(buff); err != nil {
		return false, "", err
	}

	contentType := http.DetectContentType(buff)

	// 把偏移量移回0
	f.Seek(0, 0)

	for _, v := range supportImgType {
		if contentType == v {
			return true, contentType, nil
		}
	}

	return false, contentType, nil
}

// 像素点切割和灰度化, 返回image切片
func CropAndGrayImage(f  io.Reader, re []request.MatrixPixel) ([]image.Image, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		return nil, err
	}

	reader := bytes.NewReader(buf.Bytes())

	img, err := imaging.Decode(reader)
	if err != nil {
		return nil, err
	}

	imgs := make([]image.Image, len(re))

	for i, v := range re {
		var tempImg image.Image = imaging.Crop(img, image.Rect(v.PointA.X, v.PointA.Y, v.PointB.X, v.PointB.Y))
		tempImg = imaging.Grayscale(tempImg)
		imgs[i] = tempImg
	}

	return imgs, nil
}

// 图像灰度化
func GrayImage(f multipart.File) (image.Image, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		return nil, err
	}

	reader := bytes.NewReader(buf.Bytes())

	img, err := imaging.Decode(reader)
	if err != nil {
		return nil, err
	}

	img = imaging.Grayscale(img)
	return img, nil
}
