#!/usr/bin/env sh

# TODO: add goda tool
# goda graph -cluster  ./cmd/phd | dot -Tsvg -o ./dist/graph.svg && open ./dist/graph.svg

APP_NAME=ph

# ==== app.sh =============================================== #

DAEMON_NAME=${APP_NAME}d

tests() {
    go test ./internal/...
}

build() {
    go build -o ./dist/$DAEMON_NAME ./cmd/$DAEMON_NAME
    go build -o ./dist/mi ./cmd/mi
    chmod +x ./dist/mi
}

cover() {
    go test -coverprofile=./dist/coverage.out ./internal/...
}

fmt() {
    go fmt ./...
}

start() {
    go build -o ./dist/$DAEMON_NAME ./cmd/$DAEMON_NAME

    # shellcheck disable=SC2086
    env $ENV_LOCAL ./dist/$DAEMON_NAME
}

migrate() {
    # shellcheck disable=SC2086
    env $ENV_LOCAL ./dist/mi up
}

update() {
    go get -u ./...
    go mod vendor
}

showEnv() {
    echo "$ENV_LOCAL"
}

# ==== Entry point ========================================== #

USAGE=$(
    cat <<-END
app - entry point for control

USAGE
    ./app command

COMMANDS

    env         show local environment variables
    fmt         format code (go fmt and goreturns)
    test        run tests
    start       start
    migrate     apply migrations
    update      update all deps (minor/patch) and vendor it

    help        print this docs

EXAMPLES

    $ ./app test  # starts tests on local machine
END
)

if [ ! -d ".git" ]; then
    echo "Script must be starts from root project directory"
    echo "We detect .git directory for checking that)"
    echo
    exit 1
fi

ENV_FILE_LOCAL=".env"

if [ ! -f "$ENV_FILE_LOCAL" ]; then
    echo "Can't operate with local environment, file $ENV_FILE_LOCAL does"
    echo "not exist. Please copy example file and fill the variable values."
    echo "$ cp .env.example $ENV_FILE_LOCAL"
    exit 1
fi

ENV_LOCAL=$(cat $ENV_FILE_LOCAL | grep '^[a-zA-Z]' | xargs)


case $1 in
t | test)
    tests
    ;;
b | build)
    build
    ;;
f | fmt)
    fmt
    ;;
e | env)
    showEnv
    ;;
s | start)
    start
    ;;
m | mi | migrate)
    migrate
    ;;
u | update)
    update
    ;;
h | help | "")
    echo "$USAGE"
    ;;
*)
    echo "$1 command does not exist, check help"
    echo
    exit 1
    ;;
esac
