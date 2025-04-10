# 開発用のコンテナ
FROM golang:1.24.2

WORKDIR /app

# ソースコードをコンテナにコピー
COPY . .

# Install Air and golangci-lint and modernize
RUN go install github.com/air-verse/air@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
RUN go install golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest

