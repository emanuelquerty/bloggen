package main

import (
	"net/http"
	"os"
)

// MarkdownDir implements the Value interface from the flag package
// This is used to hold the value of the flag defined by a flagset
type MardownDir struct {
	path string
}

func (m *MardownDir) String() string {
	return m.path
}

func (m *MardownDir) Set(path string) error {
	m.path = path
	return nil
}

// HTMLDir implements the FileSystem interface from net/http
type staticDir struct {
	d http.Dir
}

func (d staticDir) Open(name string) (http.File, error) {
	// Resolves path for urls enabling routing without .html
	f, err := d.d.Open(name)
	if os.IsNotExist(err) {
		// Not found, try with .html
		if f, err := d.d.Open(name + ".html"); err == nil {
			return f, nil
		}

		// Resolves path for static assets when within posts folder
		if f, err := d.d.Open(name[len("/posts/"):]); err == nil {
			return f, nil
		}

	}
	return f, err
}

