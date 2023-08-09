package revit

// LookupResult represents the result of a reverse DNS lookup.
type LookupResult struct {
	IPAddress string   // The IP address for which the lookup was performed
	DNSNames  []string // List of resolved DNS names associated with the IP address
	Error     error    // Any error that occurred during the lookup process
}
