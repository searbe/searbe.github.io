package main

import (
	"text/template"
	"path/filepath"
	"io/ioutil"
	"encoding/json"
	"time"
)

type Document struct {
	Meta     *Meta
	Id       string
	Contents *template.Template
}

type Meta struct {
	Date      time.Time
	Title     string
}

type MetaJson struct {
	Date      string
	Title     string
}

type NewestFirst []Document

func (a NewestFirst) Len() int           { return len(a) }
func (a NewestFirst) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NewestFirst) Less(i, j int) bool { return a[i].Meta.Date.After(a[j].Meta.Date) }

type OldestFirst []Document

func (a OldestFirst) Len() int           { return len(a) }
func (a OldestFirst) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a OldestFirst) Less(i, j int) bool { return a[i].Meta.Date.Before(a[j].Meta.Date) }

func ReadFrom(path string) []Document {
	postPaths, err := filepath.Glob(path + "/*.md")

	if err != nil {
		panic(err)
	}

	posts := make([]Document, len(postPaths))

	for i, path := range postPaths {
		noext := path[:len(path)-3]
		id := filepath.Base(noext)

		var metaJson MetaJson

		rawMeta, err := ioutil.ReadFile(noext + ".json")
		meta := Meta{}

		if err != nil {
			meta.Title = id
		} else {
			err = json.Unmarshal(rawMeta, &metaJson)

			if err != nil {
				panic(err)
			}

			meta.Title = metaJson.Title
			t, _ := time.Parse("02/01/2006 15:04:05", metaJson.Date)

			if err != nil {
				panic(err)
			}

			meta.Date = t
		}

		markdown, err := ioutil.ReadFile(path)

		if err != nil {
			panic(err)
		}

		contents, err := template.New("contents").Parse(string(markdown))

		if err != nil {
			panic(err)
		}

		posts[i] = Document{
			Id:       id,
			Meta:     &meta,
			Contents: contents,
		}
	}

	return posts
}