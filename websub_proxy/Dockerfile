# Use the official Go image as the base
FROM golang:1.23

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the Fiber app's port (3000)
EXPOSE 3000

# Command to run the application
CMD ["./main"]
