package main

import (
	"github.com/Azure/azure-functions-go-worker/azfunc"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(ctx azfunc.Context, ehMsg *azfunc.EventHubEvent) (ehOut string) {
	ctx.Logger.Log("Log message from function %v, invocation %v, data: %s and msg: %v", ctx.FunctionID, ctx.InvocationID, ehMsg.Data, *ehMsg)

	ehOut = "Hello from Azure Functions"
	return
}
