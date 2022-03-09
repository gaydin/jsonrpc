package jsonrpc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type (
	Registry struct {
		middlewares []MiddlewareFunc
		m           *sync.Mutex
		r           map[string]meta
	}

	meta struct {
		handler     HandlerFunc
		middlewares []MiddlewareFunc
	}
)

func New(middlewares ...MiddlewareFunc) *Registry {
	return &Registry{
		middlewares: middlewares,
		m:           &sync.Mutex{},
		r:           make(map[string]meta),
	}
}

func (mr *Registry) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requests, batch, err := ParseRequest(r)
	if err != nil {
		err := sendResponse(w, []*Response{
			{
				Version: Version,
				Error:   err,
			},
		}, false)
		if err != nil {
			fmt.Fprint(w, "Failed to encode error objects")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	resp := make([]*Response, len(requests))
	for i := range requests {
		resp[i] = mr.invokeMethod(r.Context(), requests[i])
	}

	if err := sendResponse(w, resp, batch); err != nil {
		fmt.Fprint(w, "Failed to encode result objects", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (mr *Registry) RegisterMethods(methods []Method, middlewares ...MiddlewareFunc) error {
	for i := range methods {
		if err := mr.RegisterMethod(methods[i], middlewares...); err != nil {
			return err
		}
	}

	return nil
}

func (mr *Registry) RegisterMethod(method Method, middlewares ...MiddlewareFunc) error {
	if method.method == "" || method.handler == nil {
		return errors.New("jsonrpc: method name and function should not be empty")
	}

	mr.m.Lock()
	mr.r[method.method] = meta{
		handler:     method.handler,
		middlewares: middlewares,
	}
	mr.m.Unlock()
	return nil
}

func (mr *Registry) invokeMethod(c context.Context, r *Request) *Response {
	var method meta
	res := newResponse(r)
	method, res.Error = mr.takeMethod(r)
	if res.Error != nil {
		return res
	}

	wrappedContext := WithRequestID(c, r.ID)
	wrappedContext = WithMethodName(wrappedContext, r.Method)
	handler := applyMiddleware(method.handler, mr.middlewares...)
	handler = applyMiddleware(handler, method.middlewares...)
	res.Result, res.Error = handler(wrappedContext, r.Params)
	if res.Error != nil {
		res.Result = nil
	}
	return res
}

func (mr *Registry) takeMethod(r *Request) (meta, *Error) {
	if r.Method == "" || r.Version != Version {
		return meta{}, ErrInvalidParams()
	}

	md, ok := mr.r[r.Method]
	if !ok {
		return meta{}, ErrMethodNotFound()
	}

	return md, nil
}
