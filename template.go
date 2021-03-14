package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/yuansudong/gengo"
	plugin "github.com/yuansudong/gengo/plugin"
)

const codeFileTemplate = `
// 以此等代码注释开头的,都属于工具生成.请不要人为改变.
// 寄语:  工作是为了生活,而不是生活为了工作. 如何能在最少的时间内完成工作,这就是工具存在的意义!
///////////////////////////////////////////////////////////////////////////
// 		   Package: {{ .PackageName }}
// 		   Description: 服务于火凤电艺时制作
//		   Author: 袁苏东
// 		   Email: 1145919989@qq.com
// 		   ProtocGenMock: 1.0.0
// 		   Protoc: unknown
//		   UpdateTime: {{ .DateTime }}
// 		   Company: 火凤电艺
//		   CreateTime: 2020	
///////////////////////////////////////////////////////////////////////////
package {{ .PackageName }}
{{ if .IsHave }}
import (
	"context"
	"sync"
	"dev.phoenix-ea.com/gaonankai/base.hfdy.com/gmiddle"
	"dev.phoenix-ea.com/gaonankai/base.hfdy.com/gchain"
	"dev.phoenix-ea.com/gaonankai/base.hfdy.com/state"
	"dev.phoenix-ea.com/gaonankai/proto.hfdy.com/gopb/error/bad/pbbad"
	"dev.phoenix-ea.com/gaonankai/proto.hfdy.com/gopb/error/unimpl/pbunimpl"
{{ range $importKey,$importVal :=  .Imports  }}
	"{{ $importKey }}"
{{ else }}
{{ end }}

)
{{ else }}
{{ end }}

{{ range $serviceKey,$serviceVal :=  .Services }}
// 全局变量定义处
var xxx_{{ $serviceVal.Version }}_once sync.Once

var xxx_{{ $serviceVal.Version }}_inst *Impl{{ $serviceVal.Version }}Service


{{ range $methodKey,$methodVal := $serviceVal.Methods }}
type	{{ $methodVal.Name }}Handler func(context.Context, *{{ $methodVal.TReq }})(*{{ $methodVal.TRsp }},error)
{{ else }}
{{ end  }}


type Impl{{ $serviceVal.Version }}Service struct {
	chain           gchain.ChanHandler
{{ range $methodKey,$methodVal := $serviceVal.Methods }}
	Real{{ $methodVal.Name }}Func {{ $methodVal.Name }}Handler
	Info{{ $methodVal.Name }}    *gchain.ServiceInfo
{{ else }}
{{ end  }}
}  

// NewImpl{{ $serviceVal.Version }}Service 用于新实例化一个实现服务
func NewImpl{{ $serviceVal.Version }}Service() *Impl{{ $serviceVal.Version }}Service {
	inst := new(Impl{{ $serviceVal.Version }}Service)
{{ range $methodKey,$methodVal := $serviceVal.Methods }}
	inst.Real{{ $methodVal.Name }}Func = func(ctx context.Context, req *{{ $methodVal.TReq }})(rsp *{{ $methodVal.TRsp }},err error) {
		return nil, state.Unimplemented(pbunimpl.Unimpl_U_METHOD,"{{ $methodVal.Name }} Method Not Yet Implemented")
	}
	inst.Info{{ $methodVal.Name }} = &gchain.ServiceInfo{ Service: "{{ $serviceVal.Name }}" ,Version: "{{ $serviceVal.Version }}", Method: "{{ $methodVal.Name }}"}
{{ else }}
{{ end  }}
	return inst
}
// WithChain 用于设置执行链
func (i *Impl{{ $serviceVal.Version }}Service) WithChain(chains ...gchain.ChanHandler) (*Impl{{ $serviceVal.Version }}Service) {
	i.chain = gchain.ChainUnaryServer(chains...)
	return i
}
{{ range $methodKey,$methodVal := $serviceVal.Methods }}
func (i *Impl{{ $serviceVal.Version }}Service) With{{ $methodVal.Name }}(fHandle {{ $methodVal.Name }}Handler) *Impl{{ $serviceVal.Version }}Service {
	i.Real{{ $methodVal.Name }}Func = fHandle
	return i
}

// Deal{{ $methodVal.Name }} ....
func (i *Impl{{ $serviceVal.Version }}Service) Deal{{ $methodVal.Name }}(ctx context.Context, tmpReq interface{}) (interface{}, error) {
	req := tmpReq.(*{{ $methodVal.TReq }})
	if err := req.Validate(); err != nil {
		return nil, state.Bad(pbbad.Bad_B_ARG, err.Error())
	}
	mock, err := gmiddle.Mock(ctx)
	if err != nil {
		return nil, state.Bad(pbbad.Bad_B_ARG, err.Error())
	}
	if mock {
		rsp, err := new({{ $methodVal.TRsp }}).Mock()
		if err != nil {
			return nil, state.Bad(pbbad.Bad_B_ARG, err.Error())
		}
		return rsp, nil
	}
	return i.Real{{ $methodVal.Name }}Func(ctx, req)
}


// {{ $methodVal.Name }} 用于描述一个在线请求
func (i *Impl{{ $serviceVal.Version }}Service) {{ $methodVal.Name }}(ctx context.Context, req *{{ $methodVal.TReq }}) (rsp *{{ $methodVal.TRsp }}, err error) {
	tmpRsp, err := i.chain(ctx, req, i.Info{{ $methodVal.Name }}, i.Deal{{ $methodVal.Name }})
	if err != nil {
		return nil, err
	}
	return tmpRsp.(*{{ $methodVal.TRsp }}), nil
}
{{ else }}
{{ end }}

// GetImpl{{ $serviceVal.Version }}Instance 用于获取一个实例
func GetImpl{{ $serviceVal.Version }}Instance() *Impl{{ $serviceVal.Version }}Service {
	xxx_{{ $serviceVal.Version }}_once.Do(func(){
		xxx_{{ $serviceVal.Version }}_inst = NewImpl{{ $serviceVal.Version }}Service()
	})
	return xxx_{{ $serviceVal.Version }}_inst
}
{{ else }}
{{ end }}
`

