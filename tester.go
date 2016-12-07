//测试用的小工具
//康庆 2016-11-20
package main

import (
	"strings"
	"fmt"
	"path/filepath"
	"os"
	"flag"
)

func main() {
	testFlag()
}

func testInput(){
	var use int = -1
	for true {
		fmt.Println("\n请选择工具用途[1]为生成部署文件，[2]为比较代码差异。【默认为：1,回车即可】?")
		fmt.Scanf("%d", &use)
		fmt.Println("输入的是：",use)
		if use == 1 || use == -1 || use == 2{
			break
		}else{
			fmt.Println("输入错误,您输入的是[",use,"]")
		}
	}
}

func testFlag(){
	Use := flag.String("Use","deploy",`参数Use用于指定程序用途
	[deploy]为生成部署文件
	[compare]比较代码差异`)
	flag.Parse()
	fmt.Println("Use 的 value:",*Use)
}

func testIndex(){
	str := `D:\WORKSPACES\GO_WORKSPACES\BeegoDemo\src\szboanda.com\\tools`
	fmt.Println(strings.Index(str,"BeegoDemo\\src"))
}

var javaFiles []string

func testWalk(){
	filepath.Walk("F:\\SOFTWARE",display)
}

func display(path string,info os.FileInfo,err error)error{
	if !info.IsDir(){
		if strings.HasSuffix(strings.ToLower(path),".java"){
			fmt.Println(path)
		}
	}
	return nil
}

func testAppend(){
	for i:=0;i<100;i++{
		javaFiles = append(javaFiles,"abc")
	}
	for _,str := range javaFiles {
		fmt.Println(str)
	}
}

