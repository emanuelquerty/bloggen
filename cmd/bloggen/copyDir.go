package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)


func CopyDirFromFS(sourceDir fs.FS, destDirname string) {
	srcDirname := "assets"
	err := fs.WalkDir(sourceDir, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		relPath, err := filepath.Rel(srcDirname, path)
			if err != nil {
				return err
			}
		
		dstPath := filepath.Join(destDirname, relPath)

		if d.IsDir() {
			err = os.MkdirAll(dstPath, 0755)
			return err
		}
		err = copyFileFromFS(sourceDir, path, dstPath)
		if err != nil {
			return err
		}

		fmt.Printf("Copied: %s\n", relPath)
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
}

func copyFileFromFS(srcDir fs.FS, srcPath, dstPath string) error{
	sourceFile, err := srcDir.Open(srcPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
