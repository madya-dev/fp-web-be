FROM golang:1.20.4

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod tidy

CMD ["go", "run", "cmd/main/main.go"]