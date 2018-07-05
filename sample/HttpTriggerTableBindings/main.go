package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/Azure/azure-functions-go-worker/azure"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(req *http.Request, in map[string]interface{}, out map[string]interface{}, ctx *azure.Context) {
	log.SetLevel(log.DebugLevel)

	log.Debugf("function id: %s, invocation id: %s with person name: %v", ctx.FunctionID, ctx.InvocationID, in["name"])
	out["name"] = "new name"
	out["RowKey"] = "newTestKey"
}
