/**
 * @Author: lzw5399
 * @Date: 2020/11/13 17:19
 * @Desc: 针对OcrBase.SpecialHandling进行特殊处理的
 */
package service

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"strconv"
	"strings"

	"bank/distributedquery/src/model/request"

	. "github.com/ahmetb/go-linq"
)

const (
	NULL_COLUMN_PLACEHOLDER = "{NULL}"
	SEPARATOR               = " "
)

func SpecialHandling(ocrBase request.OcrBase, str string) string {
	// HOCRMode不支持特殊处理
	if ocrBase.HOCRMode {
		return str
	}

	switch ocrBase.SpecialHandling {
	// 常规通气
	case 1:
		return conventionalVentilation(str)
	default:
		return str
	}
}

func SpecialHandlingArray(ocrBase request.OcrBase, texts []string) []string {
	finalTexts := make([]string, len(texts))
	for _, v := range texts {
		finalTexts = append(finalTexts, SpecialHandling(ocrBase, v))
	}

	return finalTexts
}

func SpecialHandlingInterface(ocrBase request.OcrBase, texts interface{}) interface{} {
	switch texts.(type) {
	case string:
		return SpecialHandling(ocrBase, texts.(string))
	case []string:
		return SpecialHandlingArray(ocrBase, texts.([]string))
	default:
		return texts
	}
}

