package main

import (
	"errors"
	// "fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]*)$")
var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html", "tmpl/list.html"))

type Page struct {
	Title string
	Body  []byte
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression.
}

func (p *Page) save() error {
	filename := "data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func loadAllPage() ([]*Page, error) {
	files, err := ioutil.ReadDir("data")
	if err != nil {
		return nil, errors.New("invalid Dir name")
	}

	pages := []*Page{}
	for _, file := range files {
		reg := regexp.MustCompile("(.*).txt")
		subMatch := reg.FindStringSubmatch(file.Name())
		title := subMatch[1]

		page, err := loadPage(title)
		if err != nil {
			return nil, errors.New("invalid Page name")
		}
		pages = append(pages, &Page{Title: title, Body: page.Body})
	}
	return pages, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	pages, err := loadAllPage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "list.html", pages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	log.Print(title)
	if title == "" {
		renderTemplate(w, "edit", nil)
		return
	}
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	if title == "" {
		title = r.FormValue("title")
	}
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// 関数リテラルとクロージャーを使って、各ハンドラーで必須の処理（ページタイトルの取得・バリデーション）を共通化
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title, err := getTitle(w, r)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, title)
	}
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/list/", listHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
