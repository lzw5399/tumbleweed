/**
 * @Author: lzw5399
 * @Date: 2021/1/16 11:21
 * @Desc:
 */
package constant

// event的类别常量
const (
	StartEvent = iota + 1
	EndEvent
)

// process instance 的type类别
const (
	MyToDo   = iota + 1 // 我的待办
	ICreated            // 我创建的
	IRelated            // 和我相关的
	All                 // 所有
)
