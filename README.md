# Savor

Savor is a HTTP middleware for filtering and sampling requests for Golang.
It provides a way to sample JSON requests based on the content of those requests.

Supported Features:
 * GJSON based filtering and message routing
 * Deduplication and count tracking of incoming messages.
 * Trigger channel notifications or method calls for matching filters.
 * Sample based on a combinarion of time, bucket size, modulos, or percentiles.

