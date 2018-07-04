package executor

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Azure/azure-functions-go-worker/azure"
	"github.com/Azure/azure-functions-go-worker/internal/logger"
	"github.com/Azure/azure-functions-go-worker/internal/registry"
	"github.com/Azure/azure-functions-go-worker/internal/rpc"
	"github.com/Azure/azure-functions-go-worker/internal/util"
	log "github.com/Sirupsen/logrus"
)

// ExecuteFunc takes an InvocationRequest and executes the function with corresponding function ID
func ExecuteFunc(req *rpc.InvocationRequest, eventStream rpc.FunctionRpc_EventStreamClient) (response *rpc.InvocationResponse) {

	log.Debugf("\n\n\nInvocation Request: %v", req)

	status := rpc.StatusResult_Success

	f, ok := registry.Funcs[req.FunctionId]
	if !ok {
		log.Debugf("function with functionID %v not loaded", req.FunctionId)
		status = rpc.StatusResult_Failure
	}
	params, outBindings, err := getFinalParams(req, f, eventStream)
	if err != nil {
		log.Debugf("cannot get params from request: %v", err)
		status = rpc.StatusResult_Failure
	}

	log.Debugf("params: %v", params)
	log.Debugf("out bindings: %v", outBindings)

	output := f.Func.Call(params)

	outputData := make([]*rpc.ParameterBinding, len(outBindings))
	i := 0
	for k, v := range outBindings {

		b, err := json.Marshal(v.Interface())
		if err != nil {
			log.Debugf("failed to marshal, %v:", err)
		}

		outputData[i] = &rpc.ParameterBinding{
			Name: k,
			Data: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: string(b),
				},
			},
		}
	}

	r := ""
	if len(output) > 0 {
		b, err := json.Marshal(output[0].Interface())
		r = string(b)
		if err != nil {
			log.Debugf("failed to marshal, %v:", err)
		}
	}

	return &rpc.InvocationResponse{
		InvocationId: req.InvocationId,
		Result: &rpc.StatusResult{
			Status: status,
		},
		ReturnValue: &rpc.TypedData{
			Data: &rpc.TypedData_Json{
				Json: r,
			},
		},
		OutputData: outputData,
	}
}

func getFinalParams(req *rpc.InvocationRequest, f *azure.Func, eventStream rpc.FunctionRpc_EventStreamClient) ([]reflect.Value, map[string]reflect.Value, error) {
	args := make(map[string]reflect.Value)
	outBindings := make(map[string]reflect.Value)

	// iterate through the invocation request input data
	// if the name of the input data is in the function bindings, then attempt to get the typed binding
	for _, input := range req.InputData {
		binding, ok := f.Bindings[input.Name]
		if ok {
			v, err := getValueFromBinding(input, binding)
			if err != nil {
				log.Debugf("cannot transform typed binding: %v", err)
				return nil, nil, err
			}
			args[input.Name] = v
		} else {
			return nil, nil, fmt.Errorf("cannot find input %v in function bindings", input.Name)
		}
	}

	ctx := &azure.Context{
		FunctionID:   req.FunctionId,
		InvocationID: req.InvocationId,
		Logger:       logger.NewLogger(eventStream, req.InvocationId),
	}

	log.Debugf("args map: %v", args)

	params := make([]reflect.Value, len(f.NamedInArgs))
	i := 0
	for _, v := range f.NamedInArgs {
		p, ok := args[v.Name]
		if ok {
			params[i] = p
			i++
		} else if v.Type == reflect.TypeOf((*azure.Context)(nil)) {
			params[i] = reflect.ValueOf(ctx)
			i++
		} else {
			b, ok := f.Bindings[v.Name]
			if ok {
				if b.Direction == rpc.BindingInfo_out {
					o, err := getOutBinding(b)
					if err != nil {
						return nil, nil, fmt.Errorf("cannot get out binding %s: %v", v.Name, err)
					}

					params[i] = o
					outBindings[v.Name] = o
					i++
				}
			}
		}
	}

	return params, outBindings, nil
}

// TODO - add here cases for all bindings supported by Azure Functions
func getValueFromBinding(input *rpc.ParameterBinding, binding *rpc.BindingInfo) (reflect.Value, error) {
	switch azure.BindingType(binding.Type) {
	case azure.HTTPTrigger:
		switch r := input.Data.Data.(type) {
		case *rpc.TypedData_Http:
			h, err := util.ConvertToNativeRequest(r.Http)
			log.Debugf("Converted Http data: %v to: %v", r.Http, *h)

			if err != nil {
				return reflect.New(nil), err
			}
			return reflect.ValueOf(h), nil
		}

	case azure.TimerTrigger:
		switch d := input.Data.Data.(type) {
		case *rpc.TypedData_Json:
			t, err := util.ConvertToTimerInput(d)
			log.Debugf("Converted timer data: %v to: %v", d, *t)

			if err != nil {
				return reflect.New(nil), err
			}

			return reflect.ValueOf(t), nil
		}

	case azure.EventGridTrigger:
		switch d := input.Data.Data.(type) {
		case *rpc.TypedData_Json:
			t, err := util.ConvertToEventGridEvent(d)
			log.Debugf("Converted event grid trigger: %v to: %v", d, *t)

			if err != nil {
				return reflect.New(nil), err
			}

			return reflect.ValueOf(t), nil
		}

	case azure.BlobBinding:
		switch d := input.Data.Data.(type) {
		case *rpc.TypedData_String_:
			b, err := util.ConvertToBlobInput(d)
			if err != nil {
				return reflect.New(nil), err
			}

			return reflect.ValueOf(b), nil
		}
	}
	return reflect.New(nil), fmt.Errorf("cannot handle binding %v", binding.Type)

}

func getOutBinding(b *rpc.BindingInfo) (reflect.Value, error) {
	switch azure.BindingType(b.Type) {
	case azure.BlobBinding:
		b := &azure.Blob{
			Data: "",
		}
		return reflect.ValueOf(b), nil

	case azure.QueueBinding:
		b := &azure.Queue{
			Data: "",
		}
		return reflect.ValueOf(b), nil
	default:
		return reflect.New(nil), fmt.Errorf("cannot handle binding %v", b.Type)
	}
}