// 获取字符串的行数
func getLineCount(str string) (int, error) {
	r := strings.NewReader(str)
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// 针对[常规通气检查报告单]的特殊处理
func conventionalVentilation(str string) string {
	lineCount, _ := getLineCount(str)
	if lineCount <= 2 {
		return str
	}

	log.Printf("开始针对[常规通气检查报告单]的特殊处理，当前总行数为:%d, 处理前的内容为:\n%s\n", lineCount, str)
	finalStr := ""
	reader := bufio.NewReader(strings.NewReader(str))

	currentLine := 0
	for {
		lineBytes, _, _ := reader.ReadLine()
		if lineBytes == nil {
			break
		}
		line := string(lineBytes)
		currentLine++

		fields := strings.Fields(line)
		prefix := defaultPrefix(currentLine)
		numbers := defaultNumbers(line, 0)

		// 4 5 6行 VT IC ERV
		if From([]int{4, 5, 6}).AnyWith(func(i interface{}) bool {
			return i == currentLine
		}) {
			if len(numbers) == 3 {
				line = genLineStr(prefix, genJoinedStr(numbers), genPlaceHolderStr(3))
			} else {
				// 最后一个字段在最后一列
				if judgeVCMAXLayout(line, numbers) {
					line = genLineStr(prefix, genJoinedStr(numbers[:len(numbers)-1]), genPlaceHolderStr(2), numbers[len(numbers)-1])
				} else {
					line = genLineStr(prefix, genJoinedStr(numbers), genPlaceHolderStr(2))
				}
			}
		}

		// 7 PIF
		if currentLine == 7 {
			if len(numbers) >= 4 {
				line = genLineStr(prefix, genPlaceHolderStr(1), numbers[0], genPlaceHolderStr(1), genJoinedStr(numbers[1:]))
			}
		}

		// 8 9 VC MAX & MV
		if currentLine == 8 || currentLine == 9 {
			log.Println(numbers, judgeVCMAXLayout(line, numbers))
			if len(numbers) == 3 {
				line = genLineStr(prefix, genJoinedStr(numbers), genPlaceHolderStr(3))
			} else {
				// 最后一个字段在最后一列
				if judgeVCMAXLayout(line, numbers) {
					line = genLineStr(prefix, genJoinedStr(numbers[:len(numbers)-1]), genPlaceHolderStr(2), numbers[len(numbers)-1])
				} else {
					line = genLineStr(prefix, genJoinedStr(numbers), genPlaceHolderStr(2))
				}
			}
		}

		// 10 11 13 FVC & FEV 1 & PEF
		if From([]int{10, 11, 13}).AnyWith(func(i interface{}) bool {
			return i == currentLine
		}) {
			line = genLineStr(prefix, genJoinedStr(numbers))
		}

		// 12 FEV 1 % VC MAX
		if currentLine == 12 {
			numbers = defaultNumbers(line, 26)
			line = genLineStr(prefix, genJoinedStr(numbers))
		}

		// 14 MVV
		if currentLine == 14 {
			numbers = defaultNumbers(line, 23)
			line = genLineStr(prefix, genJoinedStr(numbers), genPlaceHolderStr(5))
		}

		// 15 TIN/TEX
		if currentLine == 15 {
			if len(numbers) == 1 {
				line = genLineStr(prefix, genPlaceHolderStr(1), numbers[0], genPlaceHolderStr(4))
			} else if len(numbers) == 2 {
				// 最后一个字段在最后一列
				if judgeTINTEXLayout(line, numbers) {
					line = genLineStr(prefix, genPlaceHolderStr(1), numbers[0], genPlaceHolderStr(3), numbers[1])
				} else {
					line = genLineStr(prefix, genPlaceHolderStr(1), numbers[0], genPlaceHolderStr(1), numbers[1], genPlaceHolderStr(2))
				}
			} else {
				line = genLineStr(prefix, genPlaceHolderStr(1), numbers[0], genPlaceHolderStr(4))
			}
		}

		// 16 17 MIF & MEF
		if currentLine == 16 || currentLine == 17 {
			if len(fields) == 4 {
				// 最后一个字段在最后一列
				if judgeTINTEXLayout(line, fields) {
					line = genLineStr(prefix, genPlaceHolderStr(1), numbers[0], genPlaceHolderStr(3), numbers[1])
				} else {
					line = genLineStr(prefix, genPlaceHolderStr(1), numbers[0], genPlaceHolderStr(1), numbers[1], genPlaceHolderStr(2))
				}
			} else {
				line = genLineStr(prefix, genPlaceHolderStr(1), numbers[0], genPlaceHolderStr(4))
			}
		}

		// 18 19 20 FEF 25 & FEF 50 & FEF 75
		if From([]int{18, 19, 20}).AnyWith(func(i interface{}) bool {
			return i == currentLine
		}) {
			line = genLineStr(prefix, genJoinedStr(numbers))
		}

		// 21 FEF 75/85
		if currentLine == 21 {
			line = genLineStr(prefix, genPlaceHolderStr(1), numbers[0], genPlaceHolderStr(1), genJoinedStr(numbers[1:]))
		}

		// 22 MMEF 75/25
		if currentLine == 22 {
			line = genLineStr(prefix, genJoinedStr(numbers))
		}

		// 23 24 最后两行
		if currentLine == 23 || currentLine == 24 {
			numbers = defaultNumbers(line, 30)
			line = genLineStr(prefix, genPlaceHolderStr(1), numbers[0], genPlaceHolderStr(1), genJoinedStr(numbers[1:]))
		}

		finalStr += line + "\n"
	}

	log.Printf("[常规通气检查报告单]的特殊处理结束，当前总行数为:%d, 处理后的内容为:\n%s", lineCount, finalStr)

	return finalStr
}

func genPlaceHolderStr(count int) string {
	arr := make([]string, count)
	for i := 0; i < count; i++ {
		arr[i] = NULL_COLUMN_PLACEHOLDER
	}

	return strings.Join(arr, SEPARATOR)
}

func genJoinedStr(slice []string) string {
	return strings.Join(slice, SEPARATOR)
}

var prefixMap = map[int]string{
	4:  "VT[L]",
	5:  "IC[L]",
	6:  "ERV[L]",
	7:  "PIF[L/s]",
	8:  "VCMAX[L]",
	9:  "MV[L/min]",
	10: "FVC[L]",
	11: "FEV1[L]",
	12: "FEV1%VCMAX[%]",
	13: "PEF[L/s]",
	14: "MVV[L/min]",
	15: "TIN/TEX",
	16: "MIF[L/s]",
	17: "MEF[L/s]",
	18: "FEF25[L/s]",
	19: "FEF50[L/s]",
	20: "FEF75[L/s]",
	21: "FEF75/85[L/s]",
	22: "MMEF75/25[L/s]",
	23: "Vbackextrapolationex[L]",
	24: "Vbackextrapol.%FVC[%]",
}

func defaultPrefix(currentLine int) string {
	return prefixMap[currentLine]
}

func defaultNumbers(str string, location int) []string {
	// 正常的分割为24，有些只需要22个空格
	const PLAN_A_SPACE int = 24
	const PLAN_B_SPACE int = 22

	if location == 0 {
		location = PLAN_A_SPACE
	}

	if len(str) <= 24 {
		return strings.Fields(str)
	}

	fields := strings.Fields(strings.TrimLeft(str[location:], "]"))

	if len(fields) == 0 {
		return fields
	}

	firstNum, _ := strconv.Atoi(fields[0])

	// 1. 以.开头 2. 超过15 说明截过了，plan b
	if (strings.HasPrefix(fields[0], ".") || firstNum > 15) && location != PLAN_B_SPACE {
		fields = defaultNumbers(str, PLAN_B_SPACE)
	}

	return fields
}

func genLineStr(str ...string) string {
	return strings.Join(str, SEPARATOR)
}

// 针对 4 5 6 VT Ic ERV 8 9 VC MAX & MV
func judgeVCMAXLayout(line string, fields []string) bool {
	lastOne := false

	defer func() {
		if err := recover(); err != nil {
			lastOne = false
		}
	}()

	log.Println("judgeVCMAXLayout")
	log.Println(line)
	log.Println(fields)
	log.Println(fields[len(fields)-2])
	log.Println(fields[len(fields)-1])

	str := getBetweenStr(line, fields[len(fields)-2], fields[len(fields)-1])
	lastOne = len(str) > 12

	return lastOne
}

// 针对15 TIN/TEX 16 17 MIF & MEF
func judgeTINTEXLayout(line string, fields []string) bool {
	lastOne := false

	defer func() {
		if err := recover(); err != nil {
			lastOne = false
		}
	}()

	str := getBetweenStr(line, fields[len(fields)-2], fields[len(fields)-1])
	lastOne = len(str) > 18

	return lastOne
}

func getBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.LastIndex(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}
