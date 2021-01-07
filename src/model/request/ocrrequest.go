/**
 * @Author: lzw5399
 * @Date: 2020/9/30 14:50
 * @Desc:
 */
package request

import "mime/multipart"

type FileFormRequest struct {
	OcrBase
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type FileWithPixelPointRequest struct {
	FileFormRequest
	MatrixPixels []MatrixPixel `form:"-" json:"matrixPixels"` // formdata没法绑定这种对象数组
}

type Base64Request struct {
	OcrBase
	Base64 string `json:"base64" binding:"required"`
}

type Base64WithPixelPointRequest struct {
	OcrBase
	Base64       string        `json:"base64" binding:"required"`
	MatrixPixels []MatrixPixel `json:"matrixPixels" binding:"required"`
}

type OcrBase struct {
	Languages               string `form:"languages" json:"languages"` // 检测语言 eng chi_sim等
	Whitelist               string `form:"whitelist" json:"whitelist"`
	HOCRMode                bool   `form:"hocrMode" json:"hocrMode"`
	TrimLineFeed            bool   `form:"trimLineFeed" json:"trimLineFeed"`
	PreserveInterwordSpaces bool   `form:"preserveInterwordSpaces" json:"preserveInterwordSpaces"`
	SpecialHandling         int    `form:"specialHandling" json:"specialHandling"` // 特殊处理的保留字段: 0.无 | 1.常规通气
}

// 两个像素坐标点能圈出一个矩阵
type MatrixPixel struct {
	PointA Pixel `form:"pointA" json:"pointA"`
	PointB Pixel `form:"pointB" json:"pointB"`
}

// [{ "pointA": {"x": 127, "y": 249}, "pointB": {"x": 983, "y": 309} }]
// 像素坐标点
type Pixel struct {
	X int `form:"x" json:"x"`
	Y int `form:"y" json:"y"`
}
