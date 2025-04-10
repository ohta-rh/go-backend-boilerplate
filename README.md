# Go Backend Boilerplate

Gin+Airを使用したGoバックエンドのボイラープレートコードです。

## 機能

- [Gin](https://github.com/gin-gonic/gin) Webフレームワーク
- [Air](https://github.com/cosmtrek/air) ホットリロード
- Dockerベースの開発環境
- レイヤー化されたアーキテクチャ

## フォルダ構成

```
yourapp/
├── cmd/
│   └── api/
│       └── main.go        ← Ginの起動ポイント
├── internal/
│   └── handler/
│       └── router.go      ← ルーティング定義
├── go.mod
├── go.sum
└── .air.toml              ← 開発用ホットリロード設定
```

## 開発の始め方

### 必要条件

- Docker と Docker Compose
- Go 1.20以上（ローカル開発の場合）

### Dockerを使った開発

1. コンテナを起動:

```bash
docker-compose up -d
```

2. アプリケーションにアクセス:

```
http://localhost:8080
```

3. APIの例:

```
http://localhost:8080/api/v1/hello
http://localhost:8080/health
```

### Makeコマンド

Makefileには便利なコマンドが用意されています:

```bash
# パッケージの追加
make add-pkg

# 依存関係のインストール
make install

# テストの実行
make test

# カバレッジ付きテスト
make test.cover

# 特定パッケージのテスト
make test.pkg

# リンター実行
make lint

# コンテナ内のシェルに入る
make api.bash
```

## 環境変数

- `GO_ENV`: 実行環境 (`development` または `production`)
- `PORT`: サーバーポート (デフォルト: `8080`) 