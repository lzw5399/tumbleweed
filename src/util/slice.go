/**
 * @Author: lzw5399
 * @Date: 2021/3/27 20:45
 * @Desc:
 */
package util

// 取差集
func SliceDiff(a, b []int) (diff []int) {
	m := make(map[int]bool)

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
