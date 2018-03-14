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
)

type Article struct {
	Path string
	Info os.FileInfo
}

func contentHandler(w http.ResponseWriter, r *http.Request) {
	files, err := readDirectory()
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Path == r.URL.Path[1:] {
			b, err := ioutil.ReadFile(file.Path)
			if err != nil {
				log.Fatal(err)
			}
			output := blackfriday.Run(b)
			t, _ := template.ParseFiles("./templates/article.html")
			t.Execute(w, template.HTML(output))
		}
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	files, err := readDirectory()
	t, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, &files)

}

func readDirectory() ([]Article, error) {
	var articles []Article
	var err error

	err = filepath.Walk("./content", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}
		match, err := regexp.MatchString(".md", info.Name())
		if match {
			article := Article{path, info}
			articles = append(articles, article)
		}
		return nil
	})

	return articles, err
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/content/", contentHandler)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
