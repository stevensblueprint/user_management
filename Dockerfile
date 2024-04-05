FROM golang:latest as base

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN go build

CMD ["./user_management"]

EXPOSE 8080
