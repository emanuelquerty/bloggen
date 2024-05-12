package main

import (
	"embed"
	"io/fs"
	"log"

	"github.com/emanuelquerty/bloggen/parser"
)

//go:embed markdown
var markdown embed.FS

func main() {
	markdownDir, err := fs.Sub(markdown, "markdown")
	if err != nil {
		log.Fatal(err)
	}

	posts, err := parser.NewPostsFromFS(markdownDir)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", posts)
}