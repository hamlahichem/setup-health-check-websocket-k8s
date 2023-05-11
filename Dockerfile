FROM golang:1.20

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY ws_server.go .

CMD go run ws_server.go