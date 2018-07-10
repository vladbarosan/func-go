package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/Azure/azure-functions-go-worker/azfunc"
	"github.com/Azure/azure-functions-go-worker/internal/logger"
	"github.com/Azure/azure-functions-go-worker/internal/rpc"
	log "github.com/Sirupsen/logrus"
)

type fromFn func(*rpc.TypedData, map[string]*rpc.TypedData) (reflect.Value, error)

//converter transforms to and from protobuf into native types.
type converter struct {
	typeBindings map[azfunc.BindingType]fromFn
}

func newConverter() converter {
	c := converter{
		typeBindings: map[azfunc.BindingType]fromFn{},
	}

	c.typeBindings[azfunc.HTTPBinding] = ConvertToRequest
	c.typeBindings[azfunc.HTTPTrigger] = ConvertToRequest
	c.typeBindings[azfunc.BlobTrigger] = ConvertToBlob
	c.typeBindings[azfunc.BlobBinding] = ConvertToBlob
	c.typeBindings[azfunc.QueueTrigger] = ConvertToQueueMsg
	c.typeBindings[azfunc.QueueBinding] = ConvertToQueueMsg
	c.typeBindings[azfunc.TimerTrigger] = ConvertToTimer
	c.typeBindings[azfunc.EventGridTrigger] = ConvertToEventGridEvent

	return c
}

//FromProto converts protobuf parameters to golang values
func (c converter) FromProto(req *rpc.InvocationRequest, eventStream rpc.FunctionRpc_EventStreamClient, f *function) ([]reflect.Value, error) {
	args := make(map[string]reflect.Value)

	// iterate through the invocation request input data
	// if the name of the input data is in the function bindings, then attempt to get the typed binding
	for _, input := range req.InputData {
		param, ok := f.in[input.Name]
		if ok {
			v, ok := c.typeBindings[azfunc.BindingType(param.Binding.GetType())]
			if !ok {
				return nil, fmt.Errorf("cannot handle binding %s", param.Binding.GetType())
			}
			r, err := v(input.GetData(), req.GetTriggerMetadata())
			if err != nil {
				log.Debugf("cannot transform typed binding: %v", err)
				return nil, err
			}
			log.Debugf("Converted  data: %v to: %s", input.Data.Data, r.Interface())

			args[input.Name] = r
		} else {
			return nil, fmt.Errorf("cannot find input %v in function bindings", input.Name)
		}
	}

	log.Debugf("args map: %v", args)

	params := make([]reflect.Value, len(f.in))
	for _, v := range f.in {
		if v.Type == reflect.TypeOf((azfunc.Context{})) {
			ctx := azfunc.Context{
				Context:      context.Background(),
				FunctionID:   req.FunctionId,
				InvocationID: req.InvocationId,
				Logger:       logger.NewLogger(eventStream, req.InvocationId),
			}
			params[v.Position] = reflect.ValueOf(ctx)
		} else {
			params[v.Position] = args[v.Name]
		}
	}

	return params, nil
}

//ToProto converts Values to grpc protocol results
func (c converter) ToProto(values []reflect.Value, fields map[string]*funcField) ([]*rpc.ParameterBinding, *rpc.TypedData, error) {
	protoData := make([]*rpc.ParameterBinding, len(fields))

	for _, v := range fields {

		b, err := json.Marshal(values[v.Position].Interface())
		if err != nil {
			log.Debugf("failed to marshal, %v:", err)
		}

		protoData[v.Position] = &rpc.ParameterBinding{
			Name: v.Name,
			Data: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: string(b),
				},
			},
		}
	}

	// If there are named parameters or no parameters at all there is no return value
	if len(fields) > 0 || len(values) == 0 {
		return protoData, nil, nil
	}

	if len(values) > 2 {
		return nil, nil, fmt.Errorf("Expected 1 or 2 anonymous return values, got %d", len(values))
	}

	ret := ""

	b, err := json.Marshal(values[0].Interface())
	ret = string(b)
	if err != nil {
		log.Debugf("failed to marshal, %v:", err)
	}

	log.Debugf("We have return params and not out params: %s", ret)

	rv := &rpc.TypedData{
		Data: &rpc.TypedData_Json{
			Json: ret,
		},
	}
	return protoData, rv, nil
}

