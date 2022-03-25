package main

import (
	"os"
	"os/exec"
	"os/signal"
	"synk/config"
	"synk/server"
)

func main() {
	go server.Run()                 //启动gin协程
	cmd := startBrowser()           //启动ui界面，打开Chrome
	chSignal := listenToInterrupt() //监听关闭信号
	<-chSignal                      //从chan中读出中断信号
	cmd.Process.Kill()              //关闭谷歌浏览器进程
}

func startBrowser() *exec.Cmd {
	// 先写死路径，后面再照着 lorca 改
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:"+config.GetPort()+"/static/index.html")
	cmd.Start()
	return cmd
}

func listenToInterrupt() chan os.Signal {
	chSignal := make(chan os.Signal, 1)   //接收系统信号
	signal.Notify(chSignal, os.Interrupt) //用户按下ctrl+c，就发送系统信号
	return chSignal
}
