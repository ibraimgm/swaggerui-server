package main

import (
	"log"
	"net/http"

	swaggeruiserver "github.com/ibraimgm/swaggerui-server"
)

func main() {
	docs := []swaggeruiserver.Doc{
		{URL: "https://petstore.swagger.io/v2/swagger.json", Name: "PetStore"},
	}

	log.Println(http.ListenAndServe(":8080", swaggeruiserver.MustMux(docs)))
}
