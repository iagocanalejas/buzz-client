package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iagocanalejas/buzz-client/internal/api"
)

func main() {
	client := api.Init()

	var intoFolder string
	flag.StringVar(&intoFolder, "i", "", "folder to upload files into")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Usage: program -i <folder> <file-path>")
		os.Exit(1)
	}

	if intoFolder == "" {
		fmt.Println("no folder provided")
		os.Exit(1)
	}
	intoFolder = strings.ToLower(intoFolder)

	argPath := flag.Arg(0)
	p, err := os.Stat(argPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("file does not exist")
			os.Exit(1)
		}
		panic(err)
	}

	folder := retrieveFolderLinx(client, intoFolder)
	if folder == nil {
		// TODO: create new folder
		fmt.Println("folder does not exist in remote")
	}

	if p.IsDir() {
		err := filepath.Walk(argPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				client.Push(path, folder.ID())
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	} else {
		client.Push(argPath, folder.ID())
	}
}

func retrieveFolderLinx(client *api.API, name string) *api.Link {
	folders, err := client.List()
	if err != nil {
		fmt.Println("could not list folders")
		os.Exit(1)
	}

	for _, f := range folders {
		if strings.ToLower(f.Name) == name {
			return f
		}
	}
	return nil
}
