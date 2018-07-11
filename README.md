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
	"net/http"

	"github.com/Azure/azure-functions-go-worker/azfunc"
)

// Run is the entrypoint to our Go Azure Function - if you want to change it, see function.json
func Run(ctx azfunc.Context, req *http.Request, inBlob *azfunc.Blob) (outBlob string) {

	ctx.Logger.Log("function id: %s, invocation id: %s with blob : %v", ctx.FunctionID, ctx.InvocationID, *inBlob)

	outBlob = inBlob.Content
	return
}
```

Things to notice:

- we can use any vendored dependencies we might have available at compile time (everything is packaged as a Golang plugin)
- the name of the function is `Run` - can be changed, just remember to do the same in `function.json`
- the function signature - `func Run(req *http.Request, ctx *azfunc.Context) User`. Based on the `function.json` file, `ctx`, `req`, and `inBlob` are automatically populated by the worker.

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
s
