package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/net/html"

	"github.com/joho/godotenv"
)

type API struct {
	url   string
	token string

	client *http.Client
}

func Init() *API {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	return &API{
		url:    os.Getenv("BUZZHEAVIER_URL"),
		token:  os.Getenv("BUZZHEAVIER_API_KEY"),
		client: &http.Client{},
	}
}

func (a *API) List() ([]*Link, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account", a.url), nil)
	log.Println(req.URL)
	if err != nil {
		log.Fatalf("failed to create request: %s", err)
		return nil, err
	}

	req.Header.Add("Cookie", fmt.Sprintf("session=%s", a.token))

	resp, err := a.client.Do(req)
	if err != nil {
		log.Fatalf("failed to do request: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %s", err)
		return nil, err
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("failed to parse response body: %s", err)
		return nil, err
	}

	return extractLinks(doc), nil
}

func (a *API) ListFolder(folder *Link) ([]*Link, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", a.url, folder.Href), nil)
	log.Println(req.URL)
	if err != nil {
		log.Fatalf("failed to create request: %s", err)
		return nil, err
	}

	req.Header.Add("Cookie", fmt.Sprintf("session=%s", a.token))

	resp, err := a.client.Do(req)
	if err != nil {
		log.Fatalf("failed to do request: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %s", err)
		return nil, err
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("failed to parse response body: %s", err)
		return nil, err
	}

	return extractLinks(doc), nil
}

func (a *API) Push(filePath, folderID string) (*File, error) {
	fileName := filepath.Base(filePath)
	url := fmt.Sprintf("https://w.buzzheavier.com/%s?folderId=%s", fileName, folderID)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("failed to get file info: %s", err)
		return nil, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		log.Fatalf("failed to create form file: %s", err)
		return nil, err
	}

	pw := &ProgressWriter{Total: fileInfo.Size()}
	teeReader := io.TeeReader(file, pw)

	_, err = io.Copy(part, teeReader)
	if err != nil {
		log.Fatalf("failed to copy file content: %s", err)
		return nil, err
	}

	// Close the writer to finalize the multipart form data
	err = writer.Close()
	if err != nil {
		log.Fatalf("failed to close writer: %s", err)
		return nil, err
	}

	// Create the HTTP request
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		log.Fatalf("failed to create request: %s", err)
		return nil, err
	}

	// Set the Authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the HTTP request
	resp, err := a.client.Do(req)
	if err != nil {
		log.Fatalf("failed to make request: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Fatalf("failed to upload file: %s", resp.Status)
		return nil, err
	}

	var uploadResp File
	err = json.NewDecoder(resp.Body).Decode(&uploadResp)
	if err != nil {
		log.Fatalf("failed to decode response: %s", err)
		return nil, err
	}

	return &uploadResp, nil
}

func extractLinks(n *html.Node) []*Link {
	var folders []*Link = make([]*Link, 0)
	var f func(*html.Node, bool)
	f = func(n *html.Node, inTable bool) {
		if n.Type == html.ElementNode {
			if n.Data == "table" {
				inTable = true
			}
			if inTable && n.Data == "a" {
				var href, text string
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						href = attr.Val
						break
					}
				}
				if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
					text = n.FirstChild.Data
				}
				if href != "" && text != "" {
					folders = append(folders, &Link{Href: href, Name: text})
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, inTable)
		}
	}
	f(n, false)
	return folders
}
