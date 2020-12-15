
### Transports

A collection of useful `http.RoundTripper`s:

* `VerboseTransport` wraps an instance of a `http.RoundTripper` to write the `http.Request` and `http.Response`
* `TestDataTransport` enables returning the contents of local files to facilitate unit testing
* `SleepingTransport` enables testing clients which need to simulate latency in responses
* `RateLimitTransport` enables rate limit client requests by using a `golang.org/x/time/rate/Limiter` instance

See the test cases for how to use each of the `http.RoundTripper`s
