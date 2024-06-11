# Project will be using GoLang 1.22
FROM golang:1.22 as builder

# Working directory inside of Docker Container
WORKDIR /app

# Copy mod and sum files to Docker Container
COPY go.mod go.sum ./

# Downloads Go dependencies if not already existing
RUN go mod download

# Copy source from current dir to working dir
COPY . .

# Build Go app
RUN go build -o gamenet cmd/gamenet/main.go

# New stage created
FROM alpine:latest

WORKDIR /root/

# Prebuilt binary file from previous stage
COPY --from=builder /app/gamenet .

# Now the final command to run everythering
CMD ["./gamenet"]
