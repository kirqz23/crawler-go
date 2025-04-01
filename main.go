package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kirqz23/crawler-go/crawler"
)

func main() {
	startURL := flag.String("url", "", "The start URL to crawl")
	maxDepth := flag.Int("depth", 2, "Maximum depth to crawl")
	workers := flag.Int("workers", 10, "Number of concurrent workers")

	flag.Parse()

	if *startURL == "" {
		fmt.Println("Please provide a start URL using --url")
		os.Exit(1)
	}

	// For now, just print out the arguments
	log.Printf("Start URL: %s, Max Depth: %d, Workers: %d\n", *startURL, *maxDepth, *workers)
	// Later, call your crawler logic from here

	// Channels
	tasks := make(chan crawler.CrawlTask)
	results := make(chan []string)
	done := make(chan bool)

	// Launch workers
	for i := 0; i < *workers; i++ {
		go crawler.Worker(i, tasks, results, done)
	}

	// Start by sending the initial task
	go func() {
		tasks <- crawler.CrawlTask{URL: *startURL, Depth: 0}
	}()

	visited := make(map[string]bool)
	visited[*startURL] = true
	activeWorkers := *workers // how many worker goroutines we launched
	depthMap := make(map[string]int)
	depthMap[*startURL] = 0

	for {
		select {
		case links := <-results:
			// We got new links from a worker
			for _, link := range links {
				// Possibly normalize the link or skip non-HTTP links
				// Also consider absolute vs relative
				// For simplicity, assume link is absolute and HTTP
				if !visited[link] {
					// Check if we should crawl this link based on max depth
					currentDepth, exists := depthMap[link]
					if !exists {
						// If depth is not set, this is a new link
						// We should implement proper depth tracking from the parent URL
						currentDepth = 0 // placeholder, needs proper implementation
					}

					if currentDepth < *maxDepth {
						visited[link] = true
						depthMap[link] = currentDepth + 1
						tasks <- crawler.CrawlTask{URL: link, Depth: currentDepth + 1}
					}
					// Actually we need the depth from the *origin* link
					// So we need a slight tweak:
					// - pass the current depth from the worker
					// - store it in depthMap
					// - then linkDepth = currentDepth + 1
				}
			}

		case <-done:
			activeWorkers--
			if activeWorkers == 0 {
				// All workers done
				fmt.Println("All workers have finished.")
				close(results)
				return
			}
		}
	}
}
