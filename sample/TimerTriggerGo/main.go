package main

import (
	"github.com/Azure/azure-functions-go-worker/azure"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(timer *azure.Timer, ctx *azure.Context) {
	ctx.Logger.Log("Log message from function %v, invocation %v to the runtime is passed due: %t", ctx.FunctionID, ctx.InvocationID, timer.PastDue)
}
