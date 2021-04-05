package router

import (
	"context"

	"smicro/meta"
	"smicro/server"

	"smicro/tools/smicro/output/generate/com/google/hello"

	"smicro/tools/smicro/output/controller"
)

type RouterServer struct{}

func (s *RouterServer) SayHello(ctx context.Context, r *hello.HelloRequest) (resp *hello.HelloResponse, err error) {

	ctx = meta.InitServerMeta(ctx, "hello", "SayHello")
	mwFunc := server.BuildServerMiddleware(mwSayHello)
	mwResp, err := mwFunc(ctx, r)
	if err != nil {
		return
	}

	resp = mwResp.(*hello.HelloResponse)
	return
}

func mwSayHello(ctx context.Context, request interface{}) (resp interface{}, err error) {

	r := request.(*hello.HelloRequest)
	ctrl := &controller.SayHelloController{}
	err = ctrl.CheckParams(ctx, r)
	if err != nil {
		return
	}

	resp, err = ctrl.Run(ctx, r)
	return
}

func (s *RouterServer) SayHelloV2(ctx context.Context, r *hello.HelloRequest) (resp *hello.HelloResponse, err error) {

	ctx = meta.InitServerMeta(ctx, "hello", "SayHelloV2")
	mwFunc := server.BuildServerMiddleware(mwSayHelloV2)
	mwResp, err := mwFunc(ctx, r)
	if err != nil {
		return
	}

	resp = mwResp.(*hello.HelloResponse)
	return
}

func mwSayHelloV2(ctx context.Context, request interface{}) (resp interface{}, err error) {

	r := request.(*hello.HelloRequest)
	ctrl := &controller.SayHelloV2Controller{}
	err = ctrl.CheckParams(ctx, r)
	if err != nil {
		return
	}

	resp, err = ctrl.Run(ctx, r)
	return
}
