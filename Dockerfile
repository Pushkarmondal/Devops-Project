# First stage: Build the Go binary
FROM golang:1.22.3-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code to the container
COPY . .

# Build the Go binary
RUN go build -o myapp .

# Second stage: Run the binary
FROM alpine:3.18

# Set the working directory inside the second stage
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/myapp .

COPY index.html /app/

# Expose port if necessary (e.g., for an API)
EXPOSE 8080

# Command to run the binary
CMD ["./myapp"]
