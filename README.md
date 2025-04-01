# crawler-go

A **Concurrent Web Crawler** written in **Golang**. This project demonstrates how to build a scalable and efficient web crawler using Goroutines and channels for concurrency.

## Features

- Fetches HTML content from URLs.
- Extracts all `<a href="...">` links from the HTML.
- Supports concurrent crawling with configurable worker count.
- Depth-limited crawling to avoid infinite loops.
- Graceful handling of errors during HTTP requests or HTML parsing.

---

## Getting Started

### Prerequisites

- **Go 1.20+** installed on your system.
- Internet connection to fetch external URLs.

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/kirqz23/crawler-go.git
   cd crawler-go
   ```

2. Initialize the Go module (if not already done):

    ```bash
    go mod tidy
    ```

## Usage

### Running the Crawler
You can run the crawler using the main.go file. Use the following command:

    ```go
    go run [main.go](main.go) --url=<url> --depth=<max-depth> --workers=<worker-count>
    ```

#### Command-Line Flags
* `--url`: The starting URL to crawl (required).
* `--depth`: Maximum depth to crawl (default: 2).
* `--workers`: Number of concurrent workers (default: 10).

#### Example

    ```go
    go run [main.go](main.go) --url=https://example.com --depth=3 --workers=5
    ```

## Project Structure
```
crawler-go/
├── crawler/
│   ├── crawler.go         # Core crawler logic (Fetch, ParseLinks, Worker)
│   ├── crawler_test.go    # Unit tests for crawler logic
├── main.go                # Entry point for the application
├── go.mod                 # Go module file
└── README.md              # Project documentation
```

## Code Overview

### Core Functions
1. `Fetch(url string) (string, error)`
    * Fetches the HTML content of a given URL.
    * Handles HTTP errors and non-200 status codes.

2. `ParseLinks(htmlData string) ([]string, error)`
    * Parses the HTML content and extracts all <a href="..."> links.

3. `Worker(id int, tasks <-chan CrawlTask, results chan<- []string, done chan<- bool)`
    * Processes crawl tasks concurrently.
    * Fetches HTML, extracts links, and sends results back.

### Data Structures
* `CrawlTask`: Represents a single crawling task with a URL and depth.

## Testing
Unit tests are provided for the core functionality in crawler_test.go. To run the tests:

    ```bash
    go test go test -v ./... 
    ```