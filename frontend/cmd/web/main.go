package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	fmt.Println("Starting frontend service on port 80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Panic(err)
	}
}

// go:embed templates
var templateFS embed.FS

func render(w http.ResponseWriter, t string) {
	partials := []string{
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	var templates []string
	templates = append(templates, fmt.Sprintf("./cmd/web/templates/%s", t))

	for _, x := range partials {
		templates = append(templates, x)
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
