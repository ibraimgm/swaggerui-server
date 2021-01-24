package swaggeruiserver_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	swaggeruiserver "github.com/ibraimgm/swaggerui-server"
)

var testDocuments []swaggeruiserver.Doc = []swaggeruiserver.Doc{
	{URL: "fooURL", Name: "FooURL"},
	{URL: "barURL", Name: "BarURL"},
	{URL: "bazURL", Name: "BazURL"},
}

func runRequest(mux *http.ServeMux, url string) *httptest.ResponseRecorder {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	return rec
}

func checkStatus(mux *http.ServeMux, url string, status int) error {
	if rec := runRequest(mux, url); rec.Code != status {
		return fmt.Errorf("expected URL '%s' to return status %d, but got %d", url, status, rec.Code)
	}

	return nil
}

func checkSwaggerIndex(mux *http.ServeMux, url string) error {
	rec := runRequest(mux, url)

	// when a "directory" doesn't end with "/", it should redirect
	if !strings.HasSuffix(url, "/") && rec.Code == http.StatusMovedPermanently {
		return nil
	}

	if rec.Code != http.StatusOK {
		return fmt.Errorf("expected URL '%s' to give %d status, but received %d", url, http.StatusOK, rec.Code)
	}

	body := rec.Body.String()

	for _, doc := range testDocuments {
		if !strings.Contains(body, doc.Name) || !strings.Contains(body, doc.URL) {
			return fmt.Errorf("document '%s' not found on URL '%s'", doc.Name, url)
		}
	}

	return nil
}

func runSwagger(mux *http.ServeMux, swaggerRoot string) error {
	swaggerRoot = strings.TrimSuffix(swaggerRoot, "/")

	if err := checkSwaggerIndex(mux, swaggerRoot+"/"); err != nil {
		return err
	}

	if err := checkSwaggerIndex(mux, swaggerRoot+"/index.html"); err != nil {
		return err
	}

	return checkStatus(mux, swaggerRoot+"/static/favicon-32x32.png", http.StatusOK)
}

func TestHandle(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
	}{
		{name: "Docs1", pattern: "/docs"},
		{name: "Docs2", pattern: "/docs/"},
		{name: "Root", pattern: "/"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/foo", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK")) //nolint:errcheck
			})

			if err := swaggeruiserver.Handle(mux, test.pattern, testDocuments); err != nil {
				t.Fatal(err)
			}

			if err := runSwagger(mux, test.pattern); err != nil {
				t.Fatal(err)
			}

			if err := checkStatus(mux, "/api/foo", http.StatusOK); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAt(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
	}{
		{name: "Docs1", pattern: "/docs"},
		{name: "Docs2", pattern: "/docs/"},
		{name: "Root", pattern: "/"},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			mux := swaggeruiserver.MustAt(test.pattern, testDocuments)

			if err := runSwagger(mux, test.pattern); err != nil {
				t.Fatal(err)
			}

			if err := checkStatus(mux, test.pattern+"xyz", http.StatusNotFound); err != nil {
				t.Fatal(err)
			}

			if err := checkStatus(mux, "xyz"+test.pattern, http.StatusNotFound); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestMux(t *testing.T) {
	mux := swaggeruiserver.MustMux(testDocuments)

	if err := runSwagger(mux, "/"); err != nil {
		t.Fatal(err)
	}

	if err := checkStatus(mux, "/xyz", http.StatusNotFound); err != nil {
		t.Fatal(err)
	}

	if err := checkStatus(mux, "xyz/", http.StatusNotFound); err != nil {
		t.Fatal(err)
	}
}
