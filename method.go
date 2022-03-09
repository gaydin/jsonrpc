package jsonrpc

import (
	"context"
	"encoding/json"
)

type (
	HandleParamsFunc[PARAMS, RESPONSE any] func(ctx context.Context, params PARAMS) (*RESPONSE, *Error)
	HandleFunc[RESPONSE any]               func(ctx context.Context) (*RESPONSE, *Error)
	HandlerFunc                            func(ctx context.Context, params *json.RawMessage) (result interface{}, err *Error)

	Method struct {
		method  string
		handler HandlerFunc
	}
)

func NewParamsMethod[PARAMS, RESPONSE any](method string, handler HandleParamsFunc[PARAMS, RESPONSE]) Method {
	return Method{
		method: method,
		handler: func(ctx context.Context, rawParams *json.RawMessage) (interface{}, *Error) {
			if rawParams == nil {
				return nil, ErrInvalidParams()
			}

			var params PARAMS
			if err := json.Unmarshal(*rawParams, &params); err != nil {
				return nil, ErrInvalidParams()
			}

			return handler(ctx, params)
		},
	}
}

func NewMethod[RESPONSE any](method string, handler HandleFunc[RESPONSE]) Method {
	return Method{
		method: method,
		handler: func(ctx context.Context, rawParams *json.RawMessage) (interface{}, *Error) {
			return handler(ctx)
		},
	}
}
