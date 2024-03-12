FROM golang:1.16 as base
ENV GO111MODULE=on

WORKDIR /app

COPY ./src/go.mod ./src/go.sum ./src/
RUN cd src && go mod download

COPY . .

RUN cd src && go build -o main .

CMD [".src/main -prod"]

EXPOSE 8080
