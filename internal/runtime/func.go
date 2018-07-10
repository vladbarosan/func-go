package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Azure/azure-functions-go-worker/azfunc"
	"github.com/Azure/azure-functions-go-worker/internal/logger"
	"github.com/Azure/azure-functions-go-worker/internal/rpc"
	"github.com/Azure/azure-functions-go-worker/internal/util"
	log "github.com/Sirupsen/logrus"
)

// Func contains a function symbol with in and out param types
type Func struct {
	Value reflect.Value
	Type  reflect.Type
	In    map[string]*FuncField
	Out   map[string]*FuncField
}

// FuncField represents a representation of a func field
type FuncField struct {
	Name     string
	Type     reflect.Type
	Binding  *rpc.BindingInfo
	Position int
}

func (f *Func) Call(req *rpc.InvocationRequest, eventStream rpc.FunctionRpc_EventStreamClient) ([]*rpc.ParameterBinding, *string, error) {
	params, err := f.bindInputValues(req, eventStream)
	if err != nil {
		log.Debugf("cannot get params from request: %v", err)
		return nil, nil, err
	}
	output := f.Value.Call(params)

	outputData := make([]*rpc.ParameterBinding, len(f.Out))

	for _, v := range f.Out {

		b, err := json.Marshal(output[v.Position].Interface())
		if err != nil {
			log.Debugf("failed to marshal, %v:", err)
		}

		outputData[v.Position] = &rpc.ParameterBinding{
			Name: v.Name,
			Data: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: string(b),
				},
			},
		}
	}

	if len(f.Out) > 0 {
		return outputData, nil, nil
	}

	r := ""

	b, err := json.Marshal(output[0].Interface())
	r = string(b)
	if err != nil {
		log.Debugf("failed to marshal, %v:", err)
	}

	return outputData, &r, nil
}

// bindValues configures the function bindings with the values received
func (f Func) bindInputValues(req *rpc.InvocationRequest, eventStream rpc.FunctionRpc_EventStreamClient) ([]reflect.Value, error) {
	args := make(map[string]reflect.Value)

	// iterate through the invocation request input data
	// if the name of the input data is in the function bindings, then attempt to get the typed binding
	for _, input := range req.InputData {
		param, ok := f.In[input.Name]
		if ok {
			v, err := getInputValue(param, input, req.GetTriggerMetadata())
			if err != nil {
				log.Debugf("cannot transform typed binding: %v", err)
				return nil, err
			}
			args[input.Name] = v
		} else {
			return nil, fmt.Errorf("cannot find input %v in function bindings", input.Name)
		}
	}

	log.Debugf("args map: %v", args)

	params := make([]reflect.Value, len(f.In))
	for _, v := range f.In {
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

func getInputValue(field *FuncField, input *rpc.ParameterBinding, triggetMetadata map[string]*rpc.TypedData) (value reflect.Value, err error) {
	switch azfunc.BindingType(field.Binding.Type) {
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
		return reflect.New(nil), fmt.Errorf("cannot handle binding %v", field.Binding.Type)
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
