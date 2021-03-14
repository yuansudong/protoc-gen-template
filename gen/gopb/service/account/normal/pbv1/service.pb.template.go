
// 以此等代码注释开头的,都属于工具生成.请不要人为改变.
// 寄语:  工作是为了生活,而不是生活为了工作. 如何能在最少的时间内完成工作,这就是工具存在的意义!
///////////////////////////////////////////////////////////////////////////
// 		   Package: pbv1
// 		   Description: 服务于火凤电艺时制作
//		   Author: 袁苏东
// 		   Email: 1145919989@qq.com
// 		   ProtocGenMock: 1.0.0
// 		   Protoc: unknown
//		   UpdateTime: 2021-03-14 13:18:23.5141936 +0800 CST
// 		   Company: 火凤电艺
//		   CreateTime: 2020	
///////////////////////////////////////////////////////////////////////////
package pbv1

import (
	"context"
	"sync"
	"dev.phoenix-ea.com/gaonankai/base.hfdy.com/gmiddle"
	"dev.phoenix-ea.com/gaonankai/base.hfdy.com/gchain"
	"dev.phoenix-ea.com/gaonankai/base.hfdy.com/state"
	"dev.phoenix-ea.com/gaonankai/proto.hfdy.com/gopb/error/bad/pbbad"
	"dev.phoenix-ea.com/gaonankai/proto.hfdy.com/gopb/error/unimpl/pbunimpl"



)



// 全局变量定义处
var xxx_V1_once sync.Once

var xxx_V1_inst *ImplV1Service



type	InnerLoginHandler func(context.Context, *C2SInnerLogin)(*S2CInnerLogin,error)

type	InnerSigninHandler func(context.Context, *C2SInnerSignin)(*S2CInnerSignin,error)



type ImplV1Service struct {
	chain           gchain.ChanHandler

	RealInnerLoginFunc InnerLoginHandler
	InfoInnerLogin    *gchain.ServiceInfo

	RealInnerSigninFunc InnerSigninHandler
	InfoInnerSignin    *gchain.ServiceInfo

}  

// NewImplV1Service 用于新实例化一个实现服务
func NewImplV1Service() *ImplV1Service {
	inst := new(ImplV1Service)

	inst.RealInnerLoginFunc = func(ctx context.Context, req *C2SInnerLogin)(rsp *S2CInnerLogin,err error) {
		return nil, state.Unimplemented(pbunimpl.Unimpl_U_METHOD,"InnerLogin Method Not Yet Implemented")
	}
	inst.InfoInnerLogin = &gchain.ServiceInfo{ Service: "" ,Version: "V1", Method: "InnerLogin"}

	inst.RealInnerSigninFunc = func(ctx context.Context, req *C2SInnerSignin)(rsp *S2CInnerSignin,err error) {
		return nil, state.Unimplemented(pbunimpl.Unimpl_U_METHOD,"InnerSignin Method Not Yet Implemented")
	}
	inst.InfoInnerSignin = &gchain.ServiceInfo{ Service: "" ,Version: "V1", Method: "InnerSignin"}

	return inst
}
// WithChain 用于设置执行链
func (i *ImplV1Service) WithChain(chains ...gchain.ChanHandler) (*ImplV1Service) {
	i.chain = gchain.ChainUnaryServer(chains...)
	return i
}

func (i *ImplV1Service) WithInnerLogin(fHandle InnerLoginHandler) *ImplV1Service {
	i.RealInnerLoginFunc = fHandle
	return i
}

// DealInnerLogin ....
func (i *ImplV1Service) DealInnerLogin(ctx context.Context, tmpReq interface{}) (interface{}, error) {
	req := tmpReq.(*C2SInnerLogin)
	if err := req.Validate(); err != nil {
		return nil, state.Bad(pbbad.Bad_B_ARG, err.Error())
	}
	mock, err := gmiddle.Mock(ctx)
	if err != nil {
		return nil, state.Bad(pbbad.Bad_B_ARG, err.Error())
	}
	if mock {
		rsp, err := new(S2CInnerLogin).Mock()
		if err != nil {
			return nil, state.Bad(pbbad.Bad_B_ARG, err.Error())
		}
		return rsp, nil
	}
	return i.RealInnerLoginFunc(ctx, req)
}


// InnerLogin 用于描述一个在线请求
func (i *ImplV1Service) InnerLogin(ctx context.Context, req *C2SInnerLogin) (rsp *S2CInnerLogin, err error) {
	tmpRsp, err := i.chain(ctx, req, i.InfoInnerLogin, i.DealInnerLogin)
	if err != nil {
		return nil, err
	}
	return tmpRsp.(*S2CInnerLogin), nil
}

func (i *ImplV1Service) WithInnerSignin(fHandle InnerSigninHandler) *ImplV1Service {
	i.RealInnerSigninFunc = fHandle
	return i
}

// DealInnerSignin ....
func (i *ImplV1Service) DealInnerSignin(ctx context.Context, tmpReq interface{}) (interface{}, error) {
	req := tmpReq.(*C2SInnerSignin)
	if err := req.Validate(); err != nil {
		return nil, state.Bad(pbbad.Bad_B_ARG, err.Error())
	}
	mock, err := gmiddle.Mock(ctx)
	if err != nil {
		return nil, state.Bad(pbbad.Bad_B_ARG, err.Error())
	}
	if mock {
		rsp, err := new(S2CInnerSignin).Mock()
		if err != nil {
			return nil, state.Bad(pbbad.Bad_B_ARG, err.Error())
		}
		return rsp, nil
	}
	return i.RealInnerSigninFunc(ctx, req)
}


// InnerSignin 用于描述一个在线请求
func (i *ImplV1Service) InnerSignin(ctx context.Context, req *C2SInnerSignin) (rsp *S2CInnerSignin, err error) {
	tmpRsp, err := i.chain(ctx, req, i.InfoInnerSignin, i.DealInnerSignin)
	if err != nil {
		return nil, err
	}
	return tmpRsp.(*S2CInnerSignin), nil
}


// GetImplV1Instance 用于获取一个实例
func GetImplV1Instance() *ImplV1Service {
	xxx_V1_once.Do(func(){
		xxx_V1_inst = NewImplV1Service()
	})
	return xxx_V1_inst
}

