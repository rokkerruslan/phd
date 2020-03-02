# Shortcuts

Event list:
```shell script
$ curl -s localhost:3000/api/v1/events/ | jq
```

Create Event
```shell script
$ curl -X POST -d '{"Name":"Event X","Timelines":[{"Start":"2006-01-02T15:05:05Z","End":"2006-01-02T16:06:05Z"}],"Point":{"Lt":1.1,"Ln":2.2}}' -s localhost:3000/api/v1/events/ | jq
```

Offer list:
```shell script
$ curl -s localhost:3000/api/v1/offers/ | jq
```

Create offer
```shell script
$ curl -X POST -d '{"AccountID":1,"EventID":1}' -s localhost:3000/api/v1/offers/ | jq
```

Code coverage:
```shell script
$ go test -coverprofile=coverage.out ./...
$ go tool cover -func=coverage.out
```
