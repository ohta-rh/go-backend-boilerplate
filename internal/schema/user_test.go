package schema

import (
	"testing"
	"time"

	"easy-go-backend/internal/domain/user"
)

func TestUserRequest_ToEntity(t *testing.T) {
	// テストデータ
	req := &UserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}

	// 関数を実行
	entity := req.ToEntity()

	// 結果を検証
	if entity.Name != req.Name {
		t.Errorf("Expected Name %s, got %s", req.Name, entity.Name)
	}

	if entity.Email != req.Email {
		t.Errorf("Expected Email %s, got %s", req.Email, entity.Email)
	}
}

func TestFromEntity(t *testing.T) {
	// テストデータ
	now := time.Now()
	entity := &user.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: now,
	}

	// 関数を実行
	response := FromEntity(entity)

	// 結果を検証
	if response.ID != entity.ID {
		t.Errorf("Expected ID %d, got %d", entity.ID, response.ID)
	}

	if response.Name != entity.Name {
		t.Errorf("Expected Name %s, got %s", entity.Name, response.Name)
	}

	if response.Email != entity.Email {
		t.Errorf("Expected Email %s, got %s", entity.Email, response.Email)
	}

	if !response.CreatedAt.Equal(entity.CreatedAt) {
		t.Errorf("Expected CreatedAt %v, got %v", entity.CreatedAt, response.CreatedAt)
	}
}

func TestFromEntities(t *testing.T) {
	// テストデータ
	now := time.Now()
	entities := []*user.User{
		{
			ID:        1,
			Name:      "User 1",
			Email:     "user1@example.com",
			CreatedAt: now,
		},
		{
			ID:        2,
			Name:      "User 2",
			Email:     "user2@example.com",
			CreatedAt: now,
		},
	}

	// 関数を実行
	responses := FromEntities(entities)

	// 結果を検証
	if len(responses) != len(entities) {
		t.Errorf("Expected %d responses, got %d", len(entities), len(responses))
	}

	for i, entity := range entities {
		response := responses[i]

		if response.ID != entity.ID {
			t.Errorf("Entity %d: Expected ID %d, got %d", i, entity.ID, response.ID)
		}

		if response.Name != entity.Name {
			t.Errorf("Entity %d: Expected Name %s, got %s", i, entity.Name, response.Name)
		}

		if response.Email != entity.Email {
			t.Errorf("Entity %d: Expected Email %s, got %s", i, entity.Email, response.Email)
		}

		if !response.CreatedAt.Equal(entity.CreatedAt) {
			t.Errorf("Entity %d: Expected CreatedAt %v, got %v", i, entity.CreatedAt, response.CreatedAt)
		}
	}
}
