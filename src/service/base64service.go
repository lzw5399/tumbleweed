/**
 * @Author: lzw5399
 * @Date: 2020/10/25 19:51
 * @Desc:
 */
package service

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"bank/distributedquery/src/util"

	"github.com/satori/go.uuid"
)

var supportBase64TypeMap = map[string]string{
	"image/png":       "data:image/png;base64,",
	"image/jpeg":      "data:image/jpeg;base64,",
	"image/gif":       "data:image/gif;base64,",
	"image/tiff":      "data:image/tiff;base64,",
	"application/pdf": "data:application/pdf;base64,",
}

func EnsureContentType(str string) (base64 string, isPdf bool, contentType string, err error) {
	for k, v := range supportBase64TypeMap {
		if strings.HasPrefix(str, v) {
			base64 = str[len(v):]
			contentType = k
			if k == "application/pdf" {
				isPdf = true
			}
			return
		}
	}
	err = errors.New("invalid or unsupported content type")
	return
}

// 处理pdf
func PdfToImgsThenGetBytes(base64 string) ([][]byte, error) {
	// pdf先保存到本地
	filePath, err := save(base64)
	if err != nil {
		return nil, err
	}
	defer func() {
		os.Remove(filePath)
	}()

	// pdf分页转成png
	dirToSave, _ := os.Getwd()
	imgs, err := util.PdfToImgs(filePath, dirToSave)
	if err != nil {
		return nil, err
	}
	defer func() {
		for _, path := range imgs {
			os.Remove(path)
		}
	}()

	// 读取png成[]byte
	var finalArray [][]byte
	for _, imgPath := range imgs {
		byteArray, err := ioutil.ReadFile(imgPath)
		if err != nil {
			return nil, err
		}
		finalArray = append(finalArray, byteArray)
	}

	return finalArray, nil
}

func save(base64Str string) (filePath string, err error) {
	buf, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return
	}

	filePath = uuid.NewV4().String() + ".pdf"
	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	if _, err = file.Write(buf); err != nil {
		return
	}

	err = file.Sync()
	return
}
