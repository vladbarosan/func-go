# Azure Functions Go Worker

This project adds Go support to Azure Functions by implementing a [language
worker][] for Go.

[language worker]: https://github.com/Azure/azure-functions-host/wiki/Language-Extensibility

## Build and Run

### Build and run in a container

- Build the Functions runtime to include the Go worker in this repo:

`docker build -t azure-functions-go-worker .`

- Run your built Functions runtime with a connection to Storage:

`docker run --rm -p 81:80 -e AzureWebJobsStorage="$STORAGE_ACCOUNT_CONN_STRING" azure-functions-go-worker`

### Build and run locally

- Build the worker and the samples: `build.sh`
- Get and install the [functions runtime](https://github.com/Azure/azure-functions-host)
  per instructions in that repo.
- Set environment variables:

```bash
FUNCTIONS_WORKER_RUNTIME=golang
AzureWebJobsScriptRoot=               # path to user functions.
AzureWebJobsStorage=                  # Azure storage account connection string from
                                      # `az storage account show-connection-string`
```

- In `github.com/Azure/azure-functions-host`, modify
  `src/WebJobs.Script.WebHost/appsettings.json` as follows to specify the
  path to the Go worker:

```json
"langaugeWorkers": {
  "workersDirectory":
     "/home/functions-user/go/src/github.com/Azure/azure-functions-go-worker/workers"
}
```

## Test

First set up your user function, then submit an HTTP request to trigger it.

### Write a Go function

Here's how to prepare a Function triggered by an HttpTrigger.

> See the [wiki](https://github.com/Azure/azure-functions-go-worker/wiki) for
> more details on the programming model.

Prepare a `function.json` file alongside your Go code files to specify how to
bind incoming and outgoing event properties to code elements.

```json
{
  "entryPoint": "Run",
  "bindings": [
    {
      "authLevel": "anonymous",
      "type": "httpTrigger",
      "direction": "in",
      "name": "req"
    },
    {
      "name": "$return",
      "type": "http",
      "direction": "out"
    }
  ],
  "disabled": false
}
```

Write the corresponding Go function:

```go
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
func Run(ctx azfunc.Context, req *http.Request) User {

	// additional properties are bound to ctx by Azure Functions
	ctx.Logger.Log("function invoked: function %v, invocation %v",
		ctx.FunctionID, ctx.InvocationID)

	// use standard library to handle incoming request
	body, _ := ioutil.ReadAll(req.Body)

	// deserialize JSON content
	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	// get query param values
	name := req.URL.Query().Get("name")

	u := User{
		Name:          name,
		GeneratedName: fmt.Sprintf("%s-azfunc", name),
		Password:      data["password"].(string),
	}

	return u
}

// User exemplifies a struct to be returned. You can use any struct or *struct.
type User struct {
	Name          string
	GeneratedName string
	Password      string
}
```

### Trigger the function

Use a tool like [Postman](https://www.getpostman.com/apps) to execute a request
with the following parameters:

```
Method: POST
URL: http://localhost:81/api/HttpTrigger?name=helloworld
Body: { "password":"mypassword" }
```

The `Run` method from the sample should be executed.

## Things to note:

- `function.json::entryPoint` names the Go function in package main to be used
  as the Azure Function entry point. In this example that function is named
  `Run` but any name is okay as long as it is also specified in
  `function.json`.
- `main.go` is the required name for the file containing the entry point Go
  function.
- You can use any dependencies you want in your app since they'll be compiled
  into the built binary.
- Structs in the function signature are initialized based on properties in the
  incoming event and specifications in function.json. In the example
  signature of `func Run(ctx azfunc.Context, req *http.Request) User`; `ctx
  azfunc.Context`, `req *http.Request` and `User` are automatically bound to
  incoming and outgoing message properties. Properties received from the GRPC
  channel are bound to properties on the Go structs, and any Go struct
  with the named properties can be used; that is, there's nothing special about
  the default types provided in package azfunc. This is illustrated by the
  returned `User` struct in the example.
- **Properties are bound to parameters based on the name of the parameter! You
  can change the order, but the name has to be consistent with the name of the
  binding defined in `function.json`!**
- You can specify a named return type, which then needs to match an output
  binding in `function.json`. Alternatively, you can have 1 unnamed return type
  which will match the special `$return` binding.

## Disclaimer

* The project is in development; problems and frequent changes are expected.
* This project has not been evaluated for production use.
* We will reply to issues in this repo; but this project is in limited preview
  and not otherwise supported.

## Contributing

This project welcomes contributions and suggestions. Most contributions require
you to agree to a Contributor License Agreement (CLA) declaring that you have
the right to, and actually do, grant us the rights to use your contribution.
For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether
you need to provide a CLA and decorate the PR appropriately (e.g., label,
comment). Simply follow the instructions provided by the bot. You will only
need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of
Conduct](https://opensource.microsoft.com/codeofconduct/).  For more
information see the [Code of Conduct
FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or contact
[opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional
questions or comments.
