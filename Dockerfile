# 開発用のコンテナ
FROM golang:1.24.2

WORKDIR /app

# ソースコードをコンテナにコピー
COPY . .

# Install Air and golangci-lint
RUN go install github.com/air-verse/air@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

