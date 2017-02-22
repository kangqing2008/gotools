/*
文件操作相关的工具类
@康庆 2016-11-24
*/
package file

import (
	"path/filepath"
	"os"
	"os/exec"
	"strings"
	"sort"
	"io"
	"fmt"
	"bufio"
)

const StrPathSeparator string = string(os.PathSeparator)

//获取当前运行的文件的路径，并将路径中的\替换成为/
func ExecFilePath()string{
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	return strings.Replace(path, "\\", "/", -1)
}

//获取当前运行的文件所在的目录
func CurrentPath()string{
	exe := ExecFilePath()
	return Substr(exe, 0, strings.LastIndex(exe, "/"))

}

//截取字符串
func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func Exists(filename string)(bool,os.FileInfo){
	info,err := os.Stat(filename)
	return err == nil || os.IsExist(err),info
}

func ListFiles(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

//获取制定的文件所在的目录
//只针对文件
func GetFileDir(filename string)string{
	filename = strings.Replace(filename,"\\","/",-1)
	if strings.LastIndex(filename,"/") > -1{
		return Substr(filename, 0, strings.LastIndex(filename, "/"))
	}else{
		return filename
	}
}

//拷贝文件
func CopyFile(src,dest string)(int64,error){
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
	}
	defer srcFile.Close()
	desFile, err := os.Create(dest)
	if err != nil {
		fmt.Println(err)
	}
	defer desFile.Close()
	return io.Copy(desFile, srcFile)
}

func ReadLine(filename string)([]string,error){
	f, err := os.Open(filename)
	if err != nil {
		return nil,err
	}
	var results []string
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		//fmt.Println(line)
		if err != nil {
			if err == io.EOF {
				return results,nil
			}else{
				fmt.Println("读取文件出错：",filename)
			}
			return results,err
		}else{
			results = append(results,line)
		}
	}
	return results,nil
}