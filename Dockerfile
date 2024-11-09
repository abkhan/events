# Use the official Golang image as the base image
FROM golang:1.23.3-alpine3.20 AS local_builder
RUN apk update && apk add --no-cache bash git   # update local

ENV GOPATH /go
ENV GOCACHE /go/.cache
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Set the working directory inside the container
WORKDIR /app

# Copy the project's Go module files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/admin ./cmd/admin

# Expose any necessary ports
EXPOSE 8080

# Set the command to run when the container starts
CMD ["/bin/admin"]

# docker run -p 8080:8080 -it admin