// GoTemplate 用于生成相关的模板
func (g *gen) GoTemplate(file *gengo.File) *plugin.CodeGeneratorResponse_File {
	var err error
	as := &Args{}
	as.Imports = make(map[string]bool)
	as.PackageName = file.GoPkg.Name
	as.DateTime = time.Now().Local().String()
	buf := bytes.NewBuffer(make([]byte, 0, 40960))
	rspFile := new(plugin.CodeGeneratorResponse_File)
	for _, service := range file.Services {
		g.GoService(file, as, service)
	}
	if len(as.Services) != 0 {
		as.IsHave = true
	}

	tp := template.New("template.service")
	if tp, err = tp.Parse(codeFileTemplate); err != nil {
		log.Fatalln(err.Error())
	}
	if err = tp.Execute(buf, as); err != nil {
		log.Fatalln(err.Error())
	}
	name, err := g.GetAllFilePath(file)
	if err != nil {
		log.Println(PluginName, err.Error())
		os.Exit(-1)
	}
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	output := fmt.Sprintf("%s.pb.template.go", base)
	rspFile.Name = proto.String(output)
	rspFile.Content = proto.String(buf.String())
	return rspFile
}

//  GoMessage 用于处理消息
func (g *gen) GoService(rootFile *gengo.File, args *Args, service *gengo.Service) {
	sInst := Service{}
	sInst.Version = service.GetName()

	for _, mtd := range service.Methods {
		tMtd := Method{}
		tMtd.Name = mtd.GetName()
		tMtd.TReq = mtd.RequestType.GoType(rootFile.GoPkg.Path)
		if mtd.RequestType.File.GoPkg.Path != rootFile.GoPkg.Path {
			args.Imports[mtd.RequestType.File.GoPkg.Path] = true
		}
		tMtd.TRsp = mtd.ResponseType.GoType(rootFile.GoPkg.Path)
		if mtd.ResponseType.File.GoPkg.Path != rootFile.GoPkg.Path {
			args.Imports[mtd.ResponseType.File.GoPkg.Path] = true
		}
		sInst.Methods = append(sInst.Methods, tMtd)
	}
	args.Services = append(args.Services, sInst)
}

func (g *gen) GetAllFilePath(file *gengo.File) (string, error) {
	name := file.GetName()
	switch {
	case g.modulePath != "" && g.pathType != pathTypeImport:
		return "", errors.New("cannot use module= with paths=")

	case g.modulePath != "":
		trimPath, pkgPath := g.modulePath+"/", file.GoPkg.Path+"/"
		if !strings.HasPrefix(pkgPath, trimPath) {
			return "", fmt.Errorf("%v: file go path does not match module prefix: %v", file.GoPkg.Path, trimPath)
		}
		return filepath.Join(strings.TrimPrefix(pkgPath, trimPath), filepath.Base(name)), nil

	case g.pathType == pathTypeImport && file.GoPkg.Path != "":
		return fmt.Sprintf("%s/%s", file.GoPkg.Path, filepath.Base(name)), nil

	default:
		return name, nil
	}
}
