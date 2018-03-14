package main

import (
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Article struct {
	Path      string
	Info      os.FileInfo
	CreatedAt time.Time
}

type Articles []Article

func (p Articles) Len() int {
	return len(p)
}

func (p Articles) Less(i, j int) bool {
	return p[i].CreatedAt.After(p[j].CreatedAt)
}

func (p Articles) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (article Article) ToMarkdown() ([]byte, error) {
	b, err := ioutil.ReadFile(article.Path)
	if err != nil {
		return nil, err
	}
	return blackfriday.Run(b), err
}

func pathToTime(path string) (time.Time, error) {
	first := strings.IndexByte(path, '/') + 1
	last := strings.LastIndex(path, "/")
	dateString := path[first:last]
	date, err := time.Parse("2006/03/02", dateString)

	if err != nil {
		log.Fatal(err)
		return time.Time{}, err
	}

	return date, err
}

func contentHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := getArticles()
	if err != nil {
		log.Fatal(err)
	}

	for _, article := range articles {
		if article.Path == r.URL.Path[1:] {
			markdown, err := article.ToMarkdown()
			if err != nil {
				log.Fatal(err)
			}

			t, _ := template.ParseFiles("./templates/article.html")
			t.Execute(w, template.HTML(markdown))
		}
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := getArticles()
	sort.Sort(articles)
	t, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	var markdowns []byte
	for _, article := range articles {
                var markdown []byte

                articleMd, err := article.ToMarkdown()
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

func getArticles() (Articles, error) {
	var articles Articles
	var err error

	err = filepath.Walk("./content", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
			return nil
		}
		match, err := regexp.MatchString(".md", info.Name())
		if match {

			date, err := pathToTime(path)
			if err != nil {
				log.Fatal(err)
				return nil
			}

			article := Article{path, info, date}
			articles = append(articles, article)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return articles, err
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/content/", contentHandler)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
