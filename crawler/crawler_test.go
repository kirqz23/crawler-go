package crawler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetch(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Hello, World!</body></html>"))
	}))
	defer server.Close()

	// Test Fetch function
	html, err := Fetch(server.URL)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	expected := "<html><body>Hello, World!</body></html>"
	if html != expected {
		t.Errorf("Expected %q, got %q", expected, html)
	}
}

func TestParseLinks(t *testing.T) {
	// Test HTML with links
	htmlData := `
        <html>
            <body>
                <a href="http://example.com">Example</a>
                <a href="/relative">Relative</a>
            </body>
        </html>
    `

	links, err := ParseLinks(htmlData)
	if err != nil {
		t.Fatalf("ParseLinks failed: %v", err)
	}

	expected := []string{"http://example.com", "/relative"}
	if len(links) != len(expected) {
		t.Fatalf("Expected %d links, got %d", len(expected), len(links))
	}

	for i, link := range links {
		if link != expected[i] {
			t.Errorf("Expected link %q, got %q", expected[i], link)
		}
	}
}

func TestWorker(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "page1") {
			w.Write([]byte(`<a href="/page2">Page 2</a>`))
		} else if strings.Contains(r.URL.Path, "page2") {
			w.Write([]byte(`<a href="/page3">Page 3</a>`))
		} else {
			w.Write([]byte(``))
		}
	}))
	defer server.Close()

	// Channels
	tasks := make(chan CrawlTask, 1)
	results := make(chan []string, 1)
	done := make(chan bool, 1)

	// Start the worker
	go Worker(1, tasks, results, done)

	// Send a task
	tasks <- CrawlTask{URL: server.URL + "/page1", Depth: 0}
	close(tasks)

	// Collect results
	var allLinks []string
	for {
		select {
		case links := <-results:
			allLinks = append(allLinks, links...)
		case <-done:
			// Worker finished
			if len(allLinks) != 1 || allLinks[0] != "/page2" {
				t.Errorf("Expected links [/page2], got %v", allLinks)
			}
			return
		}
	}
}
