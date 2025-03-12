FROM golang:1.23

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

EXPOSE 7540

CMD ["/main"]