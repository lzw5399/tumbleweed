/**
 * @Author: lzw5399
 * @Date: 2020/9/30 13:26
 * @Desc: format response
 */
package response

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpResponse struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

func Ok(c *gin.Context) {
	result(c, nil, "success")
}

func OkWithMessage(c *gin.Context, message string) {
	result(c, nil, message)
}

func OkWithData(c *gin.Context, data interface{}) {
	result(c, data, "操作成功")
}

func OkWithDetailed(c *gin.Context, data interface{}, message string) {
	result(c, data, message)
}

// 最终返回c.PureJson
func OkWithPureData(c *gin.Context, data interface{}) {
	c.PureJSON(http.StatusOK, HttpResponse{
		true,
		"操作成功",
		data,
	})
}

func Failed(c *gin.Context, status int) {
	FailWithMsg(c, status, "操作失败")
}

func FailWithMsg(c *gin.Context, status int, err interface{}) {
	var msg interface{}

	switch err.(type) {
	case error:
		msg = err.(error).Error()

	case string:
		msg = err.(string)

	default:
		msg = ""
	}

	switch status {
	// 400
	case http.StatusBadRequest:

	// 401
	case http.StatusUnauthorized:

	// 404
	case http.StatusNotFound:
		msg = "记录未找到"

	// 500
	case http.StatusInternalServerError:
		//msg = "server internal error, please contact the maintainer"
		log.Printf("err: %s", err)
	}

	resultWithStatus(c, status, false, nil, msg)
	c.Abort()
}

func result(c *gin.Context, data interface{}, msg string) {
	resultWithStatus(c, http.StatusOK, true, data, msg)
}

func resultWithStatus(c *gin.Context, statusCode int, success bool, data interface{}, msg interface{}) {
	c.JSON(statusCode, HttpResponse{
		success,
		msg,
		data,
	})
}
