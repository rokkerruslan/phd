# Shortcuts

Database setup:
```shell script
$ docker run --name ph -p 5432:5432 -v ph:/var/lib/postgresql/data -d postgres:12 -c log_statement=all
```

Event list:
```shell script
$ curl -s localhost:3000/events/ | jq
```

Create Event
```shell script
$ curl -X POST -d @.i/event.json -s localhost:3000/events | jq
```

Update Event
```shell script
curl -X PUT -d @.i/event.json -s localhost:3000/events/2 | jq
```

Offer list:
```shell script
$ curl -s localhost:3000/offers/?account_id=5 | jq
```

Create offer
```shell script
$ curl -X POST -d @.i/offerCreate.json -s localhost:3000/offers/
```

Code coverage:
```shell script
$ go test -coverprofile=coverage.out ./...
$ go tool cover -func=coverage.out
```

SignUp
```shell script
$ curl -X POST -d @.i/account.json localhost:3000/accounts/sign-up
```
