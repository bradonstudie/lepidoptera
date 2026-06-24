FROM golang:1.25-alpine AS dev
RUN apk add --no-cache imagemagick
RUN go install github.com/air-verse/air@latest
RUN go install github.com/a-h/templ/cmd/templ@latest
WORKDIR /app
CMD ["air"]

FROM golang:1.25-alpine AS builder
RUN go install github.com/a-h/templ/cmd/templ@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN templ generate && CGO_ENABLED=0 go build -o /lepidoptera ./cmd/server

FROM alpine:3.20 AS runtime
RUN apk add --no-cache imagemagick
WORKDIR /app
COPY --from=builder /lepidoptera .
CMD ["./lepidoptera"]
