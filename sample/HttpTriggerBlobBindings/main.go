package main

import (
	"net/http"

	"github.com/Azure/azure-functions-go-worker/azfunc"
	log "github.com/Sirupsen/logrus"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(ctx azfunc.Context, req *http.Request, inBlob *azfunc.Blob, outBlob *azfunc.Blob) BlobData {
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
