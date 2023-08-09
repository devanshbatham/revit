package revit

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"math/rand"
	"time"
)

// ProcessStdin reads IP addresses from the standard input (stdin), performs DNS lookup
// for each IP, and sends the results to the results channel.
func ProcessStdin(results chan<- LookupResult, sem chan struct{}, wg *sync.WaitGroup, resolvers []string) {
	// Create a scanner to read input from the standard input (stdin)
	scanner := bufio.NewScanner(os.Stdin)

	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Loop through each line in the input
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

	// Check for errors while reading from stdin
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read from stdin: %v\n", err)
	}
}
