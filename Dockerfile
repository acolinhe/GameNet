# Use the official Golang image
FROM golang:1.22.4

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download the Go dependencies
RUN go mod download

# Copy the source code to the container
COPY . .

# Command to run the Go application
CMD ["go", "run", "cmd/gamenet/main.go"]
