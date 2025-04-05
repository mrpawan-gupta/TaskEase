FROM golang:1.24-alpine AS builder
LABEL maintainer="pawanroshankumargupta"

WORKDIR /app
RUN apk --no-cache add git ca-certificates tzdata gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
RUN go get github.com/lib/pq
COPY . .
RUN go build -o taskease ./main.go

# Runtime stage
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/taskease .
COPY --from=builder /app/config ./config

ENV GIN_MODE=release
ENV TZ=UTC
EXPOSE 8080
CMD ["/app/taskease"]