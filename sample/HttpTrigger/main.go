package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-functions-go-worker/azure"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(req *http.Request, ctx *azure.Context) User {
	ctx.Logger.Log("Log message from function %v, invocation %v to the runtime", ctx.FunctionID, ctx.InvocationID)
	body, _ := ioutil.ReadAll(req.Body)
	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	name := req.URL.Query().Get("name")
	u := User{
		Name:          name,
		GeneratedName: fmt.Sprintf("%s-azfunc", name),
		Password:      data["password"].(string),
	}

	return u
}

// User mocks any struct (or pointer to struct) you might want to return
type User struct {
	Name          string
	GeneratedName string
	Password      string
}
