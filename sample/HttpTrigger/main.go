package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-functions-go-worker/azfunc"
)

// Run runs this Azure Function because it is specified in `function.json` as
// the entryPoint. Fields of the function's parameters are also bound to
// incoming and outgoing event properties as specified in `function.json`.
func Run(ctx azfunc.Context, req *http.Request) (*User, error) {

	// additional properties are bound to ctx by Azure Functions
	ctx.Log(azfunc.LogInformation, "function invoked: functionID %v, invocationID %v", ctx.FunctionID(), ctx.InvocationID())

	// use standard library to handle incoming request
	body, _ := ioutil.ReadAll(req.Body)

	// deserialize JSON content
	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	// get query param values
	name := req.URL.Query().Get("name")

	if name == "" {
		return nil, fmt.Errorf("Missing required parameter: name")
	}

	u := &User{
		Name:          name,
		GeneratedName: fmt.Sprintf("%s-azfunc", name),
		Password:      data["password"].(string),
	}

	return u, nil
}

// User exemplifies a struct to be returned. You can use any struct or *struct.
type User struct {
	Name          string
	GeneratedName string
	Password      string
}
