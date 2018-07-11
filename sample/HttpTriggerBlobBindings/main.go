package main

import (
	"net/http"

	"github.com/Azure/azure-functions-go-worker/azfunc"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(ctx azfunc.Context, req *http.Request, inBlob *azfunc.Blob) (outBlob string) {

	ctx.Logger.Log("function id: %s, invocation id: %s with blob : %v", ctx.FunctionID, ctx.InvocationID, *inBlob)

	outBlob = inBlob.Content
	return
}
