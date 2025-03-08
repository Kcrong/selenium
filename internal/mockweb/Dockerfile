FROM golang:1.24-alpine

WORKDIR /app

COPY server.go main.go

COPY htmls htmls

RUN go build -o testserver main.go

CMD ["./testserver"]
