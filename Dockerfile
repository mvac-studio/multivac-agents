FROM golang:alpine AS builder

#Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR ./
COPY . .

# Fetch dependencies.
# Using go get.
RUN go get -d -v

# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/services

# Use scratch image for minimal size.
FROM scratch
ENV PORT=8080
EXPOSE 8080
COPY --from=builder /go/bin/service /go/bin/service
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/go/bin/service"]
