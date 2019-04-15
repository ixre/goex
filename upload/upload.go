package main

import (
	"flag"
	"fmt"
	"github.com/ixre/goex/echox"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var ch chan bool
	var port int
	var confPath string
	//var uploadDir string
	flag.IntVar(&port, "port", 1419, "listen tcp port")
	flag.StringVar(&confPath, "conf", "uams.conf", "config file")
	//flag.StringVar(&uploadDir, "upload-dir", conf.UploadSaveDir, "upload directory path")
	flag.Parse()
	log.SetFlags(log.Lmicroseconds | log.LstdFlags)
	/*
		if uploadDir != conf.UploadSaveDir { // 设置上传文件目录
			conf.UploadSaveDir = uploadDir
			log.Println("[ Upload][ CONF]: upload directory is [", uploadDir, "]")
		}*/
	go signalNotify(ch)
	go serveHttp(port)
	<-ch
}

// 监听进程信号,并执行操作。比如退出时应释放资源
func signalNotify(c chan bool) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM,
		syscall.SIGKILL, syscall.SIGUSR2)
	for {
		select {
		case sig := <-ch:
			switch sig {
			case syscall.SIGHUP, syscall.SIGKILL, syscall.SIGTERM: // 退出时
				log.Println("[ System][ TERM] - program has exit !")
				if c != nil {
					close(c)
				}
			case syscall.SIGUSR2:
			}
		}
	}
}

// 运行http服务
func serveHttp(port int) {
	// 新建应用
	e := echox.New()
	// 启动服务
	portStr := fmt.Sprintf(":%d", port)
	log.Println("[ Upload][ Serve]: start upload server on port ", portStr)
	err := http.ListenAndServe(portStr, e)
	if err != nil {
		log.Println("[ Upload][ EXIT]: ", err)
		os.Exit(1)
	}
}
