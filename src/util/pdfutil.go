/**
 * @Author: lzw5399
 * @Date: 2020/10/25 22:01
 * @Desc:
 */
package util

func PdfToImgs(filePath, dirToSave string) (imgs []string, err error) {
	//if !PathExists(filePath) {
	//	global.BANK_LOGGER.Error("不存在filePath")
	//	err = errors.New("filePath doesn't exist")
	//	return
	//}
	//
	//if !PathExists(dirToSave) {
	//	global.BANK_LOGGER.Error("不存在dirToSave")
	//	err = errors.New("dirToSave doesn't exist")
	//	return
	//}
	//
	//doc, err := fitz.New(filePath)
	//if err != nil {
	//	return
	//}
	//defer doc.Close()
	//
	//// 获取当前目录
	//dir, err := os.Getwd()
	//if err != nil {
	//	return
	//}
	//
	//// 从xxx.pdf获取xxx
	//fileName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	//
	//for n := 0; n < doc.NumPage(); n++ {
	//	img, e := doc.Image(n)
	//	if e != nil {
	//		err = e
	//		return
	//	}
	//
	//	path := filepath.Join(dir, fmt.Sprintf("%s-%03d.png", fileName, n))
	//	f, e := os.Create(path)
	//	if e != nil {
	//		err = e
	//		return
	//	}
	//
	//	e = png.Encode(f, img)
	//	if e != nil {
	//		err = e
	//		return
	//	}
	//	f.Close()
	//	imgs = append(imgs, path)
	//}

	return
}
