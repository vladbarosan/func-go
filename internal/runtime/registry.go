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
	funcs     map[string]*function
	converter converter
}

// NewRegistry returns a new function registry
func NewRegistry() *Registry {
	return &Registry{
		funcs:     map[string]*function{},
		converter: newConverter(),
	}
}

// LoadFunc populates information about the func from the compiled plugin and from parsing the source code
func (r Registry) LoadFunc(req *rpc.FunctionLoadRequest) error {
	log.Debugf("received function load request: %v", req)

	f, err := loadFuncFromPlugin(req.Metadata)
	if err != nil {
		return fmt.Errorf("cannot load function from plugin: %v", err)
	}

	ins, outs, err := loadInOut(req.Metadata, f.signature)
	if err != nil {
		return fmt.Errorf("cannot parse entrypoint: %v", err)
	}

	f.in = ins
	f.out = outs

	log.Debugf("function: %v", f)
	r.funcs[req.FunctionId] = f

	return nil
}

// ExecuteFunc takes an InvocationRequest and executes the function with corresponding function ID
func (r Registry) ExecuteFunc(req *rpc.InvocationRequest, eventStream rpc.FunctionRpc_EventStreamClient) (response *rpc.InvocationResponse) {

	log.Debugf("\n\n\nInvocation Request: %v", req)

	status := rpc.StatusResult_Success

	ir := &rpc.InvocationResponse{
		InvocationId: req.InvocationId,
		Result: &rpc.StatusResult{
			Status: status,
		},
	}
	f, ok := r.funcs[req.FunctionId]

	if !ok {
		log.Debugf("function with functionID %v not loaded", req.FunctionId)
		ir.Result.Status = rpc.StatusResult_Failure
		return ir
	}

	params, err := r.converter.FromProto(req, eventStream, f)
	if err != nil {
		ir.Result.Status = rpc.StatusResult_Failure
		return ir
	}
	output, err := f.Invoke(params)
	if err != nil {
		ir.Result.Status = rpc.StatusResult_Failure
		return ir
	}
	o, rv, err := r.converter.ToProto(output, f.out)

	if err != nil {
		log.Debugf("cannot get output data from result %v", err)
		if err != nil {
			ir.Result.Status = rpc.StatusResult_Failure
			return ir
		}
	}

	ir.ReturnValue = rv
	ir.OutputData = o
	return ir
}

// loadFuncFromPlugin takes the compiled plugin from the func's bin directory
// then reads through reflection the in and out paramns of the entrypoint
func loadFuncFromPlugin(metadata *rpc.RpcFunctionMetadata) (*function, error) {

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

	return &function{
		handler:   reflect.ValueOf(symbol),
		signature: t,
	}, nil
}

type iterator func(int) reflect.Type

// loadInOut loads the input and output types for a function
func loadInOut(metadata *rpc.RpcFunctionMetadata, funcType reflect.Type) (map[string]*funcField, map[string]*funcField, error) {

	var ins, outs map[string]*funcField
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

func extractFuncFields(fl *ast.FieldList, bindings map[string]*rpc.BindingInfo, fi iterator, l int) (map[string]*funcField, error) {
	fields := map[string]*funcField{}

	if fl.NumFields() != l {
		return nil, fmt.Errorf("Plugin %d and source %d nr of arguments are different", fl.NumFields(), l)
	}

	if l == 0 {
		return fields, nil
	}

	for i, p := range fl.List {
		t := fi(i)
		for _, n := range p.Names {
			log.Debugf("Found parameter: %s with type: %s", n, t.String())

			fields[n.Name] = &funcField{
				Name:     n.Name,
				Type:     fi(i),
				Position: i,
				Binding:  bindings[n.Name],
			}
		}
	}

	return fields, nil
}
