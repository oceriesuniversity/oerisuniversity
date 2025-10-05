# Use an official Go image
FROM golang:1.23.0

# Install GCC and SQLite dev libraries
RUN apt-get update && apt-get install -y gcc libsqlite3-dev

# Set environment variables
ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory
WORKDIR /app

# Copy the application code
COPY . .

# Build the application
RUN go build -o main .

# Run the application
CMD ["./main"]
