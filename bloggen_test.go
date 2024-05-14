package bloggen_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/emanuelquerty/bloggen"
)

func TestNewBlogSite(t *testing.T) {
	markdownFS := os.DirFS("cmd/markdown")
	blogDirectory := "testbuild"
	bloggen.NewBlogSite(markdownFS, blogDirectory)

	want := []string{
		"index.html",
		"posts/Interesting-title.html",
		"posts/Placeholder-for-First-Blog-Post.html",
	}

	got := Files(os.DirFS(blogDirectory))

	slices.Sort(want)
	slices.Sort(got)
	if !slices.Equal(want, got) {
		t.Errorf("expected %+v but got %+v", want, got)
	}
}

func Files(fsys fs.FS) (paths []string) {
    fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
        if filepath.Ext(p) == ".html" {
            paths = append(paths, p)
        }
        return nil
    })
    return paths
}