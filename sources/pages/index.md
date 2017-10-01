# Welcome
This is the index page.

{{ range $post := .Website.Posts }}
* [{{ $post.Meta.Title }}](/posts/{{$post.Id}}.html)
{{ end }}