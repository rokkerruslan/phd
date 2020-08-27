### Photographer

[![Actions Status](https://github.com/rokkerruslan/phd/workflows/Go/badge.svg)](https://github.com/rokkerruslan/phd/actions)

Backend for photographer search application.

### Quickstart

For start application you need `Docker` and `Go compiler` applications.

Build image
```shell script
$ docker build --tag phd .
```
Setup database:
```shell script
$ docker run --detach \
             --name phdb \
             --publish 5432:5432 \
             --env POSTGRES_PASSWORD=postgres \
             --volume phdb:/var/lib/postgresql/data \
             postgres:12
```
Setup migration:
```
to do
```


## Build on Docker
```shell script
docker build -t phd .
run --env-file .env phd
```

## API documentation

Open API spec [spec](docs/api.yml).
