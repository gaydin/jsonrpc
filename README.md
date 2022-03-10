# jsonrpc 2.0

## About

JSON-RPC 2.0 implementation with generics

## Install

```
$ go get -u github.com/gaydin/jsonrpc/v2
```

## Usage

```go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gaydin/jsonrpc"
)

type (
	Params struct {
		Name string `json:"name"`
	}

	Result struct {
		Message string `json:"message"`
	}
)

func NewHandler() jsonrpc.Method {
	return jsonrpc.NewParamsMethod("user.message", handler)
}

func handler(ctx context.Context, p Params) (*Result, *jsonrpc.Error) {
	return &Result{
		Message: "Hello, " + p.Name,
	}, nil
}

func main() {
	rpc := jsonrpc.New()

	if err := rpc.RegisterMethod(NewHandler()); err != nil {
		log.Fatalln(err)
	}

	http.Handle("/jrpc", rpc)

	if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
		log.Fatalln(err)
	}
}
```

## Based on

https://github.com/osamingo/jsonrpc