
### Transports

A collection of `http.RoundTripper`s useful for debugging and testing.

* `VerboseTransport` wraps an instance of a `http.RoundTripper` to write the `http.Request` and `http.Response`
* `TestDataTransport` enables returning the contents of local files to facilitate unit testing
* `SleepingTransport` enables testing clients which need to simulate latency in responses
* `RateLimitTransport` works for the Strava API but is otherwise a WIP and needs to be generalized

See the test cases for how to use each of the `http.RoundTripper`s
