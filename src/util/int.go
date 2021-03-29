/**
 * @Author: lzw5399
 * @Date: 2021/3/9 11:50
 * @Desc:
 */
package util

import "strconv"

func StringToUint(str string) uint {
	return uint(StringToInt(str))
}

func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return i
}

func InterfaceToUint(i interface{}) uint {
	str, succeed := i.(string)
	if !succeed {
		return 0
	}

	return StringToUint(str)
}

func ParseToInt64Array(arr []int) []int64 {
	intArr := make([]int64, len(arr))
	for index, item := range arr {
		intArr[index] = int64(item)
	}

	return intArr
}

func IsInteger(f float64) bool {
	f2 := float64(int64(f))

	return f == f2
}

func ParseToIntArray(arr []interface{}) []int {
	intArr := make([]int, len(arr))
	for index, item := range arr {
		intArr[index] = int(item.(float64))
	}

	return intArr
}
