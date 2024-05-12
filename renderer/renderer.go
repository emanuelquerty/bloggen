package renderer

import (
	"embed"
	"html/template"
	"io"

	"github.com/emanuelquerty/bloggen/parser"
	"github.com/gomarkdown/markdown"
	markdownparser "github.com/gomarkdown/markdown/parser"
)

var (
	//go:embed "templates/*"
	postTemplates embed.FS
)

type PostRenderer struct {
	templ *template.Template
	mdParser *markdownparser.Parser
}

func NewPostRenderer() (*PostRenderer, error) {
	templ, err := template.ParseFS(postTemplates, "templates/*.gohtml")
	if err != nil {
		return nil, err
	}

	extensions := markdownparser.CommonExtensions | markdownparser.AutoHeadingIDs
	markdownparser := markdownparser.NewWithExtensions(extensions)
	return &PostRenderer{templ: templ, mdParser: markdownparser}, nil
}

func (r *PostRenderer) Render(w io.Writer, post parser.Post) error {
	return r.templ.ExecuteTemplate(w, "blog.gohtml", newPostVM(post, r, true))
}

func (r *PostRenderer) RenderIndex(w io.Writer, posts []parser.Post) error {
	return r.templ.ExecuteTemplate(w, "index.gohtml", posts)
}

type postViewModel struct {
	parser.Post
	HTMLBody template.HTML
	IsPost bool
}

func newPostVM(p parser.Post, r *PostRenderer, isPost bool) postViewModel {
	vm := postViewModel{Post: p, IsPost: isPost}

	extensions := markdownparser.CommonExtensions | markdownparser.AutoHeadingIDs
	markdownparser := markdownparser.NewWithExtensions(extensions)

	vm.HTMLBody = template.HTML(markdown.ToHTML([]byte(p.Body),  markdownparser, nil))
	return vm
}
