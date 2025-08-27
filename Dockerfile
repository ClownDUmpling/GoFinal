FROM golang:1.23

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o todo-app main.go

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_PASSWORD=12345

CMD ["./todo-app"]