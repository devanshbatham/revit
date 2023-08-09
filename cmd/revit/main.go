package main

import (
    "bufio"
    "flag"
    "fmt"
    "math/rand"
    "os"
    "sync"
    "strings"
    "time"
    "github.com/devanshbatham/revit/internal/revit"
    "github.com/fatih/color"
    "log"
)

func main() {
    // Command-line flags
    var (
        ips         = flag.String("i", "", "target IP to scan (-i IP_ADDRESS)")
        list        = flag.String("l", "", "target list of IPs to scan (-l INPUT_FILE)")
        concurrency = flag.Int("c", 10, "level of concurrency (-c 20)")
        resolvers   = flag.String("r", "", "resolvers for reverse DNS lookup (-r 8.8.8.8 or -r resolvers.txt)")
    )

    // Parse command-line flags
    flag.Parse()
    log.SetFlags(0)

    // Print tool banner
    log.Print(`
                       _ __ 
       ________ _   __(_) /_
      / ___/ _ \ | / / / __/
     / /  /  __/ |/ / / /_  
    /_/   \___/|___/_/\__/  
                         
            -  Reverse DNS Lookup Utility
 
`)

    // Process resolver list
    var resolverList []string
    if *resolvers != "" {
        if strings.Contains(*resolvers, ".txt") {
            resolverList = readResolversFromFile(*resolvers)
        } else {
            resolverList = []string{*resolvers}
        }
    }

    // Detect if there's piped input
    info, err := os.Stdin.Stat()
    isPipedInput := err == nil && info.Mode()&os.ModeNamedPipe != 0

    // Initialize wait group, results channel, and semaphore
    var wg sync.WaitGroup
    results := make(chan revit.LookupResult)
    sem := make(chan struct{}, *concurrency)

    rand.Seed(time.Now().UnixNano())

    // Perform DNS lookup
    doLookup := func(ip string) {
        selectedResolver := ""
        if len(resolverList) > 0 {
            selectedResolver = resolverList[rand.Intn(len(resolverList))]
        }
        revit.LookupAddr(ip, results, sem, selectedResolver)
    }

    // Handle command-line flags and piped input
    if *ips != "" {
        wg.Add(1)
        go func() {
            defer wg.Done()
            doLookup(*ips)
        }()
    } else if *list != "" {
        // Handle IP list from a file
        file, err := os.Open(*list)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
            return
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            ip := scanner.Text()
            wg.Add(1)
            go func(ipAddress string) {
                defer wg.Done()
                doLookup(ipAddress)
            }(ip)
        }
    } else if isPipedInput {
        // Handle piped input
        scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
            ip := scanner.Text()
            wg.Add(1)
            go func(ipAddress string) {
                defer wg.Done()
                doLookup(ipAddress)
            }(ip)
        }
    } else {
        flag.Usage() // Display usage if no flags provided
        return
    }

    // Wait for DNS lookups to complete
    go func() {
        wg.Wait()
        close(results)
    }()

    // Color formatting functions
    yellow := color.New(color.FgYellow).SprintFunc()
    green := color.New(color.FgGreen).SprintFunc()

    // Process and display lookup results
    for res := range results {
        if res.Error != nil && res.Error.Error() == "no such host" {
            continue
        }
        for _, name := range res.DNSNames {
            name = strings.TrimSuffix(name, ".")
            ipPadding := " "
            if len(res.IPAddress) < 15 {
                ipPadding = strings.Repeat(" ", 15-len(res.IPAddress))
            }

            // Print formatted IP address and resolved name
            fmt.Printf("%s%s [%s]\n", yellow(res.IPAddress), ipPadding, green(name))
        }
    }
}

// Read resolvers from a file
func readResolversFromFile(filename string) []string {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to open resolvers file: %v\n", err)
        return nil
    }
    defer file.Close()

    var resolvers []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        resolvers = append(resolvers, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to read resolvers file: %v\n", err)
    }

    return resolvers
}
