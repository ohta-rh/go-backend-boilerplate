package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tetsuyaohta/go-backend-boilerplate/internal/usecase"
)

// SetupRouter はGinのルーターを設定し、すべてのルートを定義します.
func SetupRouter(userInteractor *usecase.UserInteractor) *gin.Engine {
	// デフォルトのミドルウェア（Logger, Recovery）を使用
	r := gin.Default()

	// ヘルスチェック用エンドポイント
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// APIグループ
	v1 := r.Group("/api/v1")
	{
		// Hello Worldエンドポイント
		v1.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello, World!",
			})
		})

		// User routes
		userHandler := NewUserHandler(userInteractor)
		users := v1.Group("/users")
		users.GET("", userHandler.GetAllUsers)
		users.POST("", userHandler.CreateUser)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}

	return r
}
