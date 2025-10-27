# --- Этап 1: Сборщик (Builder) ---
    FROM golang:1.24-alpine AS builder
    
    WORKDIR /app
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Копируем все, включая папку config
    COPY . .
    
    RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main ./cmd/backend
    
    # --- Этап 2: Финальный образ ---
    FROM alpine:latest
    
    RUN apk --no-cache add ca-certificates
    
    WORKDIR /app
    
    # Сначала копируем бинарник
    COPY --from=builder /app/main .
    
    # --- ДОБАВЛЕННАЯ СТРОКА ---
    # Копируем папку config из сборщика в финальный образ
    COPY --from=builder /app/config ./config
    
    EXPOSE 8006
    
    CMD ["/app/main"]