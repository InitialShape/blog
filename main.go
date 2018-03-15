package main

import (
	"github.com/InitialShape/blog/models"
	"html/template"
	"log"
	"net/http"
	"sort"
        "strconv"
)

const PAGE_SIZE = 4

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
        page, err := strconv.Atoi(r.URL.Query().Get("page"))
	articles, err := models.GetArticles()
	sort.Sort(articles)

        if page*PAGE_SIZE+PAGE_SIZE > len(articles) && page*PAGE_SIZE > len(articles) {
                w.WriteHeader(http.StatusNotFound)
                return
        } else {
            if page*PAGE_SIZE < len(articles) && page*PAGE_SIZE+PAGE_SIZE > len(articles) {
                    articles = articles[PAGE_SIZE*page:len(articles)]
            } else {
                    articles = articles[PAGE_SIZE*page:page*PAGE_SIZE+PAGE_SIZE]
            }
        }


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

	t.Execute(w, map[string]interface{}{
                "NextPage": page+1,
                "PrevPage": page-1,
                "Articles": articles,
        })

}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/content/", contentHandler)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
