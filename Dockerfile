FROM golang:1.14

WORKDIR /app

COPY . .

RUN go build -o ./dist ./cmd/phd

RUN ls -lah /app/dist

CMD /app/dist/phd
