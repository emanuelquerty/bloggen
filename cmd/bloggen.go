package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/emanuelquerty/bloggen"
)

//go:embed assets
var assets embed.FS
const serverPort = "3000"

func main() {
	markdownDir := MardownDir{}
	flagSet := flag.NewFlagSet("bloggen", flag.ExitOnError)
	flagSet.Var(&markdownDir, "from", "the name of the directory that contains all markdown files")
	
	flagSet.Usage = func() {programUsage(os.Stdout)}
	err := flagSet.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}

	stat, err := os.Stat(markdownDir.path)
	if  errors.Is(err, os.ErrNotExist) || !stat.IsDir() {
		fmt.Printf("bloggen: %s: no such directory\n", markdownDir.path)
		return
	}

	sourceMarkdown := os.DirFS(markdownDir.path)
	parentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	destinationDirname := filepath.Join(parentDir, "build")
	switch os.Args[1] {
	case "create":
		buildSite(sourceMarkdown, destinationDirname)
	case "preview":
		preview(sourceMarkdown, destinationDirname)
	default:
		fmt.Printf("bloggen: %s: command not found \n\nTry 'bloggen --help' for more information.\n", os.Args[1])
	}
}

func buildSite(sourceMarkdown fs.FS, destinationDirname string) {
	// const colorGreen = "\033[0;32m"
	// const colorWhite = "\033[0;37m"
	
	fmt.Printf("\nScaffolding your website in %s\n", destinationDirname)
	bloggen.NewBlogSite(sourceMarkdown, destinationDirname)
	fmt.Printf("Build directory created.\n\n")

	// Adds css static assets
	destPath := filepath.Join(destinationDirname, "assets")

	fmt.Println("Downloading static assets...")
	CopyDirFromFS(assets, destPath)
	fmt.Printf("\nDone.\n")

}

func preview(sourceMarkdown fs.FS, destinationDirname string) {
	const colorRed = "\033[0;36m"
	const colorWhite = "\033[0;37m"

	buildSite(sourceMarkdown, destinationDirname)
	fmt.Printf("\n%sYour site is available for preview on %s http://localhost:%s/ \n\n",colorWhite, colorRed, serverPort)
	startServer(serverPort)
}

func startServer(portnumber string) {
	staticFileServer := http.FileServer(staticDir{http.Dir("build")})
	log.Fatal(http.ListenAndServe(":" + portnumber, staticFileServer))
}

func programUsage(w io.Writer) {
	usageText := `
Usage: bloggen [OPTION] --from [DIRECTORY]

Generates a site from a given DIRECTORY of markdowns.
DIRECTORY is the path to the folder that contains markdown files.

Available OPTIONS:
create      creates the site from DIRECTORY
preview     creates the site from DIRECTORY and starts a server to preview the generated site

Full documentation is available at <https://github.com/emanuelquerty/bloggen>`
	
	fmt.Fprintln(w, usageText)
}