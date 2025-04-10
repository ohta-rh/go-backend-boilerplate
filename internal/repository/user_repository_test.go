package repository

import (
	"testing"
	"time"

	"easy-go-backend/ent"
)

// モックEntUserを作成するヘルパー関数.
func createMockEntUser(id int, name, email string) *ent.User {
	return &ent.User{
		ID:        id,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}
}

// mapEntUserToDomainUserのテスト.
func TestMapEntUserToDomainUser(t *testing.T) {
	// テストデータ
	mockEntUser := createMockEntUser(1, "Test User", "test@example.com")

	// 関数を実行
	var domainUser = mapEntUserToDomainUser(mockEntUser)

	// 結果を検証
	if domainUser.ID != mockEntUser.ID {
		t.Errorf("Expected ID %d, got %d", mockEntUser.ID, domainUser.ID)
	}

	if domainUser.Name != mockEntUser.Name {
		t.Errorf("Expected Name %s, got %s", mockEntUser.Name, domainUser.Name)
	}

	if domainUser.Email != mockEntUser.Email {
		t.Errorf("Expected Email %s, got %s", mockEntUser.Email, domainUser.Email)
	}

	if !domainUser.CreatedAt.Equal(mockEntUser.CreatedAt) {
		t.Errorf("Expected CreatedAt %v, got %v", mockEntUser.CreatedAt, domainUser.CreatedAt)
	}
}

// NewUserRepositoryのテスト.
func TestNewUserRepository(t *testing.T) {
	// モッククライアント
	mockClient := &ent.Client{}

	// リポジトリを作成
	repo := NewUserRepository(mockClient)

	// 結果を検証
	if repo == nil {
		t.Error("Expected repository to be created, got nil")
		return
	}

	if repo.client != mockClient {
		t.Error("Repository client does not match the provided client")
	}
}

// 注意: 実際のアプリケーションでは、より詳細なテストが必要です
// 例えば、GetByID、Create、Update、Delete、GetAllメソッドのテスト
// これらのテストでは、entのクライアントをモックする必要があります
