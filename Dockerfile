FROM golang:1.22

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

COPY *.db ./
COPY cmd/*.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go-todo

ENV TODO_PORT="7540"
ENV TODO_DBFILE="../scheduler.db"
ENV TODO_PASSWORD="password"
ENV SECRET_KEY="my_secret_key"


EXPOSE 7540

CMD ["/go-todo"]