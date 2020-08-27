FROM golang:1.14

WORKDIR /app

COPY . .

RUN go build ./cmd/phd

RUN ls -lah /app

CMD /app/phd
