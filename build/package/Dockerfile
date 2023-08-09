# Golang base image
FROM golang:1.21 as go_builder
LABEL stage="shimakaze_builder"
WORKDIR /shimakaze
COPY . .
WORKDIR /shimakaze/cmd/shimakaze
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o shimakaze

# New stage from scratch
FROM alpine:3.18
LABEL stage="shimakaze_binary"
RUN apk add --no-cache ca-certificates && apk add --no-cache tzdata
COPY --from=go_builder /shimakaze/cmd/shimakaze/shimakaze /cmd/shimakaze/shimakaze
WORKDIR /cmd/shimakaze
