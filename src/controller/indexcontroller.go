/**
 * @Author: lzw5399
 * @Date: 2021/01/07 14:22
 * @Desc: home page controller
 */
package controller

import (
	"time"
	
	"workflow/src/global/response"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	response.OkWithData(c, time.Now())
}

func Liveliness(c *gin.Context) {
	response.Ok(c)
}

func Readiness(c *gin.Context) {
	response.Ok(c)
}
