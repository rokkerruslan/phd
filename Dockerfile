FROM golang:1.15-alpine

WORKDIR /app

COPY . .

RUN go build ./cmd/phd

RUN ls -lah /app

CMD /app/phd
