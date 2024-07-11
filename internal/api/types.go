package api

import (
	"fmt"
	"strings"
)

type Link struct {
	Href string
	Name string
}

func (l *Link) ID() string {
	parts := strings.Split(l.Href, "/")
	return parts[len(parts)-1]
}

type File struct {
	ID string `json:"id"`
}

type ProgressWriter struct {
	Total      int64
	Written    int64
	Percentage int
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.Written += int64(n)
	newPercentage := int(float64(pw.Written) / float64(pw.Total) * 100)
	if newPercentage != pw.Percentage {
		pw.Percentage = newPercentage
		fmt.Printf("\rUploading... %d%% complete", pw.Percentage)
	}
	return n, nil
}
