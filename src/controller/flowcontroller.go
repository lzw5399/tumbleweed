/**
 * @Author: lzw5399
 * @Date: 2020/9/30 23:24
 * @Desc: ocr related functionality
 */
package controller

import (
	"workflow/src/global/response"
	
	"github.com/gin-gonic/gin"
)

func ScanFile(c *gin.Context) {
	response.Ok(c)
}
