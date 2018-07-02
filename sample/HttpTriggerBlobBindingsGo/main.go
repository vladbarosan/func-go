package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/Azure/azure-functions-go-worker/azure"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(req *http.Request, inBlob *azure.Blob, outBlob *azure.Blob, ctx *azure.Context) BlobData {
	log.SetLevel(log.DebugLevel)

	log.Debugf("function id: %s, invocation id: %s", ctx.FunctionID, ctx.InvocationID)

	name := req.URL.Query().Get("name")

	d := BlobData{
		Name: name,
		Data: inBlob.Data,
	}

	outBlob.Data = inBlob.Data

	return d
}

// BlobData mocks any struct (or pointer to struct) you might want to return
type BlobData struct {
	Name string
	Data string
}
