package parser_test

import (
	"errors"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/emanuelquerty/bloggen/parser"
)

type StubFailingFS struct {

}

func (s StubFailingFS) Open(name string) (fs.File, error) {
	return nil, errors.New("oh no, i always fail")
}

func TestNewBlogPosts(t *testing.T) {
	t.Run("read file directory successfully", func(t *testing.T) {
		fs := fstest.MapFS{
			"hello world.md": {Data: []byte("Title: Post 1\nDescription: Description 1\nTags: tdd, go\n---\nHello\nWorld")},
			"hello-world2.md": {Data: []byte("Title: Post 2\nDescription: Description 2\nTags: rust, borrow-checker\n---\nB\nL\nM")},
		}
	
		posts, err := parser.NewPostsFromFS(fs)
	
		if err != nil {
			t.Error(err)
		}
	
		if len(posts) != len(fs) {
			t.Errorf("got %d posts but wanted %d posts", len(posts), len(fs))
		}

		assertPost(t, posts[0], parser.Post{
			Title: "Post 1", 
			Description: "Description 1", 
			Tags: []string{"tdd", "go"},
			Body: "Hello\nWorld",
		})
	})

	t.Run("fails to read file directory", func(t *testing.T) {
		var fs StubFailingFS
		_, err := parser.NewPostsFromFS(fs)
	
		if err == nil {
			t.Error(err)
		}
	})

	t.Run("replaces white space in titles", func(t *testing.T) {
		post := parser.Post{
			Title: "Post 1", 
			Description: "Description 1", 
			Tags: []string{"tdd", "go"},
			Body: "Hello\nWorld",
		}

		got := post.SanitizeTitle()
		want := "post-1"

		if got != want {
			t.Errorf("expected %q but got %q", want, got)
		}
	})
}

func assertPost(t *testing.T, got parser.Post, want parser.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}