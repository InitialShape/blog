package utils

import (
	"log"
	"strings"
	"time"
)

func PathToTime(path string) (time.Time, error) {
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
