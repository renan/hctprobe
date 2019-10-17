FROM golang:1.13-alpine AS builder
WORKDIR /usr/src/app
COPY . .
RUN go build -o hctprobe .

FROM alpine:3.10
RUN apk --no-cache add ca-certificates
USER nobody
COPY --from=builder /usr/src/app/hctprobe /hctprobe
ENTRYPOINT ["/hctprobe"]
