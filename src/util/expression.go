/**
 * @Author: lzw5399
 * @Date: 2021/3/21 22:02
 * @Desc:
 */
package util

import (
	"fmt"

	"github.com/antonmedv/expr"
)

func CalculateExpression(expression string, env map[string]interface{}) (result bool, err error) {
	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		err = fmt.Errorf("处理失败, 请检查表达式和变量")
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		err = fmt.Errorf("处理失败, 请检查表达式和变量")
		return
	}

	if v, succeed := output.(bool); succeed {
		result = v
		return
	}
	
	err = fmt.Errorf("处理失败, 请检查表达式和变量")
	return
}
