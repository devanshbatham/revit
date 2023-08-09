package revit

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"math/rand"
	"time"
)

// ProcessFile reads IP addresses from a file, performs DNS lookup for each IP,
// and sends the results to the results channel.
func ProcessFile(filename string, results chan<- LookupResult, sem chan struct{}, wg *sync.WaitGroup, resolvers []string) {
	// Open the specified file for reading IP addresses
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Create a scanner to read IP addresses from the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := scanner.Text()

		// Increment the wait group count for each IP address
		wg.Add(1)

		// Launch a goroutine to perform DNS lookup for the IP address
		go func(ipAddress string) {
			defer wg.Done()

			// Select a random resolver from the provided list
			selectedResolver := ""
			if len(resolvers) > 0 {
				selectedResolver = resolvers[rand.Intn(len(resolvers))]
			}

			// Call the LookupAddr function to perform reverse DNS lookup
			LookupAddr(ipAddress, results, sem, selectedResolver)
		}(ip)
	}

	// Check for errors while scanning the file
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
	}
}
