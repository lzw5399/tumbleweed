/**
 * @Author: lzw5399
 * @Date: 2021/01/07 14:22
 * @Desc: home page controller
 */
package controller

import (
	"net/http"

	"bank/distributedquery/src/global"
	"bank/distributedquery/src/global/response"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"AppName": global.BANK_CONFIG.App.Name,
	})
}

func Info(c *gin.Context) {
	info, err := "ok", ""
	if err != "" {
		response.Failed(c, http.StatusInternalServerError)
		return
	}

	response.OkWithData(c, info)
}

func Liveliness(c *gin.Context) {
	response.Ok(c)
}

func Readiness(c *gin.Context) {
	response.Ok(c)
}
