# HTTP Concurrent Request

## Requirement
- A web service provides weather API: `GET http://example.server/weather?date=2021-06-29`
- query parameter:
    - date: string of date
- response data is of JSON format:
```
    {
      "status": 0,
      "result": {
        "high": 30,
        "low":  20,
        "rain": 0.75
      }
    }
```
Unfortunately, it does not provide APIs to get a range of data.
Please complete func quote() so it can take year and month as arguments,
returning a slice of results.

- for incomplete data, return error with empty slice

### Bonus
sending requests concurrency

## Ref
- [Mocking HTTP Call in Golang a Better Way](https://clavinjune.dev/en/blogs/mocking-http-call-in-golang-a-better-way/)
