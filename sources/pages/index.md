{{ range $post := .Website.Posts }}
* [{{ $post.Meta.Title }}](/posts/{{$post.Id}}.html) *{{ $post.Meta.Date.Format "Jan 02, 2006" }}*{{ end }}