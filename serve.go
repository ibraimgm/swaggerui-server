package swaggeruiserver

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"text/template"

	"github.com/ibraimgm/swaggerui-server/internal/assets"
)

// Doc represents a swagger document to be shown in the SwaggerUI.
//
// The URL can be any kind of address, as long as it is reachable by the
// web browser client. No checks or validations are don in the URL; the
// value is passed directly to the SwaggerUI javascript code.
//
// The Name, although optional, is used to show the user which of the
// documents he is browsing. If you leave it blank the UI will still
// work without problems but it will be quite odd, from the end user
// perspective, to stare at a blank combobox at the top of the page.
type Doc struct {
	URL  string
	Name string
}

type templateData struct {
	Prefix string
	Items  []Doc
}

var once sync.Once
var indexTpl *template.Template

// Handle adds the required SwaggerUI handlers to mux.
//
// This function adds two handlers to an existing mux: one on
// "{pattern}/" to serve the SwaggerUI index file, and one on
// "{pattern}/static"to serve the required static files, like
// CSS, images, etc.
//
// If no pattern is provided, "/" is assumed.
//
// This function is most useful when you need to bundle the SwaggerUI
// into an existing application.
func Handle(mux *http.ServeMux, pattern string, docs []Doc) error {
	var templateError error

	once.Do(func() {
		b, err := assets.Asset("index.template")
		if err != nil {
			templateError = fmt.Errorf("Failed to load template asset: %w", err)
			return
		}

		tpl, err := template.New("index").Parse(string(b))
		if err != nil {
			templateError = fmt.Errorf("Failed to parse template: %w", err)
			return
		}

		indexTpl = tpl
	})

	if templateError != nil {
		return templateError
	}

	pattern = strings.TrimSuffix(pattern, "/")
	rootURL := pattern + "/"
	indexURL := pattern + "/index.html"
	staticURL := pattern + "/static/"

	var sb strings.Builder
	if err := indexTpl.Execute(&sb, templateData{Prefix: pattern, Items: docs}); err != nil {
		return fmt.Errorf("Failed to execute template content: %w", err)
	}

	templateOutput := []byte(sb.String())

	mux.Handle(staticURL, http.StripPrefix(staticURL, http.FileServer(assets.AssetFile())))
	mux.HandleFunc(rootURL, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != rootURL && r.URL.Path != indexURL {
			http.NotFound(w, r)
			return
		}

		w.Write(templateOutput) //nolint:errcheck
	})

	return nil
}

// At returns a ServerMux that serves a SwaggerUI configured
// to with the provided docs.
//
// The main file is served with "{pattern}/index.html", while the
// needed static content is acessible with "{pattern}/static".
func At(pattern string, docs []Doc) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	return mux, Handle(mux, pattern, docs)
}

// MustAt works the same as At, but panics on error.
//
// Taking into account the fact that most of the errors returned
// by At are internal errors, this routine should be safe for the
// majority of the users.
func MustAt(pattern string, docs []Doc) *http.ServeMux {
	mux, err := At(pattern, docs)
	if err != nil {
		panic(err)
	}

	return mux
}

// Mux returns a *http.ServerMux serves the SwaggerUI at the root location ("/"),
// configured to show the provided swagger documents.
//
// The generated ServerMux compiles a template with the provided documents
// and serves it as "/index.html". The parsing and compilation is done only
// once, when the ServerMux is created, and apart from the index file,
// the other needed static resources are served at "/static/".
//
// If you need to customize the url prefix (e. g. to "/docs"), use the At function.
func Mux(docs []Doc) (*http.ServeMux, error) {
	return At("", docs)
}

// MustMux works the same as Mux, but panics on error.
//
// Taking into account the fact that most of the errors returned
// by Mux are internal errors, this routine should be safe for the
// majority of the users.
func MustMux(docs []Doc) *http.ServeMux {
	return MustAt("", docs)
}
