package code

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	"kangqing2008/gotools/tools/file"
)

const tips string = `**************************************************************
*  此程序用于从{MainPath}查找修改时间大于{StartTime}但小于{EndTime}的文件，自动抽取更新过的文件到{DestPath}目录下
*  能够自动识别java文件修改后对应的class文件和其它资源文件
*  配置文件config.json中参数都必须制定，具体的含义：
*  StartTime:程序将根据搜索出所有修改时间晚于这个时间的文件
*     格式为：2015-05-21 18：00
*  EndTime:程序将根据搜索出所有修改时间早于这个时间的文件，设置成NOW表示当前时间
*     格式为：NOW 或者 2015-05-22 18：00
*  MainPath:指定Maven工程源码路径
*     例如：D:\\WORKSPACES\\NEW_SVN\\LIMS\\src\\main
*  DestPath：指定更新文件生成到哪个目录
*     例如：D:\\updateCode
**************************************************************
`

type ConfModel struct {
	MainPath  string `json:"MainPath"`
	DestPath  string `json:"DestPath"`
	StartTime string `json:"StartTime"`
	EndTime   string `json:"EndTime"`
	Start     time.Time
	End       time.Time
}

var javaFiles []string
var resFiles []string
var classFiles []string
var otherFiles []string
var conf ConfModel
var counter int32 = 0

//入口函数，解析配置，查找修改过的源文件和对应的资源文件
//
func Deploy() {
	fmt.Println(tips)
	//解析配置文件
	conf = parseConfig()
	printConfig()
	//遍历java目录，找出所有修改过的java文件和资源文件
	filepath.Walk(conf.MainPath+file.StrPathSeparator+"java", checkJavaFile)
	//遍历webapp目录，找出所有修改过的其它文件
	filepath.Walk(conf.MainPath+file.StrPathSeparator+"webapp", checkOtherFile)
	//转换文件地址，找出所有修改过的java文件和资源文件对应的编译后的文件
	convertClassFiles()
	copyClassFiles()
	copyOtherFiles()
	fmt.Println("\n修改过的Java文件:", len(javaFiles))
	fmt.Println("修改过的资源文件:", len(resFiles))
	fmt.Println("更新过的Class文件:", len(classFiles)-len(resFiles))
	fmt.Println("其它文件:", len(otherFiles))
	//displayAll()
}

func printConfig() {
	fmt.Println("\nMainPath:", conf.MainPath)
	fmt.Println("DestPath:", conf.DestPath)
	fmt.Println("StartTime:", conf.StartTime)
	fmt.Println("EndTime:", conf.EndTime, "\n")
}

