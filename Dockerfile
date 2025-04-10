# 開発用のコンテナ
FROM golang:1.24.2

WORKDIR /app

# ソースコードをコンテナにコピー
COPY . .

# Goパッケージのインストール
RUN go mod download

# エントリポイント設定（開発環境用）
CMD ["air", "-c", ".air.toml"]
