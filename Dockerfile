FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o server .

FROM gcr.io/distroless/static-debian13

WORKDIR /app

COPY --from=builder /build/server /app/server
COPY resource /app/resource

EXPOSE 8000

ENTRYPOINT ["/app/server"]