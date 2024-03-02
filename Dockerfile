FROM golang:1.16 as base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

FROM base as development
RUN go build -o main .

CMD ["./main"]

EXPOSE 8080
