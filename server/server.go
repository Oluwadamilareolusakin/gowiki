package main

import (
  "fmt"
  "html/template"
  "log"
  "net/http"
  "net/url"
  "os"
  "regexp"

  "github.com/oluwadamilareolusakin/gowiki/pageio"
  "github.com/oluwadamilareolusakin/gowiki/statictemplates"
)

//go:generate go get "github.com/oluwadamilareolusakin/embedfiles"
//go:generate embedfiles "../statictemplates/template.go" "../templates" "statictemplates"

const pagePath string = ".pages"

func handleError(w http.ResponseWriter, err error) {
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

var validPath = regexp.MustCompile("^/(new|edit|view|save)/([a-zA-z0-9\\s-]+)$")

func makeHandlerFunc(fn func(w http.ResponseWriter, r *http.Request, title string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    match := validPath.FindStringSubmatch(r.URL.Path)

    if len(match) >= 2 {
      path, _ := url.QueryUnescape(match[2])
      fn(w, r, path)
      return
    }

    fn(w, r, "")
  }
}

func renderTemplate(w http.ResponseWriter, title string, page *pageio.Page) {

  templateKey := "/" + title + ".html"
  data := string(statictemplates.Get(templateKey))
  templ := template.Must(template.New("").Parse(data))
  err := templ.Execute(w, page)
  handleError(w, err)
}

func renderNotFound(w http.ResponseWriter, filename string) {
  renderTemplate(w, "404", &pageio.Page{Title: filename})
}

func viewHandler(w http.ResponseWriter, r *http.Request, filename string) {
  page, err := pageio.LoadPage(pagePath, filename)

  if err != nil {
    renderNotFound(w, filename)
    return
  }

  renderTemplate(w, "show", page)
}

func editHandler(w http.ResponseWriter, r *http.Request, filename string) {
  page, err := pageio.LoadPage(pagePath, filename)

  if os.IsNotExist(err) {
    renderNotFound(w, filename)
    return
  }

  handleError(w, err)

  renderTemplate(w, "edit", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request, filename string) {
  body := r.FormValue("body")
  title := r.FormValue("title")

  page := &pageio.Page {Title: title, Body: []byte(body)}

  err := page.Save(pagePath)

  handleError(w, err)

  renderTemplate(w, "success", page)
}

func createHandler(w http.ResponseWriter, r *http.Request, filename string) {
  renderTemplate(w, "new", nil)
}

func init() {
  if _, err := os.Stat(pagePath); os.IsNotExist(err) {
    os.Mkdir(pagePath, os.ModePerm)
  }
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
