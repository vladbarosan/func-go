package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-functions-go-worker/azure"
	"github.com/Azure/azure-functions-go-worker/internal/rpc"
)

// ConvertToNativeRequest returns a native http.Request from an rpc.HttpTrigger
func ConvertToNativeRequest(r *rpc.RpcHttp) (*http.Request, error) {

	if r == nil {
		return nil, fmt.Errorf("cannot convert nil request")
	}

	var body io.Reader
	if r.RawBody != nil {
		switch d := r.RawBody.Data.(type) {
		case *rpc.TypedData_String_:
			body = ioutil.NopCloser(bytes.NewBufferString(d.String_))
		}
	}

	req, err := http.NewRequest(r.GetMethod(), r.GetUrl(), body)

	if err != nil {
		return nil, err
	}

	for key, value := range r.GetHeaders() {
		req.Header.Set(key, value)
	}
	return req, nil
}

// ConvertToBlobInput returns a formatted BlobInput from an rpc.TypedData_String (as blob inputs are for now)
func ConvertToBlobInput(d *rpc.TypedData_String_) (*azure.Blob, error) {
	if d == nil {
		return nil, fmt.Errorf("cannot convert nil blob input")
	}

	return &azure.Blob{
		Data: d.String_,
	}, nil
}
