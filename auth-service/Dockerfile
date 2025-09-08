# Etapa 1: Build
FROM golang:1.24.5 AS builder

# Definir directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar go.mod y go.sum primero (para cachear dependencias)
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Compilar el binario para Linux (GOOS y GOARCH garantizan compatibilidad)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth-service .

# Etapa 2: Runtime
FROM alpine:3.19

# Instalar certificados SSL (importante si haces requests HTTPS)
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copiar solo el binario compilado desde la etapa anterior
COPY --from=builder /app/auth-service .

# Exponer el puerto del servicio
EXPOSE 8080

# Comando de inicio
CMD ["./auth-service"]

