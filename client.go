package fastdns

// Client is used to send DNS requests to resolvers concurrently.
type Client interface {
	Resolve(domains []string, rrtype RRtype) []Record
	Close()
}

// Record contains a DNS record entry returned by Client.Resolve.
type Record struct {
	Question string
	Type     RRtype
	Answer   string
}

// New returns a Client that will respect the retry count, queries per seconds
// and a maximum number of concurrent queries that can happen at the same time.
// It is important to call Close on the client when done in order to free the UDP
// connections it creates.
func New(resolvers []string, retryCount int, queriesPerSecond int, maxConcurrency int) Client {
	servers := newRateLimitedServerList(resolvers, queriesPerSecond)
	balancer := newServerBalancer(servers)
	parser := &msgParser{}
	resolver := newResolverDNS(retryCount, balancer, parser)

	return newClientDNS(resolver, maxConcurrency)
}

// Resolver is used to resolve a DNS query and return a list of records.
// The process is responsible for selecting which DNS servers to use and
// for rate-limiting.
type Resolver interface {
	Resolve(query string, rrtype RRtype) []Record
	Close()
}
