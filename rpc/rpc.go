package rpc

import (
	"context"

	"github.com/T-c33/smicro/meta"
	_ "github.com/T-c33/smicro/registry/etcd"
)

// service: 服务提供方的服务名
// method: 要调用服务的方法
// caller: 调用者的名字
func InitRpcMeta(ctx context.Context, service, method, caller string) context.Context {
	return meta.InitRpcMeta(ctx, service, method, caller)
}

/*
func BuildClientMiddleware(handle middleware.MiddlewareFunc) middleware.MiddlewareFunc {
	var mids []middleware.Middleware
	if len(mids) == 0 {
		return handle
	}

	m := middleware.Chain(mids[0], mids...)
	return m(handle)
}
*/
