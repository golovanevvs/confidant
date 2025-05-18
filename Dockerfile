FROM golang:1.24.2
WORKDIR /usr/src/app
COPY go.mod go.sum ./
COPY *.go ./
COPY cmd/ ./cmd
COPY internal/ ./internal
COPY resources/ ./resources
RUN go mod download
# RUN GOOS=linux GOARCH=amd64 go build -o /confidant_server ./cmd/confidant_server/main.go
# CMD ["/confidant_server"]
