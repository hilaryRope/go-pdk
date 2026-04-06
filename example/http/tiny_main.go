//go:build !std
// +build !std

package main

import (
	"github.com/extism/go-pdk"
)

//go:wasmexport http_get
func httpGet() int32 {
	// create an HTTP Request (without relying on WASI), set headers as needed
	req := pdk.NewHTTPRequest(pdk.MethodGet, "https://jsonplaceholder.typicode.com/todos/1")
	req.SetHeader("some-name", "some-value")
	req.SetHeader("another", "again")
	// send the request, get response back (can check status on response via res.Status())
	res, err := req.Send()
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	if res.Status() != 200 {
		pdk.SetErrorString("unexpected status code")
		return -1
	}

	// zero-copy output to host
	pdk.OutputMemory(res.Memory())

	return 0
}
