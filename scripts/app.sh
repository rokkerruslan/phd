#!/usr/bin/env sh

APP_NAME=photogram

# ==== app.sh =============================================== #

DAEMON_NAME=${APP_NAME}d

loc() {
    ENV_FILE_LOCAL=".env"

    if [ ! -f "$ENV_FILE_LOCAL" ]; then
        echo "Can't operate with local environment, file $ENV_FILE_LOCAL does"
        echo "not exist. Please copy example file and fill the variable values."
        echo "$ cp .env.example $ENV_FILE_LOCAL"
        exit 1
    fi

    ENV_LOCAL=$(cat $ENV_FILE_LOCAL | grep '^[a-zA-z]' | xargs)

    case $1 in
    t | test)
        tests
        ;;
    f | fmt)
        fmt
        ;;
    e | env)
        showEnv
        ;;
    s | start | "")
        start
        ;;
    *)
        echo command $1 does not exist, check help
        ;;
    esac
}

tests() {
    env $ENV_LOCAL go test -coverprofile=./dist/coverage.out ./internal/...
}

fmt() {
    go fmt ./...
}

start() {
    go build -o ./dist/$DAEMON_NAME ./cmd/$DAEMON_NAME

    env $ENV_LOCAL ./dist/$DAEMON_NAME
}

showEnv() {
    while IFS=' ' read -ra ARRAY; do
        echo "\nENVIRONMENT VARIABLES:\n"
        for i in "${ARRAY[@]}"; do
            echo "    $i\n"
        done
    done <<< "$ENV_LOCAL"
}

# ==== Entrypoint =========================================== #

USAGE=$(
    cat <<-END
app - entrypoint for control

USAGE
    ./app command subcommand

COMMANDS

    local       deal with local environment.

    SUBCOMMANDS

        env         show local environment variables
        fmt         format code (go fmt and goreturns)
        test        run tests
        start       start daemon

    help        print this docs

EXAMPLES

    $ ./app local test  # starts tests on local machine
END
)

if [ ! -d ".git" ]; then
    echo "Script must be starts from root project directory"
    echo "We detect .git directory for checking that)"
    echo
    exit 1
fi

case $1 in
l | local)
    loc "$2"
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
