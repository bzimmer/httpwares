## HTTP Transports

[![build](https://github.com/bzimmer/httpwares/actions/workflows/build.yaml/badge.svg?branch=main)](https://github.com/bzimmer/httpwares/actions/workflows/build.yaml)
[![codecov](https://codecov.io/gh/bzimmer/httpwares/branch/master/graph/badge.svg?token=JBACLW92NN)](https://codecov.io/gh/bzimmer/httpwares)

A collection of useful `http.RoundTripper` implementations:

* `VerboseTransport` wraps an instance of a `http.RoundTripper` to write the
  `http.Request` and `http.Response`
* `RateLimitTransport` enables rate limit client requests by using a
  `golang.org/x/time/rate/Limiter` instance

See the test cases for how to use each `http.RoundTripper`
