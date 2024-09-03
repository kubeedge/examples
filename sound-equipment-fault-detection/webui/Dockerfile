# Build using the official Go image
FROM golang:1.22-alpine AS build

# Set the working directory
WORKDIR /app

# Copy the Go source code to the container
COPY . .

# Compile the Go application
RUN go build -o main ./main.go

# Use a smaller base image that only contains the runtime
FROM alpine:latest

# Copy the compiled binary to the final image
COPY --from=build /app/main /app/main

# Copy the static directory to the final image
COPY --from=build /app/static /app/static

# Copy the config to the final image
# COPY --from=build /app/config /app/config

# Set the working directory
WORKDIR /app

# Set the command to run the Go application
CMD ["./main"]