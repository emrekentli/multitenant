# Use multi-stage builds to keep the final image small
FROM golang:alpine AS builder

# Install necessary packages
RUN apk update && apk add --no-cache git

# Set the working directory
WORKDIR /go/src/app/

# Copy the source code
COPY . .

# Tidy up the Go module dependencies
RUN go mod tidy

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app ./cmd/main.go

# Use a minimal base image for the final stage
FROM scratch

# Copy the built binary from the builder stage
COPY --from=builder /go/bin/app /go/bin/app
COPY --from=builder /go/src/app/views /views
COPY --from=builder /go/src/app/migrations /migrations
COPY --from=builder /go/src/app/public /public


# Expose the application port
EXPOSE 8080

# Set the entrypoint to the built binary
ENTRYPOINT ["/go/bin/app"]