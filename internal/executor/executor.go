package executor

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Azure/azure-functions-go-worker/azfunc"
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

func getFinalParams(req *rpc.InvocationRequest, f *azfunc.Func, eventStream rpc.FunctionRpc_EventStreamClient) ([]reflect.Value, map[string]reflect.Value, error) {
	args := make(map[string]reflect.Value)
	outBindings := make(map[string]reflect.Value)

	// iterate through the invocation request input data
	// if the name of the input data is in the function bindings, then attempt to get the typed binding
	for _, input := range req.InputData {
		binding, ok := f.Bindings[input.Name]
		if ok {
			v, err := getInputValue(input, binding, req.GetTriggerMetadata())
			if err != nil {
				log.Debugf("cannot transform typed binding: %v", err)
				return nil, nil, err
			}
			args[input.Name] = v
		} else {
			return nil, nil, fmt.Errorf("cannot find input %v in function bindings", input.Name)
		}
	}

	ctx := &azfunc.Context{
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
		} else if v.Type == reflect.TypeOf((*azfunc.Context)(nil)) {
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

func getInputValue(input *rpc.ParameterBinding, binding *rpc.BindingInfo, triggetMetadata map[string]*rpc.TypedData) (value reflect.Value, err error) {
	switch azfunc.BindingType(binding.Type) {
	case azfunc.HTTPTrigger:
		h, err := util.ConvertToNativeRequest(input.GetData())
		log.Debugf("Converted Http data: %v to: %v", input.Data.Data, *h)

		if err != nil {
			return reflect.New(nil), err
		}
		return reflect.ValueOf(h), nil

	case azfunc.TimerTrigger:
		t, err := util.ConvertToTimer(input.GetData())
		log.Debugf("Converted timer data: %v to: %v", input.GetData().GetData(), *t)

		if err != nil {
			return reflect.New(nil), err
		}

		return reflect.ValueOf(t), nil

	case azfunc.EventGridTrigger:
		t, err := util.ConvertToEventGridEvent(input.GetData())
		log.Debugf("Converted event grid trigger: %v to: %v", input.GetData().GetData(), *t)

		if err != nil {
			return reflect.New(nil), err
		}

		return reflect.ValueOf(t), nil

	case azfunc.BlobTrigger:
		fallthrough
	case azfunc.BlobBinding:
		b, err := util.ConvertToBlob(input.GetData())
		log.Debugf("Converted blob binding: %v to: %v", input.GetData().GetData(), b)
		if err != nil {
			return reflect.New(nil), err
		}

		return reflect.ValueOf(b), nil
	case azfunc.QueueTrigger:
		qm, err := util.ConvertToQueueMsg(input.GetData())
		log.Debugf("Converted queue mesg binding: %v to: %v", input.GetData().GetData(), qm)
		if err != nil {
			return reflect.New(nil), err
		}

		return reflect.ValueOf(qm), nil
	case azfunc.TableBinding:
		m, err := util.ConvertToMap(input.GetData())
		log.Debugf("Converted table binding: %v to: %v", input.GetData().GetData(), m)

		if err != nil {
			return reflect.New(nil), err
		}

		return reflect.ValueOf(m), nil
	default:
		return reflect.New(nil), fmt.Errorf("cannot handle binding %v", binding.Type)
	}
}

func getOutBinding(b *rpc.BindingInfo) (reflect.Value, error) {
	switch azfunc.BindingType(b.Type) {
	case azfunc.BlobBinding:
		d := &azfunc.Blob{}
		return reflect.ValueOf(d), nil
	case azfunc.QueueBinding:
		d := &azfunc.QueueMsg{}
		return reflect.ValueOf(d), nil
	case azfunc.TableBinding:
		d := map[string]interface{}{}
		return reflect.ValueOf(d), nil
	default:
		return reflect.New(nil), fmt.Errorf("cannot handle binding %v", b.Type)
	}
}
