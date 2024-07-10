package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/iagocanalejas/buzz-client/pkg"
)

func main() {
	argPath := os.Args[1]
	p, err := os.Stat(argPath)
	if err != nil && os.IsNotExist(err) {
		panic("file does not exist")
	}

	api := pkg.Init()
	folders, err := api.ListFolders()
	files, err := api.ListFiles(folders[0])
	fmt.Println(files[0].Name)

	if p.IsDir() {
		err := filepath.Walk(argPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				// TODO: handle file
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	} else {
		// TODO: handle file
	}
}
