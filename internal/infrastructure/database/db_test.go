package database

import (
	"testing"
)

func TestNewClient_InvalidDSN(t *testing.T) {
	// 無効なDSNでテスト
	_, err := NewClient("invalid-dsn", 1)
	if err == nil {
		t.Error("Expected error with invalid DSN, but got nil")
	}
}

// 注意: このテストは実際のデータベース接続を試みないようにしています
// 実際のアプリケーションでは、テスト用のデータベースを使用するか
// モックを使用してテストすることをお勧めします
