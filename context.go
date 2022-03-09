package jsonrpc

import (
	"context"
	"encoding/json"
)

type (
	requestIDKey  struct{}
	methodNameKey struct{}
)

func RequestID(c context.Context) *json.RawMessage {
	return c.Value(requestIDKey{}).(*json.RawMessage)
}

func WithRequestID(c context.Context, id *json.RawMessage) context.Context {
	return context.WithValue(c, requestIDKey{}, id)
}

func MethodName(c context.Context) string {
	return c.Value(methodNameKey{}).(string)
}

func WithMethodName(c context.Context, name string) context.Context {
	return context.WithValue(c, methodNameKey{}, name)
}
