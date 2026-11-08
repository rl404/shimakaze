# Golang base image
FROM golang:1.27.3 AS go_builder
LABEL stage="shimakaze_builder"
WORKDIR /shimakaze
COPY . .
WORKDIR /shimakaze/cmd/shimakaze
RUN CGO_ENABLED=0 go build -mod=vendor -ldflags="-s -w" -o shimakaze

# New stage from scratch
FROM alpine:3.24.3
LABEL stage="shimakaze_binary"
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 shimakaze
COPY --from=go_builder /shimakaze/cmd/shimakaze/shimakaze /usr/local/bin/shimakaze
WORKDIR /home/shimakaze
USER shimakaze
ENTRYPOINT ["/usr/local/bin/shimakaze"]