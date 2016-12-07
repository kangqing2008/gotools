package main

import (
	"fmt"
	"kangqing2008/gotools/tools/code"
	"flag"
	"strings"
)


func main() {
	use := flag.String("Use","deploy",`参数Use用于指定程序用途
	[deploy]为生成部署文件
	[compare]比较代码差异`)
	flag.Parse()
	fmt.Println("Use 的 value:",*use)
	if strings.ToLower(*use) == "deploy" {
		fmt.Println("您选择的用途为【生成部署文件】，读取deploy-config.json执行相关操作...")
		code.Deploy()
	}else if  strings.ToLower(*use) == "compare" {
		fmt.Println("您选择的用途为【比较代码差异】，读取compare-config.json执行相关操作...")
		code.Compare()
	} else {
		fmt.Println("为参数Use指定了无效的值[",*use,"],程序终止!")
	}
}
