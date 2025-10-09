# Stage 1: builder
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git, curl, bash
RUN apk add --no-cache git curl bash

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary (output: /app/app)
RUN go build -o app

# Stage 2: dev (pakai air)
FROM golang:1.23-alpine AS dev

WORKDIR /app
RUN apk add --no-cache git curl bash

# Download Air
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b /usr/local/bin

COPY go.mod go.sum ./
RUN go mod download
COPY . .

EXPOSE 8080

CMD ["air"]

# Stage 3: production
FROM alpine:latest AS prod

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]

# Final stage: conditional
# Default to "dev", override with --target prod for production