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
