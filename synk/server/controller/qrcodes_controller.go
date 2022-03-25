package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

//思路：
//1.获取文本内容
//2.将文本转为图片（用qrcode库）
//3.将图片写入HTTP响应

func QrcodesController(c *gin.Context) {
	if content := c.Query("content"); content != "" {
		png, err := qrcode.Encode(content, qrcode.Medium, 256) //把文本编程png
		if err != nil {
			log.Fatal(err)
		}
		c.Data(http.StatusOK, "image/png", png) //把图片传给前端，用于展示
	} else {
		c.Status(http.StatusBadRequest) //否则返回一个错误
	}
}
