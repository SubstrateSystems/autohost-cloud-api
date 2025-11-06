# Etapa 1: build
FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go build -o autohost-cloud-lite ./cmd/api

# Etapa 2: runtime
FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/autohost-cloud-lite /usr/local/bin/autohost-cloud-lite
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/autohost-cloud-lite"]
