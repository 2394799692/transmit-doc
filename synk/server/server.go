package server

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"synk/config"
	c "synk/server/controller"
	"synk/server/ws"
)

//go:embed frontend/dist/*
var FS embed.FS

//可以在Go语言应用程序中包含任何文件、目录的内容
//也就是说我们可以把文件以及目录中的内容都打包到生成的Go语言应用程序中了，
//部署的时候，直接扔一个二进制文件就可以了，不用再包含一些静态文件了，因为它们已经被打包到生成的应用程序中了。
//静态文件的Web托管,案例：
//var static embed.FS
//func main() {
//	http.ListenAndServe(":8080", http.FileServer(http.FS(static)))
//}

func Run() {
	hub := ws.NewHub()
	go hub.Run()               //启动websocket服务
	gin.SetMode(gin.DebugMode) //gin开发模式
	//使用不同运行模式方便应对不同场景，比如debug模式下，output format不同，logger也不同。
	router := gin.Default() //创建新的引擎，实现服务器监听状态
	//r := gin.Default() 创建带有默认中间件的路由
	//r :=gin.new()      创建带有没有中间件的路由
	//中间间：将具体业务和底层逻辑解耦的组件。
	//需要利用服务的人（前端写业务的），不需要知道底层逻辑（提供服务的）的具体实现，只要拿着中间件结果来用就好了。
	staticFiles, _ := fs.Sub(FS, "frontend/dist")          //把打包好的静态文件变成一个结构化的目录
	router.POST("/api/v1/files", c.FilesController)        //上传文件
	router.GET("/api/v1/qrcodes", c.QrcodesController)     //将局域网ip变为二维码
	router.GET("/uploads/:path", c.UploadsController)      //下载接口
	router.GET("/api/v1/addresses", c.AddressesController) //获取当前局域网ip
	router.POST("/api/v1/texts", c.TextsController)        //上传文本
	router.GET("/ws", func(c *gin.Context) {               //上下文是一个结构体
		ws.HttpController(c, hub) //websocak，实现手机穿文件到电脑,将http请求升级为websocket
	})
	router.StaticFS("/static", http.FS(staticFiles)) //访问本地文件,加载前端页面
	router.NoRoute(func(c *gin.Context) {            //设置默认路由,防止文件路径出错，如果出错返回404如果该目录下没有文件则显示默认页面index
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/static/") {
			reader, err := staticFiles.Open("index.html")
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close() //defer表示go会在恰当时间关闭（垃圾回收机制）
			stat, err := reader.Stat()
			if err != nil {
				log.Fatal(err)
			}
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html;charset=utf-8", reader, nil)
		} else {
			c.Status(http.StatusNotFound) //返回状态码404
		}
	})
	router.Run(":" + config.GetPort()) //监听端口
}
