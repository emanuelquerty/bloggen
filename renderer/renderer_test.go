package renderer_test

import (
	"bytes"
	"io"
	"testing"

	approvals "github.com/approvals/go-approval-tests"
	"github.com/emanuelquerty/bloggen/parser"
	"github.com/emanuelquerty/bloggen/renderer"
)
func TestRenderer(t *testing.T) {
	var (
		aPost = parser.Post{
			Title:       "hello world",
			Body:        "This is a post",
			Description: "This is a description",
			Tags: []string{"go", "tdd"},
		}
	)

	postRenderer, err := renderer.NewPostRenderer()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("it converts a single post into HTML", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := postRenderer.Render(&buf, aPost)
		
		if err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})

	t.Run("it renders an index of posts", func(t *testing.T) {
		buf := bytes.Buffer{}
		posts := []parser.Post{{Title: "Amazing title"}, {Title: "Another title"}}
	
		if err := postRenderer.RenderIndex(&buf, posts); err != nil {
			t.Fatal(err)
		}
	
		approvals.VerifyString(t, buf.String())
	})
}

func BenchmarkRender(b *testing.B) {
	var (
		aPost = parser.Post{
			Title:       "hello world",
			Body:        "This is a post",
			Description: "This is a description",
			Tags:        []string{"go", "tdd"},
		}
	)
	postRenderer, err := renderer.NewPostRenderer()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		postRenderer.Render(io.Discard, aPost)
	}
}