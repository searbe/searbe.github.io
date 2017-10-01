package main

import (
	"text/template"
	"path/filepath"
	"io/ioutil"
	"encoding/json"
)

type Document struct {
	Meta     *Meta
	Id       string
	Contents *template.Template
}

type Meta struct {
	Published bool
	Date      string
	Title     string
}

func ReadFrom(path string) []*Document {
	postPaths, err := filepath.Glob(path + "/*.md")

	if err != nil {
		panic(err)
	}

	posts := make([]*Document, len(postPaths))

	for i, path := range postPaths {
		noext := path[:len(path)-3]
		id := filepath.Base(noext)

		var meta Meta

		rawMeta, err := ioutil.ReadFile(noext + ".json")

		if err != nil {
			meta = Meta{
				Published: true,
				Title:     id,
			}
		} else {
			err = json.Unmarshal(rawMeta, &meta)

			if err != nil {
				panic(err)
			}
		}

		markdown, err := ioutil.ReadFile(path)

		if err != nil {
			panic(err)
		}

		contents, err := template.New("contents").Parse(string(markdown))

		if err != nil {
			panic(err)
		}

		posts[i] = &Document{
			Id:       id,
			Meta:     &meta,
			Contents: contents,
		}
	}

	return posts
}