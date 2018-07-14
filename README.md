# Azure Functions Go Worker

This project aims to add Golang support for Azure Functions.

## How to run the Go Functions Worker

### Running containerized runtime with the worker

- Build the docker image
  - `docker build --rm -f Dockerfile -t azure-functions-go-worker .`
- Run the container
  - `docker run --rm -p 81:80 -e AzureWebJobsStorage="$STORAGE_ACCOUNT_CONN_STRING" azure-functions-go-worker`

### Running locally

- Build the worker and the samples - `build.sh`
- Get the [functions runtime](https://github.com/Azure/azure-functions-host) from and follow the instructions there for setting it up

  - Set the variables:

    - `AzureWebJobsScriptRoot` - needs to point to where the function app is
    - `FUNCTIONS_WORKER_RUNTIME` - needs to be `"golang"`
    - `AzureWebJobsStorage` - needs to be your Azure Storage Account Connection String
    - In the `appsettings.json` file of the runtime set the worker directory:

```json
"langaugeWorkers": {
  "workersDirectory":
     "/home/vladdb/go/src/github.com/Azure/azure-functions-go-worker/workers"
}
```

Then, using a [REST API Client](https://www.getpostman.com/apps), if you execute a POST call to go to `localhost:81/api/HttpTrigger?name=vladdb` and with a body property called `password`, your `Run` method from the `HttpTrigger` should be executed.

## Sample Go Function

- See the [Wiki](https://github.com/Azure/azure-functions-go-worker/wiki) for mode details on the programming model.

Taking the simplest of the samples for an Http Trigger, The `function.json` looks like:

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

Now let's see the Golang function:

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-functions-go-worker/azfunc"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(ctx azfunc.Context, req *http.Request) User {
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
```

Things to notice:

- `entryPoint` - this is the name of the function used as entry point. This needs to match the function name.
- `main.go` - this needs to be the name of the file containing the entry point.
- we can use any vendored dependencies we might have available at compile time (everything is packaged as a Golang plugin).
- the name of the function is `Run` - can be changed, just remember to do the same in `function.json`
- the function signature - `func Run(ctx azfunc.Context, req *http.Request) User`. Based on the `function.json` file, `ctx`, `req`, and `inBlob` are automatically populated by the worker. Note that there is no implicit dependency on the types provided by the worker in the `azfunc` package! ( Does can be copied and modified or anything else that matches the GRPC protocol can be used.)

  > **The content of the parameters is populated based on the name of the parameter! You can change the order, but the name has to be consistent with the name of the binding defined in `function.json`!**

- you can have a named returned type that needs to match an output binding in `function.json`. You can also have 1 anonymous return value that will match the `$return` binding.

## Disclaimer

The project is currently work in progress. Please do not use in production as we expect developments over time.

## Contributing

This project welcomes contributions and suggestions. Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
