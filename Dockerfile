FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git --repository https://dl-cdn.alpinelinux.org/alpine/v3.19/main || true

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o image-palette-service \
    ./cmd/server

FROM scratch

WORKDIR /app

COPY --from=builder /app/image-palette-service .

EXPOSE 8080

CMD ["./image-palette-service"]