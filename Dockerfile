FROM golang:1.22

USER root

WORKDIR /app

COPY . .

ENV TODO_PORT="7540"
ENV TODO_DBFILE="/scheduler.db"
ENV TODO_PASSWORD="password"
ENV SECRET_KEY="my_secret_key"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go-todo

EXPOSE 7540

CMD ["/go-todo"]