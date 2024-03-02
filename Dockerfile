FROM golang:1.16 as base

COPY . .

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main -dev"]