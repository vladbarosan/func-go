package main

import (
	"fmt"

	"github.com/Azure/azure-functions-go-worker/azfunc"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(req *azfunc.HTTPRequest, ctx *azfunc.Context) User {
	ctx.Logger.Log("Log message from function %v, invocation %v to the runtime", ctx.FunctionID, ctx.InvocationID)

	u := User{
		Name:          req.Query["name"],
		GeneratedName: fmt.Sprintf("%s-azfunc", req.Query["name"]),
	}

	return u
}

// User mocks any struct (or pointer to struct) you might want to return
type User struct {
	Name          string
	GeneratedName string
}
