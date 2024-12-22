FROM golang:1.22.3-alpine

# Copy the rest of the application code
COPY . .
# Install bash and SQLite3 development libraries using apk
RUN apk add --no-cache bash sqlite sqlite-dev gcc musl-dev
RUN go build -o main ./cmd/api/main.go

# set the metadata for the image
LABEL version="0.0.1"
LABEL description="This is a Dockerfile for the forum project"
# Run the compiled Go application
CMD ["./main"]