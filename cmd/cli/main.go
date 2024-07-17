package main

import (
	"flag"
	"fmt"
	"log"
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
		log.Println("Usage: program -i <folder> <file-path>")
		os.Exit(1)
	}

	if intoFolder == "" {
		log.Println("no folder provided")
		os.Exit(1)
	}
	intoFolder = strings.ToLower(intoFolder)

	argPath := flag.Arg(0)
	p, err := os.Stat(argPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("file does not exist")
			os.Exit(1)
		}
		panic(err)
	}

	var folder *api.Link
	folder = retrieveFolderLink(client, intoFolder)
	if folder == nil {
		folder, err = client.CreateFolder(intoFolder)
		if err != nil || folder == nil {
			log.Println(err)
			os.Exit(1)
		}
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

func retrieveFolderLink(client *api.API, name string) *api.Link {
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
