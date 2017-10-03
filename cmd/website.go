package main

import (
	"html/template"
	"path/filepath"
	"os"
	"bufio"
	"bytes"
	"github.com/russross/blackfriday"
	"sort"
)

type Website struct {
	template *template.Template
	Posts    []Document
	Pages    []Document
}

type PageInfo struct {
	Document Document
	Website  Website
	Contents template.HTML
}

func New(postPath string, pagePath string, templatePath string) Website {
	t, err := template.ParseFiles(templatePath)

	if err != nil {
		panic(err)
	}

	posts := ReadFrom(postPath)
	pages := ReadFrom(pagePath)

	sort.Sort(NewestFirst(posts))
	sort.Sort(NewestFirst(pages))

	return Website {
		template: t,
		Posts:    posts,
		Pages:    pages,
	}
}

func (w Website) WriteTo(outputPath string) {
	for _, doc := range w.Pages {
		write(outputPath, doc, w)
	}

	postPath := filepath.Join(outputPath, "/posts")
	_ = os.Mkdir(postPath, 644)

	for _, doc := range w.Posts {
		write(postPath, doc, w)
	}
}

func write (path string, d Document, s Website) {
	path, err := filepath.Abs(filepath.Join(path, d.Id + ".html"))

	if err != nil {
		panic(err)
	}

	f, err := os.Create(path)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	pageInfo := PageInfo{
		Document: d,
		Website: s,
	}

	buf := new(bytes.Buffer)

	if err := d.Contents.Execute(buf, pageInfo); err != nil {
		panic(err)
	}

	pageInfo.Contents = template.HTML(blackfriday.MarkdownCommon(buf.Bytes()))

	w := bufio.NewWriter(f)

	if err = s.template.Execute(w, pageInfo); err != nil {
		panic(err)
	}

	w.Flush()
}
