package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	swaggeruiserver "github.com/ibraimgm/swaggerui-server"
)

func main() {
	addr := flag.String("addr", ":8080", "the address and port to listen")
	location := flag.String("location", "/", "the url location to use for the documentation")
	docStr := flag.String("docs", "", "a comma-separated list of documents in the format NAME=URL")
	docsFile := flag.String("file", "", "a file with the list of documents in the format NAME=URL, separated by newline")

	flag.Parse()

	if !strings.HasSuffix(*location, "/") {
		*location += "/"
	}

	// document list from command-line
	docs := buildDocs(nil, *docStr)

	// document list from file
	docs, err := docsFromFile(docs, *docsFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if len(docs) == 0 {
		fmt.Fprintln(os.Stderr, "No document list provided.")
		os.Exit(1)
	}

	log.Println("Serving documents:")
	for _, d := range docs {
		if d.URL != d.Name {
			log.Printf("\t%s: %s", d.Name, d.URL)
		} else {
			log.Printf("\t%s", d.Name)
		}
	}

	log.Println("Listening at ", *addr)
	log.Println(http.ListenAndServe(*addr, swaggeruiserver.MustAt(*location, docs)))
}

func buildDocs(docs []swaggeruiserver.Doc, value string) []swaggeruiserver.Doc {
	for _, s := range strings.Split(value, ",") {
		item := strings.Split(s, "=")
		name, url := item[0], item[0]

		if name == "" {
			continue
		}

		if len(item) > 1 {
			url = item[1]
		}

		docs = append(docs, swaggeruiserver.Doc{URL: url, Name: name})
	}

	return docs
}

func docsFromFile(docs []swaggeruiserver.Doc, filename string) ([]swaggeruiserver.Doc, error) {
	if filename == "" {
		return docs, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return docs, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing documents file: %v", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		docs = buildDocs(docs, scanner.Text())
	}

	return docs, scanner.Err()
}
