package main

import (
	log "github.com/Sirupsen/logrus"

	"github.com/Azure/azure-functions-go-worker/azfunc"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(ctx *azfunc.Context, queueMsg *azfunc.QueueMsg) {
	log.SetLevel(log.DebugLevel)
	log.Debugf("function id: %s, invocation id: %s with queue data: %v", ctx.FunctionID, ctx.InvocationID, queueMsg.Data)
	ctx.Logger.Log("function id: %s, invocation id: %s with queue data: %v", ctx.FunctionID, ctx.InvocationID, queueMsg.Data)
}
