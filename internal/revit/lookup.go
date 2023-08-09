package revit

import (
	"net"
	"context"
)

// LookupAddr performs reverse DNS lookup for the given IP address using the specified resolver.
// It sends the lookup results (IP address, resolved DNS names, and any errors) to the results channel.
func LookupAddr(ip string, results chan<- LookupResult, sem chan struct{}, resolver string) {
	// Acquire a semaphore to limit the level of concurrency
	sem <- struct{}{}
	defer func() { <-sem }()

	// Initialize a net.Resolver for performing DNS lookups
	var r *net.Resolver
	if resolver != "" {
		// Create a custom dialer for the resolver
		dialer := &net.Dialer{}
		r = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return dialer.DialContext(ctx, "udp", resolver+":53")
			},
		}
	} else {
		r = net.DefaultResolver
	}

	// Perform a reverse DNS lookup using the Resolver
	names, err := r.LookupAddr(context.Background(), ip)

	// Send the lookup results (IP address, resolved DNS names, and errors) to the results channel
	results <- LookupResult{
		IPAddress: ip,
		DNSNames:  names,
		Error:     err,
	}
}
