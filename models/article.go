package models

import (
	"gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
        "bytes"
	"time"
        "github.com/InitialShape/blog/utils"
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

func (article Article) ToMarkdown(excerpt bool) ([]byte, error) {
	b, err := ioutil.ReadFile(article.Path)
	if err != nil {
		return nil, err
	}
        markdown := blackfriday.Run(b)
        if excerpt {
                excerptEnd := bytes.Index(markdown, []byte("</p>"))
                if excerptEnd != -1 {
                        markdown = markdown[:excerptEnd]
                }
        }
        return markdown, err
}

func GetArticles() (Articles, error) {
	var articles Articles
	var err error

	err = filepath.Walk("./content", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
			return nil
		}
		match, err := regexp.MatchString(".md", info.Name())
		if match {

			date, err := utils.PathToTime(path)
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