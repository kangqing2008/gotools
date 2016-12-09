package main

import (
	"net/http"
	"log"
	"flag"
	"strconv"
	"fmt"
)

func main() {
	port := flag.Int("port",10000,"请指定监听端口")
	fmt.Println("启动http服务,监听[",*port,"]端口......")
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(*port), http.FileServer(http.Dir("./"))))
}
