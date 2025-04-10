package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler はヘルスチェック関連のエンドポイントを処理します.
type HealthHandler struct{}

// NewHealthHandler は新しいHealthHandlerを作成します.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// RegisterRoutes はヘルスチェック関連のルートを登録します.
func (h *HealthHandler) RegisterRoutes(r *gin.Engine) {
	// ヘルスチェック用エンドポイント
	r.GET("/health", h.HealthCheck)
}

// HealthCheck はシステムの健全性を確認するエンドポイントです.
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
