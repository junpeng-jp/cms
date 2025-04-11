package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	pageTemplateDir = "./pages"
)

var (
	host       string
	port       string
	staticPath string
	pagePath   string
)

func main() {
	ctx := context.Background()

	flag.StringVar(&host, "host", "localhost", "host name")
	flag.StringVar(&port, "port", "8080", "port")
	flag.StringVar(&staticPath, "static", "./static", "path to static assets")
	flag.StringVar(&pagePath, "page", "./page", "path to page templates")

	flag.Parse()

	if err := runServer(ctx, host, port, staticPath, pagePath); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func runServer(
	ctx context.Context,
	host, port, staticPath, pagePath string,
) error {
	server := newServer(staticPath, pagePath)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: server,
	}

	log.Printf("Server listening on %s", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error listening and serving: %s\n", err)
	}
	return nil
}

func newServer(staticPath string, pagePath string) http.Handler {
	mux := http.NewServeMux()

	// get page template
	pages := make(map[string]*template.Template, 1)

	err := filepath.Walk(pagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Base(path) != "index.html" {
			return nil
		}
		// Extract the template name based on the directory structure
		// Example: template/about/index.html becomes "about/index.html"
		name, err := filepath.Rel(pagePath, path)
		if err != nil {
			log.Printf("Error parsing template path '%s': %v", path, err)
			return err
		}
		name = filepath.Join("/", filepath.Dir(name))
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			log.Printf("Error parsing template '%s': %v", path, err)
			return err
		}
		pages[name] = tmpl
		log.Printf("Loaded template: %s as %s", path, name)

		return nil
	})
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	setupRoutes(mux, pages, staticPath)

	return mux
}

func setupRoutes(
	handler *http.ServeMux,
	templateMap map[string]*template.Template,
	staticPath string,
) {
	fs := http.FileServer(http.Dir(staticPath))
	handler.Handle("/static/", http.StripPrefix("/static/", fs))
	handler.HandleFunc("GET /api/render", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<article><p>Hello World</p></article>"))
	})
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path
		// Construct the expected template name
		templateName := filepath.Join("/", strings.TrimSuffix(path, "/index.html"))
		// Check if the template exists
		tmpl, found := templateMap[templateName]
		if !found {
			http.Error(w, fmt.Sprintf("%s not found", templateName), http.StatusNotFound)
			return
		}
		// Execute the template
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			log.Printf("Error executing template '%s': %v", templateName, err)
		}
	})
}
