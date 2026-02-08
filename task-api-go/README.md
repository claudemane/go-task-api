# Task API (Go, net/http)

Minimal in-memory task service that matches the assignment requirements.

## Run
```bash
go run ./cmd/api
```

Server starts on `:8080`.

## Auth
All endpoints require header:
- `X-API-KEY: secret12345`

## Examples (curl)

Create:
```bash
curl -i -X POST http://localhost:8080/tasks \
  -H "X-API-KEY: secret12345" \
  -H "Content-Type: application/json" \
  -d '{"title":"Write unit tests"}'
```

List:
```bash
curl -i http://localhost:8080/tasks -H "X-API-KEY: secret12345"
```

Get by id:
```bash
curl -i "http://localhost:8080/tasks?id=1" -H "X-API-KEY: secret12345"
```

Patch done:
```bash
curl -i -X PATCH "http://localhost:8080/tasks?id=1" \
  -H "X-API-KEY: secret12345" \
  -H "Content-Type: application/json" \
  -d '{"done":true}'
```
