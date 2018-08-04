# Azure Functions for Go

[![Build Status](https://travis-ci.com/vladbarosan/func-go.svg?token=pzfiiBDjqjzLCtQCMpq1&branch=dev)](https://travis-ci.com/vladbarosan/func-go)

This project adds Go support to Azure Functions by implementing a [language
worker][] for Go. Note that only Go 1.10+ is supported. Supported platforms are Linux and Mac.

[language worker]: https://github.com/Azure/azure-functions-host/wiki/Language-Extensibility

## Contents:

- [Run a Go Functions instance](#run-a-go-functions-instance)
- [Write and deploy a Go Function](#write-and-deploy-a-go-function)

# Run a Go Functions instance

Clone the repo and run one of the following to deploy an instance with
batteries (i.e. prebuilt samples) included to:

- your friendly local Docker daemon: `make local-instance`
- Azure App Service: `make azure-instance`

**NOTE** that to use Azure App Service you must specify a public registry for
`RUNTIME_IMAGE_REGISTRY` in `.env` rather than `local`.

The local-instance creates and utilizes an Azure Storage account. The
azure-instance also creates and utilizes an App Service plan and functionapp.

Some triggers are triggered after the instance is created. (TODO: make them
skippable; and verify success.)

Default instance configuration options can be overriden by maintaining your own
`.env` file in the root of your clone. When no `.env` file is found one is
created based on `.env.tpl`.

### Other Options

- Build the Functions runtime including the Go worker in this repo:
  `test/build.sh`. Add `1` as a parameter to also push to a registry. The image
  name is built from configuration in `.env`.

- Run an instance of the runtime: `docker run --rm --publish 8080:80 --env "AzureWebJobsStorage=$STORAGE_CONNSTR" local/azure-functions-go-worker:dev`.
  (The image name chosen reflects the defaults in `.env.tpl`.)

To discover your Storage account connection string consider using `az storage account show-connection-string ...`.

To run **Event Hubs** samples, set a namespace connection string in an environment
variable, and specify that environment variable name as the value of the
`connection` field in the functionapp's `function.json`. The variable name used
in the samples is "EventHubConnectionSetting". This applies for **Service Bus** and **CosmosDb** samples.

To discover the connection strings required you can use:

- for **Event Hubs**: `az eventhubs namespace authorization-rule list ...`.
- for **CosmosDb**: `az cosmosdb list-keys` and then the format `connstr="AccountEndpoint=https://$account_name.documents.azure.com:443/;AccountKey=$account_key;"`
- for **Service Bus**: `az servicebus namespace authorization-rule keys list ...`

### Run locally without containers

- Build the worker and the samples: `build.sh native bundle`
- Get and install the [functions runtime](https://github.com/Azure/azure-functions-host)
  per instructions in that repo.
- Set environment variables:

```bash
FUNCTIONS_WORKER_RUNTIME=golang              # intended target language worker
AzureWebJobsScriptRoot=/home/site/wwwroot    # path in container fs to user code
AzureWebJobsStorage=                         # Storage account connection string
EventHubConnectionSetting=                   # Event Hubs namespace connection string
ServiceBusConnectionString=                  # Service Bus namespace SAS Policy ConnectionString
CosmosDBConnectionString=                     # Cosmos Db account connection string
```

- In `github.com/Azure/azure-functions-host`, modify
  `src/WebJobs.Script.WebHost/appsettings.json` as follows to specify the
  path to the Go worker:

```json
"langaugeWorkers": {
  "workersDirectory":
     "/home/functions-user/go/src/github.com/vladbarosan/func-go/workers"
}
```

# Write and deploy a Go Function

Follow these high-level steps to create Go Functions:

1.  Write a Go Function.
2.  Deploy it.
3.  Trigger and watch it.

Following are step-by-step instructions to prepare a Go Function triggered by
an HttpTrigger, as demonstrated in [the HttpTrigger sample][].

> See [the wiki][] and [Things to Note](#things-to-note) below for more details
> on the programming model.

[the httptrigger sample]: ./sample/HttpTrigger
[the wiki]: https://github.com/vladbarosan/func-go/wiki/Programming-Model

## Write a Go Function

1.  Create a directory with the files for your Go Function: `mkdir myfunc && cd myfunc && touch main.go; touch function.json`.

1.  Put the following code in `main.go`.

    ```go
    package main

    import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"

        "github.com/vladbarosan/func-go/azfunc"
    )

    // Run runs this Azure Function if/because it is specified in `function.json` as
    // the entryPoint. Fields of the function's parameters are also bound to
    // incoming and outgoing event properties as specified in `function.json`.
    func Run(ctx azfunc.Context, req *http.Request) (User, error) {

        // additional properties are bound to ctx by Azure Functions
        ctx.Log(azfunc.LogInformation,"function invoked: function %v, invocation %v", ctx.FunctionID(), ctx.InvocationID())

        // use Go's standard library to:
    	//  handle incoming request:
        body, _ := ioutil.ReadAll(req.Body)

        // to deserialize JSON content:
        var data map[string]interface{}
        var err error
        err = json.Unmarshal(body, &data)
        if err != nil {
            return nil, fmt.Errorf("failed to unmarshal JSON: %s\n", err)
        }

        // and to get query param values:
        name := req.URL.Query().Get("name")

        if name == "" {
            return nil, fmt.Errorf("missing required query parameter: name")
        }

    	// Prepare a struct to return. The special output binding name
    	// `$return` transforms the struct into near-equivalent JSON.
        u := &User{
            Name:     name,
            Greeting: fmt.Sprintf("Hello %s. %s\n", name, data["greeting"].(string)),
        }

        return u, nil
    }

    // User exemplifies a struct to be returned. You can use any struct or *struct.
    type User struct {
        Name     string
        Greeting string
    }
    ```

1.  Put the following configuration in the `function.json` file next to
    `main.go`. `function.json` specifies bindings between incoming and outgoing
    event properties and the structs and types in your code.

    For more details see [the function.json
    wiki](https://github.com/Azure/azure-functions-host/wiki/function.json).

    ```json
    {
      "entryPoint": "Run",
      "bindings": [
        {
          "name": "req"
          "type": "httpTrigger",
          "direction": "in",
          "authLevel": "anonymous",
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

## Deploy it

With your Function written, package and deploy it to a Go Functions instance.

If you need an instance see [Run an instance][].

[run an instance]: #run-a-go-functions-instance

**TODO(joshgav)**: add scripts and instructions.

## Trigger and watch it

Now your Function is live and ready to handle events. Time to trigger it!

1.  Use a tool like [Postman](https://www.getpostman.com/apps) or `curl` to
    execute a request with the following parameters (in this case for a local
    instance on port 8080):

        ```
        HTTP Method: `POST`
        URL: `http://localhost:8080/api/HttpTrigger?name=world`
        Headers: `Content-Type: application/json`
        Body: `{"greeting": "How are you?"}`
        ```

        ```bash
        declare PORT=8080 PERSON_NAME=world
        curl -L "http://localhost:${PORT}/api/HttpTrigger?name=${PERSON_NAME}" \
            --data '{ "greeting": "How are you?" }' \
            --header 'Content-Type: application/json' \
        ```

        The `Run` method from the sample should be executed and a User object with
        Name and Greeting properties like the following should be returned:

        ```json
        {
          "Name": "world",
          "Greeting": "Hello world. How are you?\n"
        }
        ```

# More information

## Things to note

- `function.json::entryPoint` names the Go function in package main to be used
  as the Azure Function entry point. In this example that function is named
  `Run` but any name is okay as long as it is also specified in
  `function.json`.
- `main.go` is the required name for the file containing the entry point Go
  function.
- You can use any dependencies you want in your app since they'll be compiled
  into the built binary.
- Structs in the function signature are initialized based on properties in the
  incoming event and specifications in function.json. In the example signature
  of `func Run(ctx azfunc.Context, req *http.Request) (User, error)`; `ctx azfunc.Context`, `req *http.Request` and `User` are automatically bound to
  incoming and outgoing message properties. Properties received from the GRPC
  channel are bound to properties on the Go structs, and any Go struct with the
  named properties can be used; that is, there's nothing special about the
  default types provided in package azfunc. This is illustrated by the returned
  `User` struct in the example.
- **Properties are bound to parameters based on the name of the parameter! You
  can change the order, but the name has to be consistent with the name of the
  binding defined in `function.json`!**
- You can specify a named return type, which then needs to match an output
  binding in `function.json`. Alternatively, you can have 1 unnamed return type
  which will match the special `$return` binding.
- You can also have an optional `error` return (named or anonymous) value to
  signal that the function execution failed for whatever reason.
- Having pointer types is preferred, but you can also have parameters and
  return values as non-pointer types for your functions.

## Disclaimer

- The project is in development; problems and frequent changes are expected.
- This project has not been evaluated for production use.
- We will reply to issues in this repo; but this project is in limited preview
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
Conduct](https://opensource.microsoft.com/codeofconduct/). For more
information see the [Code of Conduct
FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or contact
[opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional
questions or comments.
