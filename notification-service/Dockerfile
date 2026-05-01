# notification-service/Dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем всё содержимое папки notification-service в Docker
# Это избавит нас от проблем с путями
COPY . .

# Собираем проект
# -mod=vendor указывает брать зависимости локально из папки vendor
RUN go build -mod=vendor -o main .

FROM alpine:latest
WORKDIR /root/
# Копируем бинарник из стадии builder
COPY --from=builder /app/main .

# Запускаем
CMD ["./main"]