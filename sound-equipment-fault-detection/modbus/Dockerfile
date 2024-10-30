# Build the application using the official Go 1.22 image
FROM golang:1.22-alpine AS build

# Set the working directory
WORKDIR /app

# Copy the source code to the container
COPY . .

# Compile the Go application
RUN go build -o main ./main.go

# Run the application using a smaller base image
FROM alpine:3.18

# Set the working directory
WORKDIR /app

# Copy the binary generated in the build phase to the final image
COPY --from=build /app/main /app/main

# # Set the entry point of the application
# ENTRYPOINT ["/app/main"]

# If there are default parameters, you can set them here (optional)
# CMD ["--help"]
