FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o app .

FROM gcr.io/distroless/static-debian13

WORKDIR /app

COPY --from=builder /build/app .

EXPOSE 8000

EntryPoint ["./app"]