package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/yuansudong/gengo"
)

// PluginName 用于描述一个插件名称
const PluginName = "protoc-gen-template:"

// TestFile 用于描述一个测试文件
const TestFile = "C:/Users/Administrator/Desktop/new/golang/src/dev.phoenix-ea.com/yuansudong/tools.hfdy.com/protoc-gen-mock/t1.test.bin"

var (
	_Lang             = flag.String("lang", "", "指定一个开发语言")
	_ImportPrefix     = flag.String("import_prefix", "", "prefix to be added to go package paths for imported proto files")
	_ImportPath       = flag.String("import_path", "", "used as the package if no input files declare go_package. If it contains slashes, everything up to the rightmost slash is ignored.")
	_PathTypeString   = flag.String("paths", "", "specifies how the paths of generated files are structured")
	_ModulePathString = flag.String("module", "", "specifies a module prefix that will be stripped from the go package to determine the output directory")
)

func main() {
	flag.Parse()
	startFromStdin()
	//startFromFile()
}

// SaveReq 用于保存请求到一个文件
func SaveReq() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Println(err.Error())
	}
	ioutil.WriteFile(TestFile, data, 0755)
}

// _StartFromFile 测试的时候从文件中读取数据
func _StartFromFile() {
	in, err := os.Open(TestFile)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(-1)
	}
	req, err := gengo.GetRequest(in)
	if err != nil {
		gengo.WriteError(err)
		return
	}
	if req.Parameter != nil {
		req.GetParameter()
	}
	reg := gengo.NewRegistry()
	if err = reg.Load(req); err != nil {
		gengo.WriteError(err)
		return
	}
	var targets []*gengo.File
	for _, target := range req.FileToGenerate {
		f, err := reg.LookupFile(target)
		if err != nil {
			log.Fatal(err)
		}
		targets = append(targets, f)
	}
	g := New(reg)
	out, err := g.Generate(targets)
	if err != nil {
		gengo.WriteError(err)
		return
	}
	gengo.WriteFiles(out)
}

// startFromStdin 用于开始从标准输入
func startFromStdin() {
	in := os.Stdin
	req, err := gengo.GetRequest(in)
	if err != nil {
		gengo.WriteError(err)
		return
	}

	if req.Parameter != nil {
		// 解析参数
	}
	reg := gengo.NewRegistry()
	if err = reg.Load(req); err != nil {
		gengo.WriteError(err)
		return
	}
	var targets []*gengo.File
	for _, target := range req.FileToGenerate {
		f, err := reg.LookupFile(target)
		if err != nil {
			log.Fatal(err)
		}
		targets = append(targets, f)
	}
	g := New(reg)
	out, err := g.Generate(targets)
	if err != nil {
		gengo.WriteError(err)
		return
	}
	gengo.WriteFiles(out)
}
