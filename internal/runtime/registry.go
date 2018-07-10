package runtime

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"plugin"
	"reflect"

	"github.com/Azure/azure-functions-go-worker/internal/rpc"
	log "github.com/Sirupsen/logrus"
)

// Registry contains all information about user functions and how to execute them
type Registry struct {
	Funcs map[string]*Func
}

// NewRegistry returns a new function registry
func NewRegistry() *Registry {
	return &Registry{
		Funcs: map[string]*Func{},
	}
}

// LoadFunc populates information about the func from the compiled plugin and from parsing the source code
func (r Registry) LoadFunc(req *rpc.FunctionLoadRequest) error {
	log.Debugf("received function load request: %v", req)

	f, err := loadFuncFromPlugin(req.Metadata)
	if err != nil {
		return fmt.Errorf("cannot load function from plugin: %v", err)
	}

	ins, outs, err := loadInOut(req.Metadata, f.Type)
	if err != nil {
		return fmt.Errorf("cannot parse entrypoint: %v", err)
	}

	f.In = ins
	f.Out = outs

	log.Debugf("function: %v", f)
	r.Funcs[req.FunctionId] = f

	return nil
}

// ExecuteFunc takes an InvocationRequest and executes the function with corresponding function ID
func (r Registry) ExecuteFunc(req *rpc.InvocationRequest, eventStream rpc.FunctionRpc_EventStreamClient) (response *rpc.InvocationResponse) {

	log.Debugf("\n\n\nInvocation Request: %v", req)

	status := rpc.StatusResult_Success

	f, ok := r.Funcs[req.FunctionId]

	if !ok {
		log.Debugf("function with functionID %v not loaded", req.FunctionId)
		status = rpc.StatusResult_Failure
	}

	out, ret, err := f.Call(req, eventStream)

	if err != nil {
		log.Debugf("cannot get params from request: %v", err)
		status = rpc.StatusResult_Failure
	}

	var rv *rpc.TypedData

	if ret != nil {
		log.Debugf("We have return params and not out params: %s", *ret)

		rv = &rpc.TypedData{
			Data: &rpc.TypedData_Json{
				Json: *ret,
			},
		}
	}

	return &rpc.InvocationResponse{
		InvocationId: req.InvocationId,
		Result: &rpc.StatusResult{
			Status: status,
		},
		ReturnValue: rv,
		OutputData:  out,
	}
}

// loadFuncFromPlugin takes the compiled plugin from the func's bin directory
// then reads through reflection the in and out paramns of the entrypoint
func loadFuncFromPlugin(metadata *rpc.RpcFunctionMetadata) (*Func, error) {

	path := fmt.Sprintf("%s/bin/%s.so", metadata.Directory, metadata.Name)
	plugin, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot get .so object from path %s: %v", path, err)
	}

	symbol, err := plugin.Lookup(metadata.EntryPoint)
	if err != nil {
		return nil, fmt.Errorf("cannot look up symbol for entrypoint function %s: %v", metadata.EntryPoint, err)
	}

	t := reflect.TypeOf(symbol)
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("symbol is not func, but %v", t.Kind())
	}

	return &Func{
		Value: reflect.ValueOf(symbol),
		Type:  t,
	}, nil
}

type iterator func(int) reflect.Type

// loadInOut loads the input and output types for a function
func loadInOut(metadata *rpc.RpcFunctionMetadata, funcType reflect.Type) (map[string]*FuncField, map[string]*FuncField, error) {

	var ins, outs map[string]*FuncField
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, metadata.ScriptFile, nil, parser.AllErrors)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot parse file %v: %v", metadata.ScriptFile, err)
	}

	// traverse the AST and inspect the nodes
	// if the node is a func declaration, check if entrypoint and get input params names and types (as string)
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			log.Debugf("found function: %v", x.Name.Name)
			if x.Name.Name != metadata.EntryPoint {
				log.Debugf("not function entrypoint, moving on...")

				// not the entrypoint, go further into the AST
				return true
			}

			ins, err = extractFuncFields(x.Type.Params, metadata.GetBindings(), funcType.In, funcType.NumIn())
			outs, err = extractFuncFields(x.Type.Results, metadata.GetBindings(), funcType.Out, funcType.NumOut())

			// this is the entrypoint, no need to traverse the AST any longer
			return false

		default:
			// not a func declaration, need to go further in the AST
			return true
		}
	})

	return ins, outs, nil
}

func extractFuncFields(fl *ast.FieldList, bindings map[string]*rpc.BindingInfo, fi iterator, l int) (map[string]*FuncField, error) {
	fields := map[string]*FuncField{}

	if fl.NumFields() != l {
		return nil, fmt.Errorf("Plugin %d and source %d nr of arguments are different", fl.NumFields(), l)
	}

	for i, p := range fl.List {
		t := fi(i)
		for _, n := range p.Names {
			log.Debugf("Found parameter: %s with type: %s", n, t.String())

			fields[n.Name] = &FuncField{
				Name:     n.Name,
				Type:     fi(i),
				Position: i,
				Binding:  bindings[n.Name],
			}
		}
	}

	return fields, nil
}
