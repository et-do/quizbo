# Use the official Golang image as a base
FROM golang:1.22.4

# Set environment variables
ENV PORT=8080

# Set the current working directory inside the container
WORKDIR /workspace

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if go.mod and go.sum are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Expose port 8080 for the app
EXPOSE 8080

# Build the Go application
RUN go build -o main .

# Set the default command for the container
CMD ["./main"]
