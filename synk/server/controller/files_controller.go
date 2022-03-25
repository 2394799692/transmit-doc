package controller

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//思路：
//1.获取go执行文件所在目录
//2.在该目录创建uploads目录
//3.将文件保存为另一个文件
//4.返回后者的下载路径

func FilesController(c *gin.Context) {
	file, err := c.FormFile("raw") //读取用户上传的文件
	if err != nil {
		log.Fatal(err)
	}
	exe, err := os.Executable() //输出一个临时文件的路径
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exe) //用于返回指定路径中除最后一个元素以外的所有元素。
	if err != nil {
		log.Fatal(err)
	}
	filename := uuid.New().String() //创建uploads文件
	uploads := filepath.Join(dir, "uploads")
	err = os.MkdirAll(uploads, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	fullpath := path.Join("uploads", filename+filepath.Ext(file.Filename)) //获取本地文件路径
	fileErr := c.SaveUploadedFile(file, filepath.Join(dir, fullpath))      //存储用户上传的文件
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath})
}
