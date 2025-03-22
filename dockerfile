# Use the official Golang image for building the application
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the Go binary for Linux (CGO_ENABLED=1 enables C bindings for SQLite)
RUN CGO_ENABLED=1 GOOS=linux go build -o /gdragon ./cmd

# Expose the application port (adjust this based on your app's port)
EXPOSE 8080

# Set the entrypoint to run the compiled Go binary when the container starts
CMD ["/gdragon"]
