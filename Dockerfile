FROM golang:1.19-alpine AS builder

WORKDIR /app

# Bağımlılıkları kopyala ve indir
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodları kopyala
COPY . .

# Uygulamayı derle
RUN go build -o main

# Çalıştırma aşaması
FROM alpine:3.14

WORKDIR /app

# Derlenen uygulamayı kopyala
COPY --from=builder /app/main .
# Gerekli dosyaları kopyala
COPY --from=builder /app/database/schema.sql ./database/

# Uygulamayı çalıştır
CMD ["./main"] 