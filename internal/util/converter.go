package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-functions-go-worker/azfunc"
	"github.com/Azure/azure-functions-go-worker/internal/rpc"
)

// ConvertToNativeRequest returns a native http.Request from an rpc.HttpTrigger
func ConvertToNativeRequest(d *rpc.TypedData) (*http.Request, error) {

	t, ok := d.Data.(*rpc.TypedData_Http)

	if !ok {
		return nil, fmt.Errorf("cannot convert non http Http request")
	}

	r := t.Http
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

// ConvertToBlob returns a formatted Blob from an rpc.TypedData_String (as blob inputs are for now)
func ConvertToBlob(d *rpc.TypedData) (*azfunc.Blob, error) {

	t, ok := d.Data.(*rpc.TypedData_String_)

	if !ok {
		return nil, fmt.Errorf("cannot convert blob input")
	}

	return &azfunc.Blob{
		Data: t.String_,
	}, nil
}

// ConvertToQueueMsg returns a formatted Queue from an rpc.TypedData_String
func ConvertToQueueMsg(d *rpc.TypedData) (*azfunc.QueueMsg, error) {
	t, ok := d.Data.(*rpc.TypedData_String_)

	if !ok {
		return nil, fmt.Errorf("cannot convert queue message input")
	}

	return &azfunc.QueueMsg{
		Data: t.String_,
	}, nil
}

//ConvertToTimer returns a formatted TimerInput from an rpc.
func ConvertToTimer(d *rpc.TypedData) (*azfunc.Timer, error) {

	t, ok := d.Data.(*rpc.TypedData_Json)

	if !ok {
		return nil, fmt.Errorf("cannot convert non json timer")
	}

	timer := azfunc.Timer{}
	if err := json.Unmarshal([]byte(t.Json), &timer); err != nil {
		return nil, fmt.Errorf("cannot unmarshal timer object")
	}

	return &timer, nil
}

// ConvertToMap returns a map object
func ConvertToMap(d *rpc.TypedData) (map[string]interface{}, error) {

	t, ok := d.Data.(*rpc.TypedData_Json)

	if !ok {
		return nil, fmt.Errorf("cannot convert non json map input")
	}

	m := map[string]interface{}{}
	if err := json.Unmarshal([]byte(t.Json), &m); err != nil {
		return nil, fmt.Errorf("cannot unmarshal map")
	}

	return m, nil
}

// ConvertToEventGridEvent returns an EventGridEvent
func ConvertToEventGridEvent(d *rpc.TypedData) (*azfunc.EventGridEvent, error) {

	t, ok := d.Data.(*rpc.TypedData_Json)

	if !ok {
		return nil, fmt.Errorf("cannot convert non json event grid event input")
	}

	eventGridEvent := azfunc.EventGridEvent{}
	if err := json.Unmarshal([]byte(t.Json), &eventGridEvent); err != nil {
		return nil, fmt.Errorf("cannot unmarshal event grid event object")
	}

	return &eventGridEvent, nil
}
