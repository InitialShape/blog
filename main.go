package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "log"
    "net/http"
    "gopkg.in/russross/blackfriday.v2"
)

func contentHandler(w http.ResponseWriter, r *http.Request) {
    files := readFiles()
    for _, file := range files {
        if "content/"+file.Name()  == r.URL.Path[1:] {
            filePath := "./content/" + file.Name()
            b, err := ioutil.ReadFile(filePath)
            if err != nil {
                log.Fatal(err)
            }
            output := blackfriday.Run(b)
            fmt.Fprintf(w, string(output))
        }
    }
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    files := readFiles()
    for _, file := range files {
        fmt.Fprintf(w, "<a href=\"content/"+file.Name()+"\">"+file.Name()+"</a>")
    }
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
    log.Fatal(http.ListenAndServe(":8080", nil))
}
