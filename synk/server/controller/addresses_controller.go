package controller

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

//思路：
//1.获取电脑在各个局域网的IP地址
//2.转为json写入HTTP响应

func AddressesController(c *gin.Context) {
	addrs, _ := net.InterfaceAddrs() //获取当前电脑的所有ip地址
	var result []string
	for _, address := range addrs { //遍历所有ip地址
		// check the address type and if it is not a loopback（回环） the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil { //将地址存储到result切片中
				result = append(result, ipnet.IP.String())
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"addresses": result}) //作为一个json返回给前端
}
