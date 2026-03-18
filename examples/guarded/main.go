package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/FelixBitSoul/go-http-leak-guard/utils"
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

			// Cross-platform way to count file descriptors
			fdCount := utils.GetFDCount()

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
