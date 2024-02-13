FROM golang:1.16 as base

FROM base as dev

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main -dev"]