//将classes目录下的所有文件copy到目标文件夹下
func copyClassFiles() {
	s := file.StrPathSeparator
	rpath := s + "webapp" + s + "WEB-INF" + s + "classes"
	//源文件路径
	src := conf.MainPath + rpath
	//目标文件路径
	dest := conf.DestPath + rpath
	for _, class := range classFiles {
		dclass := strings.Replace(class, src, dest, -1)
		fpath := file.GetFileDir(dclass)
		if exists, _ := file.Exists(fpath); !exists {
			os.MkdirAll(fpath, os.ModeDir)
		}
		_, err := file.CopyFile(class, dclass)
		if counter++; counter >= 66 {
			fmt.Println(".")
			counter = 0
		} else {
			fmt.Print(".")
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

//将其他文件copy到目标文件夹下
func copyOtherFiles() {
	s := file.StrPathSeparator
	//源文件路径
	src := conf.MainPath + s + "webapp"
	//目标文件路径
	dest := conf.DestPath + s + "webapp"
	for _, res := range otherFiles {
		dres := strings.Replace(res, src, dest, -1)
		fpath := file.GetFileDir(dres)
		if exists, _ := file.Exists(fpath); !exists {
			os.MkdirAll(fpath, os.ModeDir)
		}
		_, err := file.CopyFile(res, dres)
		if counter++; counter >= 66 {
			fmt.Println(".")
			counter = 0
		} else {
			fmt.Print(".")
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

func convertClassFiles() {
	s := file.StrPathSeparator
	p1 := conf.MainPath + s + "java"
	p2 := conf.MainPath + s + "webapp" + s + "WEB-INF" + s + "classes"
	s1 := ".java"
	s2 := ".class"
	for _, java := range javaFiles {
		class := strings.Replace(java, p1, p2, -1)
		class = strings.Replace(class, s1, s2, -1)
		//判断java文件对应的class文件是否存在，如果存在保存起来
		if exists, _ := file.Exists(class); exists {
			classFiles = append(classFiles, class)
			//获取到class文件对应的目录文件对象
			dir := file.Substr(class, 0, strings.LastIndex(class, "\\"))
			if exists, _ := file.Exists(dir); !exists {
				panic("无法找到Java文件[" + java + "]对应的文件class文件[" + class + "]可能是项目没有编译成功")
			}
			childens, err := file.ListFiles(dir)
			if err != nil {
				fmt.Println(err)
				panic("罗列目录下的文件列表出错")
			}
			//去掉class文件路径中的后缀,并将路径符号替换过来，以便于dir+children拼接
			className := strings.Replace(class, ".class", "", -1)
			for _, child := range childens {
				//如果class文件目录下的其它文件，满足三个条件，那么也需要copy出去
				//比class文件名称长，完全包含class，后缀是.class
				child = dir + s + child
				if strings.Index(child, className) > -1 && strings.HasSuffix(child, ".class") && len(child) > len(class) && strings.Index(child, "$") > -1 {
					classFiles = append(classFiles, child)
				}
			}
		} else {
			fmt.Println("******无法找到Java文件[" + java + "]对应的文件class文件[" + class + "]可能是项目没有编译成功")
		}
	}
	//将资源文件的路径也做转换后，放入classFiles
	for _, res := range resFiles {
		res = strings.Replace(res, p1, p2, -1)
		if exists, _ := file.Exists(res); exists {
			classFiles = append(classFiles, res)
		} else {
			fmt.Println("无法在classes目录找到资源文件：", res)
		}
	}
}

func displayAll() {
	fmt.Println("JavaFiles")
	for _, str := range javaFiles {
		fmt.Println(str)
	}
	fmt.Println("ResFiles")
	for _, str := range resFiles {
		fmt.Println(str)
	}
	fmt.Println("OtherFiles")
	for _, str := range otherFiles {
		fmt.Println(str)
	}
	fmt.Println("ClassFiles")
	for _, str := range classFiles {
		fmt.Println(str)
	}
}

//解析配置文件
func parseConfig() ConfModel {
	b, err := ioutil.ReadFile(file.CurrentPath() + file.StrPathSeparator + "deploy-config.json")
	if err != nil {
		fmt.Println(err)
		panic("无法正确的读取配置文件deploy-config.json,请确认文件是否存在!")
	}
	var conf = ConfModel{}
	var err1, err2 error
	err = json.Unmarshal(b, &conf)
	if err != nil {
		fmt.Println(err)
		panic("无法正确的解析配置文件deploy-config.json,请确认格式是否正确!")
	}
	conf.Start, err1 = time.Parse("2006-01-02 15:04", conf.StartTime)
	if err1 != nil {
		fmt.Println(err1, err2)
		panic("设置的StartTime无法解析成日期，请注意格式!")
	}

	conf.End, err2 = time.Parse("2006-01-02 15:04", conf.EndTime)
	if err2 != nil {
		conf.End = time.Now()
	}
	return conf
}

func checkJavaFile(path string, info os.FileInfo, err error) error {
	if info == nil {
		return nil
	}
	if !info.IsDir() {
		//先判断时间,如果修改时间不在监控范围内，就不做任何处理
		if info.ModTime().Before(conf.Start) || info.ModTime().After(conf.End) {
			return nil
		}
		//如果是java文件就放入javaFiles切片
		if strings.HasSuffix(strings.ToLower(path), ".java") {
			javaFiles = append(javaFiles, path)
			//否则放入resFiles切片
		} else {
			resFiles = append(resFiles, path)
		}
	}
	return nil
}

func checkOtherFile(path string, info os.FileInfo, err error) error {
	if info == nil {
		return nil
	}
	if !info.IsDir() {
		//需要忽略掉WEB-INF/classes目录下的文件
		if strings.Index(path, "WEB-INF\\classes") > -1 {
			return nil
		}
		if info.ModTime().After(conf.Start) && info.ModTime().Before(conf.End) {
			otherFiles = append(otherFiles, path)
		}
	}
	return nil
}
