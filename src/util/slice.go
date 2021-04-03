/**
 * @Author: lzw5399
 * @Date: 2021/3/27 20:45
 * @Desc:
 */
package util

// 取差集
func SliceDiff(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func SliceAnyString(s []string, it string) bool {
	for _, item := range s {
		if item == it {
			return true
		}
	}

	return false
}

func SliceMinMax(array []int) (min int, max int) {
	max = array[0]
	min = array[0]

	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return
}
