package main

import "github.com/Azure/azure-functions-go-worker/azfunc"

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(ctx azfunc.Context, timer *azfunc.Timer) {
	ctx.Logger.Log("Log message from function %v, invocation %v to the runtime with timer: %v", ctx.FunctionID, ctx.InvocationID, *timer)
}
