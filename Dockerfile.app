# Production

# FROM golang:1.21.1

# WORKDIR /app
# COPY go.mod go.sum ./

# RUN go mod download

# COPY . .

# RUN go build -o main ./cmd/main.go

# EXPOSE 8080

# CMD ["./main"]


# dev
FROM golang:1.21.1
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8080
EXPOSE 8000

CMD ["go", "run", "cmd/main.go"]
