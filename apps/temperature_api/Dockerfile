FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

# Инициализируем модуль, если его нет (для корректной сборки в новых версиях Go)
RUN [ ! -f go.mod ] && go mod init temperature-api || true

RUN go build -o temperature-api main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/temperature-api .

EXPOSE 8081

CMD ["./temperature-api"]
