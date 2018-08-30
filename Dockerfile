# Build the dashboard binary
FROM golang:1.10.3 as builder

# Copy in the go src
WORKDIR /go/src/github.com/presslabs/dashboard
COPY pkg/    pkg/
COPY cmd/    cmd/
COPY vendor/ vendor/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o dashboard github.com/presslabs/dashboard/cmd/dashboard

# Copy the dashboard binary into a thin image
FROM scratch
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/github.com/presslabs/dashboard/dashboard /
ENTRYPOINT ["/dashboard"]
CMD ["help"]
