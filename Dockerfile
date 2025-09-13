FROM golang:1.23-alpine

WORKDIR /app

# Install git, curl, dan dependencies lain yang diperlukan
RUN apk add --no-cache git curl bash

# Download dan install Air binary secara langsung
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b /usr/local/bin

# Copy go.mod dan go.sum untuk cache dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy semua source code
COPY . .

# Expose port aplikasi
EXPOSE 8080

# Gunakan air untuk hot reload
CMD ["air"]