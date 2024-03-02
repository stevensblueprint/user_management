FROM golang:1.16 as base
ENV GO111MODULE=on

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main -prod"]

EXPOSE 8080
