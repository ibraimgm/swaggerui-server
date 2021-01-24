# SwaggerUI Server

The easiest way to run SwaggerUI in your Go application.

![CI](https://github.com/ibraimgm/swaggerui-server/workflows/CI/badge.svg)
[![codecov](https://codecov.io/gh/ibraimgm/swaggerui-server/branch/master/graph/badge.svg?token=jX55quRBda)](https://codecov.io/gh/ibraimgm/swaggerui-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/ibraimgm/swaggerui-server)](https://goreportcard.com/report/github.com/ibraimgm/swaggerui-server)
[![SwaggerUI](https://img.shields.io/badge/SwaggerUI-v3.40.0-blue)](https://github.com/swagger-api/swagger-ui/releases/tag/v3.40.0)
[![Docker](https://img.shields.io/badge/docker-latest-blue)](https://hub.docker.com/r/ibraimgm/swaggerui-server)

## Usage

### Command-line application

Use the `swaggerui-server`  command-line application to easily display one or more Swagger definitions. First, install the tool with:

```shell
# install
$ go install github.com/ibraimgm/swaggerui-server/cmd/...

# check the available options
$ swaggerui-server -h
Usage of ./swaggerui-server:
  -addr string
        the address and port to listen (default ":8080")
  -docs string
        a comma-separated list of documents in the format NAME=URL
  -file string
        a file with the list of documents in the format NAME=URL, separated by newline
  -location string
        the url location to use for the documentation (default "/")
```

To provide one or more documents to display, use the `-docs` command-line argument:

```shell
# Serves only the 'PetStore' demo
$ swaggerui-server -docs https://petstore.swagger.io/v2/swagger.json

# Same as above, but show a nice name instead of the URL in the UI
$ swaggerui-server -docs PetStore=https://petstore.swagger.io/v2/swagger.json

# You can specify more than one definition as a comma-separated list
$ swaggerui-server -docs PetStore=https://petstore.swagger.io/v2/swagger.json,Logz=https://raw.githubusercontent.com/logzio/public-api/master/alerts/swagger.json
```

If you need to show multiple documents, passing all of them via a command-line parameter is a bit cumbersome. A better option is to define the list of documents into a file and use the `-file` option:

```shell
# The file below has a list of documents, in the same format
# accepted by the `-docs` parameter
$ cat documents.txt
PetStore=https://petstore.swagger.io/v2/swagger.json
Logz.io=https://raw.githubusercontent.com/logzio/public-api/master/alerts/swagger.json
SwaggerTools=https://raw.githubusercontent.com/apigee-127/swagger-tools/master/examples/2.0/api/swagger.json

# Just provide the file and you're good to go
$ swaggerui-server -file documents.txt
```

In all of the above example, the UI is available at `http://localhost:8080/`. You can change the address(`-addr`) or the location(`-location`) to more suitable values, if you wish. For example, `./swaggerui-server -addr :9090 -location /docs` will serve the UI at `http://localhost:9090/docs`.

### Docker image

If you don't want to compile the server yourself, another alternative is to use a prebuilt docker image:

```shell
# Try this and take a look at http://localhost:8080
$ docker run --rm -it -p 8080:8080 ibraimgm/swaggerui-server
```

The prebuilt image is based on alpine and have `swaggerui-server` installed on `/usr/local/bin`. The default `CMD` of the image starts with the PetShop demo document, so you certainly want to customize it to a more sensible value.

### As a library

You can use `swaggerui-server` as a library to server the UI directly inside your application. For example, if you with to serve the UI on the `/docs` URL, you can do:

```go
package main

import (
  "log"
  "net/http"

  swaggeruiserver "github.com/ibraimgm/swaggerui-server"
)

func main() {
  // create a server mux and add your API routes
  mux := http.NewServeMux()
  mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, this is my cool API"))
  })

  // add one or more swagger documentation files
  // of course, you can also serve one or more
  // files from the same ServerMux.
  //
  // as long as you remember that the url must be
  // reachable by the web browser, anything goes!
  docs := []swaggeruiserver.Doc{
    {URL: "https://petstore.swagger.io/v2/swagger.json", Name: "PetStore"},
  }

  // register the handler in the correct path
  if err := swaggeruiserver.Handle(mux, "/docs", docs); err != nil {
    log.Fatal(err)
  }

  // run you application
  log.Println(http.ListenAndServe(":8080", mux))
}
```

## Re-generating static resources

All the needed SwaggerUI files are already bundled in `internal/assets/static.go`. If wish to regenerate this file
(for example, to put a custom CSS or to use a different SwaggerUI version), you should follow these
steps:

1. Put the SwaggerUI bundle (the `dist` folder fo the [official repository](https://github.com/swagger-api/swagger-ui)) in the root folder of this project, with the name `static`.
2. Run `make static`. This will apply a patch to the original `index.html` and rename it to `index.template`. Then, the `*.map` files will be discarded and the `static.go` file will be recreated.

If the patching doesn't workwith the specific SwaggerUI that you're using (or if you're using a customized `index.html`), you can just manually create your `static/index.template` first, and then run `make static`. If the template file already exists (and this is the case) only the `static.go` file will be generated. In this scenario, please make sure that you removed all the uneeded files (e. g. `*.map`), to avoid store uneeded data.
