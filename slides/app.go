package slides

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/gopherguides/hype"
)

//go:embed templates/slides.html
var HTMLTemplate string

//go:embed templates/assets/app.js
var AppJS string

//go:embed templates/assets/*.*
var AssetsFS embed.FS

// func init() {
// 	err := fs.WalkDir(AssetsFS, ".", func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		fmt.Println(path)

// 		return nil
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// }

type App struct {
	FileName string
	PWD      string
	Parser   *hype.Parser // If nil, a default parser is used.

}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a == nil {
		http.Error(w, "nil app", http.StatusInternalServerError)
		return
	}
	if err := a.serve(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// a.once.Do(func() {
	// 	a.mux = http.NewServeMux()
	// 	// a.mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(a.Parser.FS))))
	// 	a.mux.Handle("/assets/", http.FileServer(http.FS(a.Parser.FS)))
	// 	a.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 		fmt.Printf("TODO >> app.go:60 r.URL %[1]T %+[1]v\n", r.URL)
	// 		if err := a.serve(w, r); err != nil {
	// 			http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		}
	// 	})
	// })

	// a.mux.ServeHTTP(w, r)
}

func (a *App) serve(w http.ResponseWriter, r *http.Request) error {
	if r == nil {
		return fmt.Errorf("nil request")
	}

	cab := os.DirFS(a.PWD)
	p := hype.NewParser(cab)
	p.Root = a.PWD

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, a.FileName)
	if err != nil {
		return fmt.Errorf("pwd: %q, file: %q, err: %w", a.PWD, a.FileName, err)
	}

	pages, err := doc.Pages()
	if err != nil {
		return err
	}

	var nodes hype.Nodes
	for i, page := range pages {
		page.Set("id", fmt.Sprintf("page-%d", i))
		nodes = append(nodes, page)
	}

	bs := nodes.String()

	tmpl, err := template.New("slides").Parse(HTMLTemplate)
	if err != nil {
		return err
	}

	data := map[string]any{
		"title": doc.Title,
		"body":  template.HTML(bs),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
