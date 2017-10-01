package main

import (
	"flag"
	"path/filepath"
)

func main() {
	inputPath := flag.String("posts", "./sources", "Path containing \"Posts\" directory (containing .md Posts) and template.tpl file")
	outputPath := flag.String("outputPath", "./", "Path in which to write static outputPath")

	posts, err := filepath.Abs(filepath.Join(*inputPath, "posts"))

	if err != nil {
		panic(err)
	}

	pages, err := filepath.Abs(filepath.Join(*inputPath, "pages"))

	if err != nil {
		panic(err)
	}

	template, err := filepath.Abs(filepath.Join(*inputPath, "template.gohtml"))

	if err != nil {
		panic(err)
	}

	output, err := filepath.Abs(*outputPath)

	if err != nil {
		panic(err)
	}

	site := New(posts, pages, template)
	site.WriteTo(output)
}

