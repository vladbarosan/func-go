# Azure Functions Go Worker

This project aims to add Golang support for Azure Functions.

## How to run the sample

- build the the worker:
  - `docker build -t azure-functions-go-worker-dev .`
- build the sample:
  - `docker run -it -v ${PWD}:/go/src/github.com/Azure/azure-functions-go-worker -w /go/src/github.com/Azure/azure-functions-go-worker golang:1.10 /bin/bash -c "go build -buildmode=plugin -o sample/${SampleName}/${SampleName}.so sample/${SampleName}/main.go"`
- run the worker with the sample
  - `docker run -v ${PWD}/sample/HttpTriggerGo:/home/site/wwwroot/HttpTriggerGo -p 81:80 azure-functions-go-worker-dev`

Then, if you go to `localhost:81/api/HttpTriggerGo`, your `Run` method from the sample should be executed.

Things to notice:

- `entryPoint` - this is the name of the function used as entrypioint

Now let's see the Golang function:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-functions-go-worker/azure"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(req *http.Request, ctx *azure.Context) User {
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
```

Things to notice:

- we can use any vendored dependencies we might have available at compile time (everything is packaged as a Golang plugin)
- the name of the function is `Run` - can be changed, just remember to do the same in `function.json`
- the function signature - `func Run(req *http.Request, ctx *azure.Context) User`. Based on the `function.json` file, `req`, `User`, `outBlob` and `ctx` are automatically populated by the worker.

  > **The content of the parameters is populated based on the name of the parameter! You can change the order, but the name has to be consistent with the name of the binding defined in `function.json`!**

- you can have a return type from the function that, in the case of the `HTTPTrigger` is packaged back as the response body.

## Disclaimer

The project is currently work in progress. Please do not use in production as we expect developments over time.
