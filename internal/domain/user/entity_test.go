package user

import (
	"testing"
	"time"
)

func TestUserEntity(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name  string
		user  User
		check func(*testing.T, User)
	}{
		{
			name: "valid user",
			user: User{
				ID:        1,
				Name:      "Test User",
				Email:     "test@example.com",
				CreatedAt: now,
			},
			check: func(t *testing.T, u User) {
				if u.ID != 1 {
					t.Errorf("Expected ID to be 1, got %d", u.ID)
				}
				if u.Name != "Test User" {
					t.Errorf("Expected Name to be 'Test User', got %s", u.Name)
				}
				if u.Email != "test@example.com" {
					t.Errorf("Expected Email to be 'test@example.com', got %s", u.Email)
				}
				if !u.CreatedAt.Equal(now) {
					t.Errorf("Expected CreatedAt to be %v, got %v", now, u.CreatedAt)
				}
			},
		},
		{
			name: "zero value user",
			user: User{},
			check: func(t *testing.T, u User) {
				if u.ID != 0 {
					t.Errorf("Expected ID to be 0, got %d", u.ID)
				}
				if u.Name != "" {
					t.Errorf("Expected Name to be empty, got %s", u.Name)
				}
				if u.Email != "" {
					t.Errorf("Expected Email to be empty, got %s", u.Email)
				}
				if !u.CreatedAt.IsZero() {
					t.Errorf("Expected CreatedAt to be zero time, got %v", u.CreatedAt)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, tt.user)
		})
	}
}

func TestUserFieldValidation(t *testing.T) {
	// 現在はフィールドの存在確認のみ実装
	t.Run("user structure has expected fields", func(t *testing.T) {
		user := User{
			ID:        123,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
		}

		// フィールドの存在確認
		if user.ID != 123 {
			t.Error("ID field not properly set/accessed")
		}

		if user.Name != "John Doe" {
			t.Error("Name field not properly set/accessed")
		}

		if user.Email != "john@example.com" {
			t.Error("Email field not properly set/accessed")
		}

		if user.CreatedAt.IsZero() {
			t.Error("CreatedAt field not properly set/accessed")
		}
	})
}
