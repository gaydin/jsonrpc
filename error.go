package jsonrpc

import "fmt"

const (
	ErrorCodeParse          ErrorCode = -32700
	ErrorCodeInvalidRequest ErrorCode = -32600
	ErrorCodeMethodNotFound ErrorCode = -32601
	ErrorCodeInvalidParams  ErrorCode = -32602
	ErrorCodeInternal       ErrorCode = -32603
)

type (
	ErrorCode int

	Error struct {
		Code    ErrorCode   `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}
)

func (e *Error) Error() string {
	return fmt.Sprintf("jsonrpc: code: %d, message: %s, data: %+v", e.Code, e.Message, e.Data)
}

func ErrParse() *Error {
	return &Error{
		Code:    ErrorCodeParse,
		Message: "Parse error",
	}
}

func ErrInvalidRequest() *Error {
	return &Error{
		Code:    ErrorCodeInvalidRequest,
		Message: "Invalid Request",
	}
}

func ErrMethodNotFound() *Error {
	return &Error{
		Code:    ErrorCodeMethodNotFound,
		Message: "Method not found",
	}
}

func ErrInvalidParams() *Error {
	return &Error{
		Code:    ErrorCodeInvalidParams,
		Message: "Invalid params",
	}
}

func ErrInternal() *Error {
	return &Error{
		Code:    ErrorCodeInternal,
		Message: "Internal error",
	}
}
