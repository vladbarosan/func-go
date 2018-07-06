package main

import (
	"github.com/Azure/azure-functions-go-worker/azfunc"
	log "github.com/Sirupsen/logrus"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(ctx azfunc.Context, in *azfunc.Blob) {
	log.SetLevel(log.DebugLevel)

	log.Debugf("function id: %s, invocation id: %s with blob data: %v", ctx.FunctionID, ctx.InvocationID, in.Data)
}
