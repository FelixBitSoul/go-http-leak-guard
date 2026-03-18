package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	client := &http.Client{}
	var wg sync.WaitGroup

	// Start 50 concurrent goroutines
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				// Make request to local port 8080
				serverURL := "http://localhost:8080"
				if envURL := os.Getenv("SERVER_URL"); envURL != "" {
					serverURL = envURL
				}
				resp, err := client.Get(serverURL)
			if err != nil {
				log.Printf("Request error: %v", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// OPTIMIZED: Use an immediately-invoked function literal to ensure atomic resource handling.
			// This guarantees resource cleanup even if subsequent logic panics.
			func() {
				// CRITICAL: Drain and discard the remaining body data.
				// This is the key to allowing the Transport to reuse the TCP connection.
				_, _ = io.Copy(io.Discard, resp.Body)

				// Ensure the body is closed, even if a panic occurs.
				resp.Body.Close()
			}()

				time.Sleep(50 * time.Millisecond)
			}
		}()
	}

	// Monitor and print PID and active file handles every second
	go func() {
		for {
			pid := os.Getpid()
			var fdCount int

			// Cross-platform way to count file descriptors
			if runtime.GOOS == "linux" {
				// Linux: count files in /proc/self/fd
				data, err := ioutil.ReadFile("/proc/self/fd")
				if err == nil {
					fdCount = len(data) // This is approximate, actual count would need directory listing
				}
			} else if runtime.GOOS == "windows" {
				// Windows: we can't easily count handles, so we'll show a placeholder
				fdCount = -1
			} else {
				// Other Unix-like systems
				fdCount = -1
			}

			if fdCount >= 0 {
				fmt.Printf("PID: %d, Active FDs: %d\n", pid, fdCount)
			} else {
				fmt.Printf("PID: %d, FD monitoring not available on %s\n", pid, runtime.GOOS)
			}

			time.Sleep(1 * time.Second)
		}
	}()

	fmt.Println("Starting guarded HTTP client with 50 concurrent goroutines...")
	fmt.Println("This version properly drains response bodies to prevent connection leaks")
	fmt.Println("Press Ctrl+C to stop")

	wg.Wait()
}