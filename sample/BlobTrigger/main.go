package main

import (
	log "github.com/Sirupsen/logrus"

	"github.com/Azure/azure-functions-go-worker/azure"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(in *azure.Blob, ctx *azure.Context) {
	log.SetLevel(log.DebugLevel)

	log.Debugf("function id: %s, invocation id: %s with blob data: %v", ctx.FunctionID, ctx.InvocationID, in.Data)
}
