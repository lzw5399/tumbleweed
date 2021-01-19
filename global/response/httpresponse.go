/**
 * @Author: lzw5399
 * @Date: 2020/9/30 13:26
 * @Desc: format response
 */
package response

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HttpResponse struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

func Ok(c echo.Context) error {
	return result(c, nil, "success")
}

func OkWithMessage(c echo.Context, message string) error {
	return result(c, nil, message)
}

func OkWithData(c echo.Context, data interface{}) error {
	return result(c, data, "操作成功")
}

func OkWithDetailed(c echo.Context, data interface{}, message string) error {
	return result(c, data, message)
}

func Failed(c echo.Context, status int) error {
	return FailWithMsg(c, status, "操作失败")
}

func FailWithMsg(c echo.Context, status int, err interface{}) error {
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

	return resultWithStatus(c, status, false, nil, msg)
}

func result(c echo.Context, data interface{}, msg string) error {
	return resultWithStatus(c, http.StatusOK, true, data, msg)
}

func resultWithStatus(c echo.Context, statusCode int, success bool, data interface{}, msg interface{}) error {
	return c.JSON(statusCode, HttpResponse{
		success,
		msg,
		data,
	})
}
