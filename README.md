### Photographer

Backend for photographer search application.

### Quickstart

For start applicatino you need `Docker` and `Go compiler` applications.

Setup database:
```shell script
$ docker run --detach \
             --name phdb \
             --publish 5432:5432 \
             --volume phdb:/var/lib/postgresql/data \
             postgres:12
```

Build application and tools:
```shell script
$ ./app build
```

Setup table structure:
```
$ ./app migrate
```

Run application:
```
$ ./app start
```

## API documentation

Open API spec [spec](docs/api.yml).
