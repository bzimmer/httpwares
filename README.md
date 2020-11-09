### Transports

A collection of `http.RoundTripper`s useful for debugging and testing.

* `VerboseTransport` wraps an instance of a `http.RoundTripper` to write the `http.Request` and `http.Response`
* `TestDataTransport` enables returning the contents of local files to facilitate unit testing
