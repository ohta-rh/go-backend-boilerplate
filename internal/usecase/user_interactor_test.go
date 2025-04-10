package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	domainUser "easy-go-backend/internal/domain/user"
)

// MockUserRepository はテスト用のモックリポジトリです.
type MockUserRepository struct {
	users map[int]*domainUser.User
}

// NewMockUserRepository はモックリポジトリを作成します.
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[int]*domainUser.User),
	}
}

// GetByID はIDによるユーザー取得をモックします.
func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*domainUser.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}

	// nilとnilを返す代わりに、センチネルエラーを使用
	return nil, domainUser.ErrUserNotFound
}

// Create はユーザー作成をモックします.
func (m *MockUserRepository) Create(ctx context.Context, u *domainUser.User) (*domainUser.User, error) {
	m.users[u.ID] = u
	return u, nil
}

// Update はユーザー更新をモックします.
func (m *MockUserRepository) Update(ctx context.Context, u *domainUser.User) (*domainUser.User, error) {
	if _, ok := m.users[u.ID]; ok {
		m.users[u.ID] = u
		return u, nil
	}

	// nilとnilを返す代わりに、センチネルエラーを使用
	return nil, domainUser.ErrUserNotFound
}

// Delete はユーザー削除をモックします.
func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
	delete(m.users, id)
	return nil
}

// GetAll は全ユーザー取得をモックします.
func (m *MockUserRepository) GetAll(ctx context.Context) ([]*domainUser.User, error) {
	users := make([]*domainUser.User, 0, len(m.users))
	for _, u := range m.users {
		users = append(users, u)
	}

	return users, nil
}

func TestUserInteractor_GetUser(t *testing.T) {
	// モックリポジトリの準備
	repo := NewMockUserRepository()
	// テスト用ユーザーデータの作成
	testUser := &domainUser.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
	}
	repo.users[testUser.ID] = testUser

	// テスト対象のInteractorを作成
	interactor := NewUserInteractor(repo)

	// テストケース
	t.Run("existing user", func(t *testing.T) {
		ctx := context.Background()
		result, err := interactor.GetUser(ctx, 1)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if result == nil {
			t.Fatal("Expected user, got nil")
		}

		if result.ID != testUser.ID {
			t.Errorf("Expected ID %d, got %d", testUser.ID, result.ID)
		}

		if result.Name != testUser.Name {
			t.Errorf("Expected Name %s, got %s", testUser.Name, result.Name)
		}

		if result.Email != testUser.Email {
			t.Errorf("Expected Email %s, got %s", testUser.Email, result.Email)
		}
	})

	t.Run("non-existing user", func(t *testing.T) {
		ctx := context.Background()
		result, err := interactor.GetUser(ctx, 999)

		if !errors.Is(err, domainUser.ErrUserNotFound) {
			t.Fatalf("Expected error %v, got %v", domainUser.ErrUserNotFound, err)
		}

		if result != nil {
			t.Errorf("Expected nil, got %+v", result)
		}
	})
}

func TestUserInteractor_CreateUser(t *testing.T) {
	// モックリポジトリの準備
	repo := NewMockUserRepository()

	// テスト対象のInteractorを作成
	interactor := NewUserInteractor(repo)

	// テストケース
	t.Run("create new user", func(t *testing.T) {
		ctx := context.Background()
		newUser := &domainUser.User{
			ID:        2,
			Name:      "New User",
			Email:     "new@example.com",
			CreatedAt: time.Now(),
		}

		result, err := interactor.CreateUser(ctx, newUser)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if result == nil {
			t.Fatal("Expected user, got nil")
		}

		if result.ID != newUser.ID {
			t.Errorf("Expected ID %d, got %d", newUser.ID, result.ID)
		}

		// リポジトリに保存されていることを確認
		savedUser, _ := repo.GetByID(ctx, newUser.ID)
		if savedUser == nil {
			t.Fatal("User not saved in repository")
		}
	})
}

func TestUserInteractor_GetAllUsers(t *testing.T) {
	// モックリポジトリの準備
	repo := NewMockUserRepository()

	// テスト用ユーザーデータの作成
	user1 := &domainUser.User{ID: 1, Name: "User 1", Email: "user1@example.com", CreatedAt: time.Now()}
	user2 := &domainUser.User{ID: 2, Name: "User 2", Email: "user2@example.com", CreatedAt: time.Now()}

	repo.users[user1.ID] = user1
	repo.users[user2.ID] = user2

	// テスト対象のInteractorを作成
	interactor := NewUserInteractor(repo)

	// テストケース
	t.Run("get all users", func(t *testing.T) {
		ctx := context.Background()
		results, err := interactor.GetAllUsers(ctx)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(results) != 2 {
			t.Errorf("Expected 2 users, got %d", len(results))
		}

		// 特定のユーザーが含まれていることを確認（順序は保証されないため、IDのみチェック）
		foundIDs := make(map[int]bool)
		for _, u := range results {
			foundIDs[u.ID] = true
		}

		if !foundIDs[1] || !foundIDs[2] {
			t.Errorf("Expected users with IDs 1 and 2, got %v", foundIDs)
		}
	})
}
