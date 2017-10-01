package main

import (
	"html/template"
)

type TemplateInput struct {
	Document *Document
	Website  *Website
	Contents template.HTML
}
