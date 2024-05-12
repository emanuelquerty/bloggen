package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"strings"
)

const (
	titleSeparator = "Title: "
	descriptionSeparator = "Description: "
	tagsSeparator = "Tags: "
)

type Post struct {
	Title string
	Description string
	Tags []string
	Body string
}

func (p Post) SanitizeTitle() string{
	return strings.ToLower(strings.Replace(p.Title, " ", "-", -1))
}

func NewPostsFromFS(fileSystem fs.FS) ([]Post, error) {
	dir, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		return nil, err
	}

	var posts []Post
	for _, f := range dir {
		post, err := getPost(fileSystem, f.Name())
		if err != nil {
			return nil, err //todo: needs clarification, should we totally fail if one file fails? or just ignore?
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func getPost(fileSystem fs.FS, filename string) (Post, error) {
	postfile, err := fileSystem.Open(filename)
	if err != nil {
		return Post{}, err
	}
	defer postfile.Close()

	return newPost(postfile)
}

func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	readLine := func() string {
		scanner.Scan()
		return scanner.Text()
	}

	title := readLine()[len(titleSeparator):]
	description := readLine()[len(descriptionSeparator):]
	tags := strings.Split(readLine()[len(tagsSeparator):], ", ")
	body := readBody(scanner)

	post := Post{Title: title, Description: description, Tags: tags, Body: body}

	return post, nil
}

func readBody(scanner *bufio.Scanner) string {
	scanner.Scan() // ignore a line

	buff := bytes.Buffer{}
	for scanner.Scan() {
		fmt.Fprintln(&buff, scanner.Text())
	}

	return strings.TrimSuffix(buff.String(), "\n")
}