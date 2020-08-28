FROM golang:1.15-alpine

WORKDIR /app

COPY . .

RUN go build ./pkg/mi

CMD /app/mi up
