package registry

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"plugin"
	"reflect"

	"github.com/Azure/azure-functions-go-worker/azure"
	"github.com/Azure/azure-functions-go-worker/internal/rpc"
	log "github.com/Sirupsen/logrus"
)

var (
	// Funcs contains function id and all information the registry gets from compiled plugin and source code
	Funcs = make(map[string]*azure.Func)
)

// LoadFunc populates information about the func from the compiled plugin and from parsing the source code
func LoadFunc(req *rpc.FunctionLoadRequest) error {
	log.Debugf("received function load request: %v", req)

	f, err := loadFuncFromPlugin(req.Metadata)
	if err != nil {
		return fmt.Errorf("cannot load function from plugin: %v", err)
	}

	namedIn, err := parseEntrypoint(req.Metadata)
	if err != nil {
		return fmt.Errorf("cannot parse entrypoint: %v", err)
	}

	f.Bindings = req.Metadata.Bindings
	f.NamedInArgs = namedIn

	log.Debugf("function: %v", f)
	Funcs[req.FunctionId] = f

	return nil
}

// loadFuncFromPlugin takes the compiled plugin from the func's bin directory
// then reads through reflection the in and out paramns of the entrypoint
func loadFuncFromPlugin(metadata *rpc.RpcFunctionMetadata) (*azure.Func, error) {

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

	in := make([]reflect.Type, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		in[i] = t.In(i)
	}

	out := make([]reflect.Type, t.NumOut())
	for i := 0; i < t.NumOut(); i++ {
		out[i] = t.Out(i)
	}

	return &azure.Func{
		Func: reflect.ValueOf(symbol),
		In:   in,
		Out:  out,
	}, nil
}

func parseEntrypoint(metadata *rpc.RpcFunctionMetadata) ([]*azure.Arg, error) {
	var namedInArgs []*azure.Arg

	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, metadata.ScriptFile, nil, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("cannot parse file %v: %v", metadata.ScriptFile, err)
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
			for _, p := range x.Type.Params.List {
				for _, n := range p.Names {
					// TODO - can any of the values here be nil?
					// TODO - handle cases when in user func there is no pointer type
					key := fmt.Sprintf("*%v.%v", p.Type.(*ast.StarExpr).X.(*ast.SelectorExpr).X.(*ast.Ident).Name, p.Type.(*ast.StarExpr).X.(*ast.SelectorExpr).Sel.Name)
					t, ok := azure.StringToType[key]
					if ok {
						namedInArgs = append(namedInArgs, &azure.Arg{
							Name: n.Name,
							Type: t,
						})
					} else {
						log.Debugf("cannot find key %v in type map", key)
					}
				}
			}

			// this is the entrypoint, no need to traverse the AST any longer
			return false

		default:
			// not a func declaration, need to go further in the AST
			return true
		}
	})

	return namedInArgs, nil
}
