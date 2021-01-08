/**
 * @Author: lzw5399
 * @Date: 2021/01/07 14:22
 * @Desc: home page controller
 */
package controller

import (
	"net/http"

	"workflow/src/global"
	"workflow/src/global/response"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"AppName": global.BANK_CONFIG.App.Name,
	})
}

func Liveliness(c *gin.Context) {
	response.Ok(c)
}

func Readiness(c *gin.Context) {
	response.Ok(c)
}
