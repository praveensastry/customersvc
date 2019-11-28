# customersvc

This example demonstrates how to use Go kit to implement a REST-y HTTP service.
It leverages the excellent [gorilla mux package](https://github.com/gorilla/mux) for routing.

Run the example with the optional port address for the service:

```bash
$ go run ./cmd/customersvc/main.go -http.addr :8080
ts=2018-05-01T16:13:12.849086255Z caller=main.go:47 transport=HTTP addr=:8080
```

Create a Customer:

```bash
curl -d '{"id":"1234","Name":"Go Kit"}' -H "Content-Type: application/json" -X POST http://localhost:8080/customers/
{}
```

Get the customer you just created

```bash
curl localhost:8080/customers/1234
{"customer":{"id":"1234","name":"Go Kit"}}
```
