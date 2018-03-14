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
			markdown, err := article.ToMarkdown(false)
			if err != nil {
				log.Fatal(err)
			}

			t, _ := template.ParseFiles("./templates/article.html")
			t.Execute(w, template.HTML(markdown))
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

	var markdowns []byte
	for _, article := range articles {
		var markdown []byte

		articleMd, err := article.ToMarkdown(true)
		if err != nil {
			log.Fatal(err)
		}
		markdown = append(markdown, []byte("<section>")...)
		markdown = append(markdown, articleMd...)
		markdown = append(markdown, []byte("</section>")...)
		markdowns = append(markdowns, markdown...)
	}
	t.Execute(w, template.HTML(markdowns))

}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/content/", contentHandler)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
