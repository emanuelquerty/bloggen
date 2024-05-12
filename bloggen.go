package bloggen

import (
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/emanuelquerty/bloggen/parser"
	"github.com/emanuelquerty/bloggen/renderer"
)

func NewBlogSite(sourceMarkdown fs.FS, destinationDirname string) {
	posts, err := parser.NewPostsFromFS(sourceMarkdown)
	if err != nil {
		log.Fatal(err)
	}

	renderer, err := renderer.NewPostRenderer()
	if err != nil {
		log.Fatal(err)
	}

	createDirectory(destinationDirname)
	file := createFile(destinationDirname, "index")

	err = renderer.RenderIndex(file, posts)
	if err != nil {
		log.Fatal(err)
	}

	postDirectory := filepath.Join(destinationDirname, "posts")
	createDirectory(postDirectory)
	for _, post := range posts {
		pageName := strings.ReplaceAll(post.Title, " ", "-")
		file := createFile(postDirectory, pageName)
		renderer.Render(file, post)
	}
	rewriteURL(destinationDirname)
}

func createDirectory(dirName string) {
	err := os.MkdirAll(dirName, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}

func createFile(dirName, filename string) io.Writer {
	newFilePath := filepath.Join(dirName, filename + ".html")
	file, err := os.Create(newFilePath)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func rewriteURL(dirName string) {
	bytes := []byte(
`RewriteEngine on

# Remove .html (or htm) from visible URL (permanent redirect)
RewriteCond %{REQUEST_URI} ^/(.+)\.html?$ [nocase]
RewriteRule ^ /%1 [L,R=301]

# Quietly point back to the HTML file (temporary/undefined redirect):
RewriteCond %{REQUEST_FILENAME} !-d
RewriteCond %{REQUEST_FILENAME}.html -f
RewriteRule ^ %{REQUEST_URI}.html [END]`)

	urlRewritePath := filepath.Join(dirName, ".htaccess")
	err := os.WriteFile(urlRewritePath, bytes, 0755)
	if err != nil {
		log.Fatal(err)
	}
}