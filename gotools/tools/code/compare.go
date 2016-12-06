package code

import (
	"os"
	"fmt"
	"strings"
	"io/ioutil"
	"encoding/json"
	"io"
	"path/filepath"
	"bytes"
	"kangqing/gotools/tools/file"
)

const btips string = `**************************************************************
*  配置文件compare-config.json中参数的含义：
*  PathA:第一个目录
*  PathB:第二个目录
*  DestPath：目标目录
*     例如：D:\\updateCode
*  Mode:运行模式
      A
      B
      AB
**************************************************************`
const(
	MODE_A		string = "A"
	MODE_B		string = "B"
	MODE_AB		string = "A-B"

)

type CompConfModel struct {
	PathA  string `json:"PathA"`
	PathB  string `json:"PathB"`
	DestPath  string `json:"DestPath"`
	Mode   string `json:"Mode"`
}

var configModel CompConfModel
var pathAFiles	map[string]string = map[string]string{}
var pathBFiles	map[string]string = map[string]string{}
var copyFilesCount	int = 0

//入口函数，解析配置，查找修改过的源文件和对应的资源文件
func Compare() {
	fmt.Println(btips)
	//解析配置文件
	configModel = parseConfigModel()
	printCompConfig()
	findFiles(configModel.PathA,pathAFinder)
	findFiles(configModel.PathB,pathBFinder)
	results := compare(pathAFiles,pathBFiles)
	printDifference(results)
	copyDifference(results)
}

func printDifference(results map[string]string){
	fmt.Println("PathA目录下的文件数量：",len(pathAFiles))
	fmt.Println("PathB目录下的文件数量：",len(pathBFiles))
	fmt.Println("有差异文件的总数量：",len(results))

	var a,b,ab int
	for _,v := range results{
		if v == MODE_A{
			a++
		}else if v == MODE_B{
			b++
		}else if v == MODE_AB{
			ab++
		}
	}
	fmt.Println("PathA存在，PathB不存在的文件：",a)
	fmt.Println("PathB存在，PathA不存在的文件：",b)
	fmt.Println("PathA和PathB中都存在，但内容不同的文件数量：",ab)
}

func copyDifference(results map[string]string){
	for path,mode := range results{
		if MODE_A == mode{
			if "ALL" == configModel.Mode{
				src := configModel.PathA + path
				dest := configModel.DestPath + path + ".A"
				if _,err := copyFileAndMkDir(src,dest);err != nil{
					fmt.Println("拷贝文件[",src,"]到[",dest,"]时出错:",err)
				}
			}else if "A" == configModel.Mode{
				src := configModel.PathA + path
				dest := configModel.DestPath + path
				if _,err := copyFileAndMkDir(src,dest);err != nil{
					fmt.Println("拷贝文件[",src,"]到[",dest,"]时出错:",err)
				}
			}
		}else if MODE_B == mode {
			if "ALL" == configModel.Mode{
				src := configModel.PathB + path
				dest := configModel.DestPath + path + ".B"
				if _,err := copyFileAndMkDir(src,dest);err != nil{
					fmt.Println("拷贝文件[",src,"]到[",dest,"]时出错:",err)
				}
			}else if "B" == configModel.Mode{
				src := configModel.PathB + path
				dest := configModel.DestPath + path
				if _,err := copyFileAndMkDir(src,dest);err != nil{
					fmt.Println("拷贝文件[",src,"]到[",dest,"]时出错:",err)
				}

			}
		}else if MODE_AB == mode || "ALL" == configModel.Mode || "AB" == configModel.Mode {
			srcA := configModel.PathA + path
			destA := configModel.DestPath + path + ".A"
			if _,err := copyFileAndMkDir(srcA,destA);err != nil{
				fmt.Println("拷贝文件[",srcA,"]到[",destA,"]时出错:",err)
			}
			srcB := configModel.PathA + path
			destB := configModel.DestPath + path + ".B"
			if _,err := copyFileAndMkDir(srcB,destB);err != nil{
				fmt.Println("拷贝文件[",srcB,"]到[",destB,"]时出错:",err)
			}
		}
	}
}

func copyFileAndMkDir(src,dest string)(int64,error){
	fpath := file.GetFileDir(dest)
	if exists, _ := file.Exists(fpath); !exists {
		os.MkdirAll(fpath, os.ModeDir)
	}
	return file.CopyFile(src,dest);
}

func findFiles(path string,finder filepath.WalkFunc){
	filepath.Walk(path,finder)
}

func compare(fas,fbs map[string]string)map[string]string{
	results := map[string]string{}
	for ka := range fas{
		kb := strings.Replace(ka,configModel.PathA,configModel.PathB,-1)
		path := strings.Replace(ka,configModel.PathA,"",-1)
		if _,ok := fbs[kb];!ok{
			results[path] = MODE_A
		}else{
			if !equal(ka,kb){
				results[path] = MODE_AB
			}
		}
	}
	for kb := range fbs{
		ka := strings.Replace(kb,configModel.PathB,configModel.PathA,-1)
		path := strings.Replace(kb,configModel.PathB,"",-1)
		if _,ok := fas[ka];!ok{
			results[path] = MODE_B
		}
	}
	return results
}

func equal(fa,fb string)bool{
	sinfo, err := os.Lstat(fa)
	if err != nil {
		return false
	}
	dinfo, err := os.Lstat(fb)
	if err != nil {
		return false
	}
	if sinfo.Size() != dinfo.Size() {
		return false
	}
	return compareFile(fa, fb)
}

func compareFile(spath, dpath string) bool {
	sFile, err := os.Open(spath)
	if err != nil {
		return false
	}
	dFile, err := os.Open(dpath)
	if err != nil {
		return false
	}
	b := comparebyte(sFile, dFile)
	sFile.Close()
	dFile.Close()
	return b
}
//下面可以代替md5比较.
func comparebyte(sfile *os.File, dfile *os.File) bool {
	var sbyte []byte = make([]byte, 512)
	var dbyte []byte = make([]byte, 512)
	var serr, derr error
	for {
		_, serr = sfile.Read(sbyte)
		_, derr = dfile.Read(dbyte)
		if serr != nil || derr != nil {
			if serr != derr {
				return false
			}
			if serr == io.EOF {
				break
			}
		}
		if bytes.Equal(sbyte, dbyte) {
			continue
		}
		return false
	}
	return true
}

func pathAFinder(path string, info os.FileInfo, err error) error{
	if !info.IsDir(){
		pathAFiles[path] = "1"
	}
	return nil
}


func pathBFinder(path string, info os.FileInfo, err error) error{
	if !info.IsDir(){
		pathBFiles[path] = "1"
	}
	return nil
}

func printCompConfig() {
	fmt.Println("PathA:", configModel.PathA)
	fmt.Println("PathB:", configModel.PathB)
	fmt.Println("DestPath:", configModel.DestPath)
	fmt.Println("Mode:", configModel.Mode)
}

//解析配置文件
func parseConfigModel() CompConfModel {
	b, err := ioutil.ReadFile(file.CurrentPath() + file.StrPathSeparator + "compare-config.json")
	if err != nil {
		fmt.Println(err)
		panic("无法正确的读取配置文件compare-config.json,请确认文件是否存在!")
	}
	var config = CompConfModel{}
	err = json.Unmarshal(b, &config)
	if err != nil {
		fmt.Println(err)
		panic("无法正确的解析配置文件compare-config.json,请确认格式是否正确!")
	}
	return config
}

