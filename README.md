# SwaggerUI Server

The easiest way to run SwaggerUI in your Go application.

![CI](https://github.com/ibraimgm/swaggerui-server/workflows/CI/badge.svg)

## Usage

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

All the needed SwaggerUI files are already bundled in `static.go`. If wish to regenerate this file
(for example, to put a custom CSS or to use a different SwaggerUI version), you should follow these
steps:

1. Put the SwaggerUI bundle (the `dist` folder fo the [official repository](https://github.com/swagger-api/swagger-ui)) in the root folder of this project, with the name `static`.
2. Run `make static`. This will apply a patch to the original `index.html` and rename it to `index.template`. Then, the `static.go` file will be recreated.

If the patching doesn't workwith the specific SwaggerUI that you're using (or if you're using a customized `index.html`), you can just manually create your `static/index.template` first, and then run `make static`. If the template file already exists (and this is the case) only the `static.go` file will be generated.
