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

# Install godotenv
RUN go get github.com/joho/godotenv

# Copy the .env file to the container (if needed for build or runtime)
COPY .env .env

# Build the Go application
RUN go build -o main ./cmd/gamenet

# Command to run the Go application
CMD ["./main"]
