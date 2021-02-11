# external-dns-endpoint-server

A tiny server providing a connector source implementation for
the [`external-dns`](https://github.com/kubernetes-sigs/external-dns).

It also offers an HTTP server that accepts IP address in plaintext via a POST
request body. The IP address received is then stored internally, and sent over
to the `external-dns` when it connects to connector source server implemented
at this app.
