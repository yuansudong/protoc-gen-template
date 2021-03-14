package main

import (
	"log"

	"github.com/yuansudong/gengo"
	plugin "github.com/yuansudong/gengo/plugin"
)

type pathType int

const (
	pathTypeImport pathType = iota
	pathTypeSourceRelative
)

type gen struct {
	reg        *gengo.Registry
	pathType   pathType
	modulePath string
}

// Method 用于描述一个服务名称
type Method struct {
	Name string
	TReq string
	TRsp string
}

// Service 用于描述一个服务
type Service struct {
	Version string
	Methods []Method
	Name    string
}

// Args 用于描述一个参数
type Args struct {
	DateTime    string
	IsHave      bool
	PackageName string
	Services    []Service
	Imports     map[string]bool
}

// New 新建立一个
func New(reg *gengo.Registry) *gen {
	var pathType pathType
	switch *_PathTypeString {
	case "", "import":
		// paths=import is default
	case "source_relative":
		pathType = pathTypeSourceRelative
	default:
		log.Fatalf(`Unknown path type %s: want "import" or "source_relative".`, *_PathTypeString)
	}

	return &gen{
		reg:        reg,
		pathType:   pathType,
		modulePath: *_ModulePathString,
	}
}

// Generate 用于执行生成操作
func (g *gen) Generate(targets []*gengo.File) ([]*plugin.CodeGeneratorResponse_File, error) {
	var files []*plugin.CodeGeneratorResponse_File
	for _, file := range targets {
		files = append(files, g.GoTemplate(file))
	}
	return files, nil
}
