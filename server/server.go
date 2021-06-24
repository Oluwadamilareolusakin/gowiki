package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
    "regexp"
	"github.com/oluwadamilareolusakin/gowiki/io"
)

var absPath, _ = os.Getwd()

var templateFiles = []string{"templates/404.html", "templates/edit.html", "templates/new.html", "templates/success.html", "templates/show.html"}

func parseTemplates() *template.Template {
  var templates *template.Template

  for _, filename := range templateFiles {
    filename = absPath + "/" + filename
    templates = template.Must(template.ParseFiles(filename))
  }

  return templates
}

var templates = parseTemplates()

const pagePath string = ".pages"

func handleError(w http.ResponseWriter, err error) {
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

var validPath = regexp.MustCompile("^/(new|edit|view|save)/([a-zA-z0-9]+)$")

func makeHandlerFunc(fn func(w http.ResponseWriter, r *http.Request, title string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    match := validPath.FindStringSubmatch(r.URL.Path)

    if len(match) >= 2 {
      fn(w, r, match[2])
      return
    }

    fn(w, r, "")
  }
}

func renderTemplate(w http.ResponseWriter, title string, page *io.Page) {
  err := templates.ExecuteTemplate(w, title + ".html", page)

  handleError(w, err)
}

func viewHandler(w http.ResponseWriter, r *http.Request, filename string) {
  page, err := io.LoadPage(pagePath, filename)

  if err != nil {
    renderTemplate(w, "404", &io.Page{Title: filename})
    return
  }

  renderTemplate(w, "show", page)
}

func editHandler(w http.ResponseWriter, r *http.Request, filename string) {
  page, err := io.LoadPage(pagePath, filename)

  handleError(w, err)

  renderTemplate(w, "edit", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request, filename string) {
  body := r.FormValue("body")
  title := r.FormValue("title")

  page := &io.Page {Title: title, Body: []byte(body)}

  err := page.Save(pagePath)

  handleError(w, err)

  renderTemplate(w, "success", page)
}

func createHandler(w http.ResponseWriter, r *http.Request, filename string) {
  renderTemplate(w, "new", nil)
}

func main() {
  http.HandleFunc("/view/", makeHandlerFunc(viewHandler))
  http.HandleFunc("/new", makeHandlerFunc(createHandler))
  http.HandleFunc("/save", makeHandlerFunc(saveHandler))
  http.HandleFunc("/edit/", makeHandlerFunc(editHandler))
  fmt.Println("Starting server on port 8080...")
  fmt.Println("Server started, access it at http://localhost:8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