// ConvertToRequest returns a native http.Request from an rpc.HttpTrigger
func ConvertToRequest(d *rpc.TypedData, tm map[string]*rpc.TypedData) (reflect.Value, error) {

	t, ok := d.Data.(*rpc.TypedData_Http)

	if !ok {
		return reflect.Value{}, fmt.Errorf("cannot convert non http Http request")
	}

	if t.Http == nil {
		return reflect.Value{}, fmt.Errorf("cannot convert nil request")
	}

	var body io.Reader
	if t.Http.RawBody != nil {
		switch d := t.Http.RawBody.Data.(type) {
		case *rpc.TypedData_String_:
			body = ioutil.NopCloser(bytes.NewBufferString(d.String_))
		}
	}

	req, err := http.NewRequest(t.Http.GetMethod(), t.Http.GetUrl(), body)

	if err != nil {
		return reflect.Value{}, err
	}

	for key, value := range t.Http.GetHeaders() {
		req.Header.Set(key, value)
	}

	return reflect.ValueOf(req), nil
}

// ConvertToBlob returns a formatted Blob from an rpc.TypedData_String (as blob inputs are for now)
func ConvertToBlob(d *rpc.TypedData, tm map[string]*rpc.TypedData) (reflect.Value, error) {

	t, ok := d.Data.(*rpc.TypedData_String_)

	if !ok {
		return reflect.Value{}, fmt.Errorf("cannot convert blob input")
	}

	b := &azfunc.Blob{
		Data: t.String_,
	}

	return reflect.ValueOf(b), nil
}

// ConvertToQueueMsg returns a formatted Queue from an rpc.TypedData_String
func ConvertToQueueMsg(d *rpc.TypedData, tm map[string]*rpc.TypedData) (reflect.Value, error) {
	t, ok := d.Data.(*rpc.TypedData_String_)

	if !ok {
		return reflect.Value{}, fmt.Errorf("cannot convert queue message input")
	}

	qm := &azfunc.QueueMsg{
		Text: t.String_,
	}

	return reflect.ValueOf(qm), nil
}

//ConvertToTimer returns a formatted TimerInput from an rpc.
func ConvertToTimer(d *rpc.TypedData, tm map[string]*rpc.TypedData) (reflect.Value, error) {

	t, ok := d.Data.(*rpc.TypedData_Json)

	if !ok {
		return reflect.Value{}, fmt.Errorf("cannot convert non json timer")
	}

	timer := &azfunc.Timer{}
	if err := json.Unmarshal([]byte(t.Json), &timer); err != nil {
		return reflect.Value{}, fmt.Errorf("cannot unmarshal timer object")
	}

	return reflect.ValueOf(timer), nil
}

// ConvertToMap returns a map object
func ConvertToMap(d *rpc.TypedData, tm map[string]*rpc.TypedData) (reflect.Value, error) {

	t, ok := d.Data.(*rpc.TypedData_Json)

	if !ok {
		return reflect.Value{}, fmt.Errorf("cannot convert non json map input")
	}

	m := map[string]interface{}{}
	if err := json.Unmarshal([]byte(t.Json), &m); err != nil {
		return reflect.Value{}, fmt.Errorf("cannot unmarshal map")
	}

	return reflect.ValueOf(m), nil
}

// ConvertToEventGridEvent returns an EventGridEvent
func ConvertToEventGridEvent(d *rpc.TypedData, tm map[string]*rpc.TypedData) (reflect.Value, error) {

	t, ok := d.Data.(*rpc.TypedData_Json)

	if !ok {
		return reflect.Value{}, fmt.Errorf("cannot convert non json event grid event input")
	}

	e := &azfunc.EventGridEvent{}
	if err := json.Unmarshal([]byte(t.Json), &e); err != nil {
		return reflect.Value{}, fmt.Errorf("cannot unmarshal event grid event object")
	}

	return reflect.ValueOf(e), nil
}
