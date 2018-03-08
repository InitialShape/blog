package main

import (
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func contentHandler(w http.ResponseWriter, r *http.Request) {
	files := readFiles()
	for _, file := range files {
		if "content/"+file.Name() == r.URL.Path[1:] {
		        filePath := "./content/" + file.Name()
			b, err := ioutil.ReadFile(filePath)
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
	files := readFiles()
	t, err := template.ParseFiles("./templates/index.html")
        if err != nil {
            log.Fatal(err)
        }
	t.Execute(w, &files)

}

func readFiles() []os.FileInfo {
	files, err := ioutil.ReadDir("./content")
	if err != nil {
		log.Fatal(err)
	}

	return files
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/content/", contentHandler)
        fs := http.FileServer(http.Dir("static"))
        http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
