package handler

import (
	"easy-go-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

// SetupRouter はGinのルーターを設定し、すべてのルートを定義します.
func SetupRouter(userInteractor *usecase.UserInteractor) *gin.Engine {
	// デフォルトのミドルウェア（Logger, Recovery）を使用
	r := gin.Default()

	// ヘルスチェックハンドラーの登録
	healthHandler := NewHealthHandler()
	healthHandler.RegisterRoutes(r)

	// APIグループ
	v1 := r.Group("/api/v1")
	{
		// Hello Worldハンドラーの登録
		helloHandler := NewHelloHandler()
		helloHandler.RegisterRoutes(v1)

		// ユーザーハンドラーの登録
		userHandler := NewUserHandler(userInteractor)
		userHandler.RegisterRoutes(v1)
	}

	return r
}
