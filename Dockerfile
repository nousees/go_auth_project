FROM golang:1.23-alpine

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]

# FROM golang:1.23-alpine AS builder

# WORKDIR /app

# COPY go.* ./

# RUN go mod download

# COPY . .

# RUN go build -o main ./cmd/main.go


# FROM alpine:latest

# WORKDIR /app

# COPY --from=builder /app/main .

# EXPOSE 8080

# CMD ["./main"]