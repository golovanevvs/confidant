FROM golang:1.24.2
#ENV ADDRESS=":3000"
#ENV DATABASE_DSN="postgres://postgres:password@db:5432/confidant?sslmode=disable"
#"host=db port=5432 user=postgres password=password dbname=confidant sslmode=disable"
#ENV POSTGRES_USER="postgres"
#ENV POSTGRES_PASSWORD="password"
#ENV POSTGRES_DB="confidant"
WORKDIR /usr/src/app
COPY go.mod go.sum ./
COPY *.go ./
COPY cmd/ ./cmd
COPY internal/ ./internal
COPY resources/ ./resources
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -o /confidant_server ./cmd/confidant_server/main.go
CMD ["/confidant_server"]
