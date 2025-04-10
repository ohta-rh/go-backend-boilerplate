package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HelloHandler はHello World関連のエンドポイントを処理します
type HelloHandler struct{}

// NewHelloHandler は新しいHelloHandlerを作成します
func NewHelloHandler() *HelloHandler {
	return &HelloHandler{}
}

// RegisterRoutes はHello World関連のルートを登録します
func (h *HelloHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Hello Worldエンドポイント
	router.GET("/hello", h.HelloWorld)
}

// HelloWorld は基本的なHello Worldレスポンスを返すエンドポイントです
func (h *HelloHandler) HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}