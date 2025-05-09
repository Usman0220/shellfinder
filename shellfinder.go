package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	websitesFile  = "websites.txt"
	endpointsFile = "endpoints.txt"
	numWorkers    = 50 // Number of concurrent workers (goroutines)
)

// Task represents a single URL to check
type Task struct {
	Website  string
	Endpoint string
	FullURL  string
}

func main() {
	fmt.Println(`
   _____ __         ___________           __
  / ___// /_  ___  / / / ____(_)___  ____/ /__  _____
  \__ \/ __ \/ _ \/ / / /_  / / __ \/ __  / _ \/ ___/
 ___/ / / / /  __/ / / __/ / / / / / /_/ /  __/ /
/____/_/ /_/\___/_/_/_/   /_/_/ /_/\__,_/\___/_/

	`)
	fmt.Println("Shell Finder v0.0.1 by Ahmed Lekssays (0x70776e) - Go Version\n")
	fmt.Println("[*] INFO: Websites should start with 'http://' or 'https://' and should not end with '/'")

	fmt.Println("[*] INFO: Loading websites...")
	websites, err := loadLinesFromFile(websitesFile)
	if err != nil {
		fmt.Printf("[-] ERROR: Failed to load websites: %v\n", err)
		os.Exit(1)
	}
	if len(websites) == 0 {
		fmt.Println("[-] ERROR: websites.txt is empty or not found.")
		os.Exit(1)
	}
	fmt.Printf("[*] INFO: Loaded %d websites.\n", len(websites))

	fmt.Println("[*] INFO: Loading endpoints...")
	endpoints, err := loadLinesFromFile(endpointsFile)
	if err != nil {
		fmt.Printf("[-] ERROR: Failed to load endpoints: %v\n", err)
		os.Exit(1)
	}
	if len(endpoints) == 0 {
		fmt.Println("[-] ERROR: endpoints.txt is empty or not found.")
		os.Exit(1)
	}
	fmt.Printf("[*] INFO: Loaded %d endpoints.\n", len(endpoints))

	tasks := make(chan Task)
	results := make(chan string)
	var wg sync.WaitGroup

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second, // Adjust timeout as needed
	}

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(&wg, tasks, results, client)
	}

	// Goroutine to collect and print results
	var foundCount int
	var printWg sync.WaitGroup
	printWg.Add(1)
	go func() {
		defer printWg.Done()
		for result := range results {
			fmt.Println(result)
			if strings.HasPrefix(result, "[+] SUCCESS:") {
				foundCount++
			}
		}
	}()

	// Send tasks to workers
	fmt.Println("[*] INFO: Searching...")
	totalTasks := 0
	for _, website := range websites {
		cleanWebsite := strings.TrimSpace(website)
		if cleanWebsite == "" {
			continue
		}

		fmt.Printf("[*] INFO: Queuing tasks for target: %s\n", cleanWebsite)
		for _, endpoint := range endpoints {
			cleanEndpoint := strings.TrimSpace(strings.TrimLeft(endpoint, "/"))
			if cleanEndpoint == "" {
				continue
			}
			fullURL := cleanWebsite + "/" + cleanEndpoint
			tasks <- Task{Website: cleanWebsite, Endpoint: cleanEndpoint, FullURL: fullURL}
			totalTasks++
		}
	}
	close(tasks) // Signal workers that no more tasks will be sent

	wg.Wait()      // Wait for all workers to finish
	close(results) // Signal result printer that no more results will be sent
	printWg.Wait() // Wait for the printer to finish

	fmt.Printf("\n[*] INFO: Search complete. Checked %d potential URLs.\n", totalTasks)
	fmt.Printf("[*] INFO: Found %d potential shells.\n", foundCount)
}

// loadLinesFromFile reads lines from a file into a slice of strings.
func loadLinesFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" { // Skip empty lines
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

// worker processes tasks from the tasks channel and sends results to the results channel.
func worker(wg *sync.WaitGroup, tasks <-chan Task, results chan<- string, client *http.Client) {
	defer wg.Done()
	for task := range tasks {
		req, err := http.NewRequest("GET", task.FullURL, nil)
		if err != nil {
			// Silently skip malformed URLs or log them differently if needed
			continue
		}
		req.Header.Set("User-Agent", "ShellFinderGo/0.0.1 (github.com/0x70776e)")

		resp, err := client.Do(req)
		if err != nil {
			// Connection errors are common, so not printing them by default to reduce noise
			// results <- fmt.Sprintf("[-] WARN: Failed to connect to %s: %v", task.FullURL, err)
			continue
		}

		io.Copy(io.Discard, resp.Body) // Read the body to EOF to allow connection reuse
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK { // 200 OK
			results <- fmt.Sprintf("[+] SUCCESS: Shell found at: %s (Endpoint: %s)", task.FullURL, task.Endpoint)
		} else if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized { // 403 or 401
			// You could also include the endpoint here if desired for these statuses
			// results <- fmt.Sprintf("[!] INFO: Potential protection at %s (Status: %d, Endpoint: %s)", task.FullURL, resp.StatusCode, task.Endpoint)
		}
		// Add more status code checks if needed
	}
}
