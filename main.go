package main

import (
	"github.com/InitialShape/blog/models"
	"html/template"
	"log"
	"net/http"
	"sort"
)

func contentHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := models.GetArticles()
	if err != nil {
		log.Fatal(err)
	}

	for _, article := range articles {
		if article.Path == r.URL.Path[1:] {
			html, err := article.ToHTML(false)
			if err != nil {
				log.Fatal(err)
			}

			t, _ := template.ParseFiles("./templates/article.html")
			t.Execute(w, html)
		}
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := models.GetArticles()
	sort.Sort(articles)
	t, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	for i, article := range articles {

                html, err := article.ToHTML(true)
                article.HTML = html
		if err != nil {
			log.Fatal(err)
		}
                articles[i] = article
	}
	t.Execute(w, articles)

}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/content/", contentHandler)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
