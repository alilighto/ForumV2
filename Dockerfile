FROM golang:latest
WORKDIR /usr/src/app
COPY . .
RUN go build -o main ./cmd/api/main.go
# Run the compiled Go application
CMD ["./main"]