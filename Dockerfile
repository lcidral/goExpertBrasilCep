FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

# Build com CGO desabilitado
RUN CGO_ENABLED=0 GOOS=linux go build -o cepfinder

# Imagem final
FROM alpine:latest

# Instala certificados
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/cepfinder .

ENTRYPOINT ["./cepfinder"]