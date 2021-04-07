package middleware

import (
	"context"
)

type MiddlewareFunc func(ctx context.Context, req interface{}) (resp interface{}, err error)

type Middleware func(MiddlewareFunc) MiddlewareFunc

func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next MiddlewareFunc) MiddlewareFunc {
		for i := len(others) - 1; i >= 0; i-- {
			next = others[i](next)
		}
		return outer(next)
	}
}
