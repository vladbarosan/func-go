package main

import (
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-functions-go-worker/azfunc"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(req *http.Request, outQueue *azfunc.QueueMsg, ctx *azfunc.Context) {
	ctx.Logger.Log("Log message from function %v, invocation %v to the runtime", ctx.FunctionID, ctx.InvocationID)

	body, _ := ioutil.ReadAll(req.Body)
	outQueue.Data = string(body)
}
