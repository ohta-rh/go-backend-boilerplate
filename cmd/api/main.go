package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/tetsuyaohta/go-backend-boilerplate/internal/handler"
	"github.com/tetsuyaohta/go-backend-boilerplate/internal/infrastructure/database"
	"github.com/tetsuyaohta/go-backend-boilerplate/internal/repository"
	"github.com/tetsuyaohta/go-backend-boilerplate/internal/usecase"
)

// run is the actual application logic, separate from main to properly handle defer.
func run() error {
	// 環境変数の取得
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development" // デフォルトは開発環境
	}

	// 環境に応じた設定
	if env == "development" {
		// 開発環境の場合はデバッグモードを有効化
		log.Println("Running in development mode")
		gin.SetMode(gin.DebugMode)
	} else {
		// 本番環境の場合はリリースモードを設定
		log.Println("Running in production mode")
		gin.SetMode(gin.ReleaseMode)
	}

	// PostgreSQLの接続情報（Docker Composeの設定に基づく）
	dsn := "postgresql://root:password@db:5432/app?sslmode=disable"

	// Database client initialization
	client, err := database.NewClient(dsn, 5)
	if err != nil {
		return err
	}
	defer client.Close()

	log.Println("Successfully connected to database")

	// リポジトリの作成
	userRepo := repository.NewUserRepository(client)

	// ユースケースの作成
	userInteractor := usecase.NewUserInteractor(userRepo)

	// ルーターのセットアップ
	r := handler.SetupRouter(userInteractor)

	// サーバーの起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // デフォルトポート
	}

	log.Printf("Server starting on port %s...\n", port)

	// Run the server and handle the error
	return r.Run(":" + port)
}

func main() {
	if err := run(); err != nil {
		log.Printf("Application error: %v", err)
		os.Exit(1)
	}
}
