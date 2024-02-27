FROM golang:1.21

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o /user_management
ENV PATH="/var/lib/authelia-main/users.yml"

EXPOSE 8080

CMD [ "/user_management", "-prod" ]
