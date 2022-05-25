Go event handler
=================

## Summary
HTTP server that accepts any POST request (JSON) from multiple clients' websites. Each request forms part of a struct (for that particular visitor) that will be printed to the terminal when the struct is fully complete. 

Frontend intentionally left blank.

### Example JSON Requests
```javascript
{
  "eventType": "copyAndPaste",
  "websiteUrl": "https://test.test",
  "sessionId": "123123-123123-123123123",
  "pasted": true,
  "formId": "inputCardNumber"
}

{
  "eventType": "timeTaken",
  "websiteUrl": "https://test.test",
  "sessionId": "123123-123123-123123123",
  "timeTaken": 12,
}

```

## Backend (Go)

Server is written using only the go base package. Whole thing is done to learn how to use the language


### Running the backend

Navigate to the `code-test` directory. Run `go run main.go`. Logging should indicate whether the server has started up correctly. The server is served on port 8000 by default. This can be changed by passing a different value to the `listenAddr` var in the first line of the main routine in `main.go`.

To run tests for the backend, run `go test` from the same directory.

### Benchmarking the server against many requests

Run `ab -c 100 -n 100 http://localhost:8000/` to run the concurrency benchmark on the server when it is being served on the 8000 port. This will send 100 requests to the server and provide information about response times, number of failed requests, etc. Running this command in my environment yielded no failures and maximum response times of \~7ms.



