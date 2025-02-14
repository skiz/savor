# Savor

Savor is a HTTP middleware for filtering and sampling requests for Golang.
It provides a way to sample JSON requests based on the content of those requests.

Supported Features:
 * GJSON based filtering and message routing
 * Deduplication and count tracking of incoming messages.
 * Trigger channel notifications or method calls for matching filters.
 * Sample based on a combinarion of time, bucket size, modulos, or percentiles.


Maybe this should be a router?

```


// implement a filter
type Filter interface {
    func matcher() string
    func 
}

config := &savor.Config{
    filters: []savor.Filter{
        matcher: "
    }
}

s := savor.NewSavor(config)


client := &http.Client{
	CheckRedirect: redirectPolicyFunc,
}

```


