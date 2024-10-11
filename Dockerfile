# Use the official Golang image as the build stage
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go application
RUN go build -o gamenet ./cmd/gamenet

# Use the official minimal base image for Go
FROM golang:1.20

WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/gamenet .

# Expose the necessary port (8080 in this example)
EXPOSE 8080

# Run the Go application
CMD ["./gamenet"]
