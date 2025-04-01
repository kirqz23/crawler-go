package crawler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// Fetch downloads the content at the given URL and returns the raw HTML as a string.
func Fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching URL %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(bodyBytes), nil
}

// ParseLinks extracts all <a href="..."> links from the HTML data
func ParseLinks(htmlData string) ([]string, error) {
	links := []string{}

	doc, err := html.Parse(strings.NewReader(htmlData))
	if err != nil {
		return nil, err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links, nil
}

type CrawlTask struct {
	URL   string
	Depth int
}

func Worker(
	id int,
	tasks <-chan CrawlTask, // receive-only channel
	results chan<- []string, // send-only channel (for extracted links)
	done chan<- bool, // signal that this worker is done
) {
	for task := range tasks {
		htmlData, err := Fetch(task.URL)
		if err != nil {
			// Log error, continue
			continue
		}

		links, err := ParseLinks(htmlData)
		if err != nil {
			// Log error, continue
			continue
		}

		// Send extracted links back
		results <- links
	}
	// When the tasks channel closes, the worker finishes.
	done <- true
}
