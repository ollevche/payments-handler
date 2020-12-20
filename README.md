# API Design Example
A simple handler which represents extendable API design in Go.

## Requirements
Go runtime with version 1.15 or higher.

## How to run the server?
In the root directory run the following commands:
```
go build
./payments-handler serve
```

You should see something like this:
```
10:48PM INF cmd/root.go:36 > Log level set to debug
10:48PM INF cmd/serve.go:54 > Listening and serving HTTP requests... address=:8080
```

## How to run tests?
In the root directory run the following command:
```
go test ./http/rest
```

## Payments handler
The API provides a single HTTP endpoint by the following path: `/payments/options`. It expectes to read `product_id` from query params. It's possible to import `insomnia.json` and send requests via [Insomnia tool](https://insomnia.rest/). Note: third-party dependencies are emulated.

Example request-response via curl (with beatified JSON):
```
curl -v http://localhost:8080/payments/options?product_id=5fdf8dd752da22ffcb1cf412
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /payments/options?product_id=5fdf8dd752da22ffcb1cf412 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json
< X-Response-Time: 49
< X-Server-Name: r5u
< Date: Sun, 20 Dec 2020 20:56:16 GMT
< Content-Length: 422
< 
{
   "product":{
      "id":"5fdf8dd752da22ffcb1cf412"
   },
   "options":[
      {
         "provider":{
            "id":0,
            "name":"PayPal"
         },
         "button_url":"https://payments.paypal.com/button"
      },
      {
         "provider":{
            "id":1,
            "name":"Apple Pay"
         },
         "button_url":"https://payments.applepay.com/button"
      },
      {
         "provider":{
            "id":2,
            "name":"Google Pay"
         },
         "button_url":"https://payments.googlepay.com/button"
      },
      {
         "provider":{
            "id":3,
            "name":"Stripe"
         },
         "button_url":"https://payments.stripe.com/button"
      }
   ]
}
* Connection #0 to host localhost left intact

```
