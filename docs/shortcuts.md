# Shortcuts

Event list:
```shell script
$ curl -s localhost:3000/api/v1/events/ | jq
```

Create Event
```shell script
echo '{"Name":"Event X","OwnerID": 1,"Timelines":[{"Start":"2006-01-02T15:05:05Z","End":"2006-01-02T16:06:05Z"}],"Point":{"Lt":1.1,"Ln":2.2}}' | curl -X POST -d @- -s localhost:3000/api/v1/events/ | jq
```

Update Event
```shell script
$ curl -X PUT -d '{"ID":8,"Name":"Event X8","Timelines":[{"Start":"2006-01-02T15:05:05Z","End":"2006-01-02T16:06:05Z"}]}' -s localhost:3000/api/v1/events/ | jq
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

SignUp
```shell script
$ curl -X POST -d '{"Email":"em","Password":"1234567890"}' http://localhost:3000/api/v1/auth/sign-up
```
