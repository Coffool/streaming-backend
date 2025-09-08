
# ======================================
# Etapa 1: Build
# ======================================
FROM golang:1.24.5 AS builder

# Definir directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar archivos de dependencias primero para aprovechar cache
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Compilar el binario para Linux
ARG SERVICE_NAME=history-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${SERVICE_NAME} .

# ======================================
# Etapa 2: Runtime
# ======================================
FROM alpine:3.19

# Instalar certificados SSL (para HTTPS)
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copiar solo el binario compilado desde la etapa anterior
ARG SERVICE_NAME=history-service
COPY --from=builder /app/${SERVICE_NAME} .

# Exponer el puerto que use el microservicio
EXPOSE 8081

# Comando de inicio
CMD ["./history-service"]
