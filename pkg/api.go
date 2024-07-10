package pkg

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

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

type Link struct {
	Href string
	Name string
}

func (a *API) ListFolders() ([]*Link, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account", a.url), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Cookie", fmt.Sprintf("session=%s", a.token))

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return extractLinks(doc), nil
}

func (a *API) ListFiles(folder *Link) ([]*Link, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", a.url, folder.Href), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Cookie", fmt.Sprintf("session=%s", a.token))

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return extractLinks(doc), nil
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
