package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	domainUser "easy-go-backend/internal/domain/user"

	"github.com/gin-gonic/gin"
)

// MockUserInteractor はテスト用のユースケースモックです.
type MockUserInteractor struct {
	users map[int]*domainUser.User
	err   error
}

// NewMockUserInteractor はモックインタラクターを作成します.
func NewMockUserInteractor() *MockUserInteractor {
	return &MockUserInteractor{
		users: make(map[int]*domainUser.User),
	}
}

// SetError はテスト用にエラーを設定します.
func (m *MockUserInteractor) SetError(err error) {
	m.err = err
}

// GetUser はユーザー取得をモックします.
func (m *MockUserInteractor) GetUser(ctx context.Context, id int) (*domainUser.User, error) {
	if m.err != nil {
		return nil, m.err
	}

	if u, ok := m.users[id]; ok {
		return u, nil
	}

	// 存在しないユーザーの場合は、ErrUserNotFoundを返す
	return nil, domainUser.ErrUserNotFound
}

// CreateUser はユーザー作成をモックします.
func (m *MockUserInteractor) CreateUser(ctx context.Context, u *domainUser.User) (*domainUser.User, error) {
	if m.err != nil {
		return nil, m.err
	}

	if u.ID == 0 {
		// 実際のアプリケーションでは自動ID採番されるロジックをモック
		u.ID = len(m.users) + 1
	}

	m.users[u.ID] = u

	return u, nil
}

// UpdateUser はユーザー更新をモックします.
func (m *MockUserInteractor) UpdateUser(ctx context.Context, u *domainUser.User) (*domainUser.User, error) {
	if m.err != nil {
		return nil, m.err
	}

	if _, ok := m.users[u.ID]; ok {
		m.users[u.ID] = u
		return u, nil
	}

	return nil, domainUser.ErrUserNotFound
}

// DeleteUser はユーザー削除をモックします.
func (m *MockUserInteractor) DeleteUser(ctx context.Context, id int) error {
	if m.err != nil {
		return m.err
	}

	if _, ok := m.users[id]; !ok {
		return domainUser.ErrUserNotFound
	}

	delete(m.users, id)

	return nil
}

// GetAllUsers は全ユーザー取得をモックします.
func (m *MockUserInteractor) GetAllUsers(ctx context.Context) ([]*domainUser.User, error) {
	if m.err != nil {
		return nil, m.err
	}

	users := make([]*domainUser.User, 0, len(m.users))
	for _, u := range m.users {
		users = append(users, u)
	}

	return users, nil
}

// UserHandlerHelper はハンドラーのヘルパーメソッドを定義します.
type UserHandlerHelper interface {
	handleGetUser(c *gin.Context, id int) (*domainUser.User, error)
	handleCreateUser(c *gin.Context, userDTO *domainUser.User) (*domainUser.User, error)
	handleUpdateUser(c *gin.Context, id int, userDTO *domainUser.User) (*domainUser.User, error)
	handleDeleteUser(c *gin.Context, id int) error
	handleGetAllUsers(c *gin.Context) ([]*domainUser.User, error)
}

// handleGetUserRequest handles GET /users/:id requests.
func handleGetUserRequest(c *gin.Context, handler UserHandlerHelper) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	u, err := handler.handleGetUser(c, id)
	if err != nil {
		if errors.Is(err, domainUser.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, u)
}

// handleCreateUserRequest handles POST /users requests.
func handleCreateUserRequest(c *gin.Context, handler UserHandlerHelper) {
	var userDTO domainUser.User
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := handler.handleCreateUser(c, &userDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// handleUpdateUserRequest handles PUT /users/:id requests.
func handleUpdateUserRequest(c *gin.Context, handler UserHandlerHelper) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var userDTO domainUser.User
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userDTO.ID = id

	updatedUser, err := handler.handleUpdateUser(c, id, &userDTO)
	if err != nil {
		if errors.Is(err, domainUser.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// handleDeleteUserRequest handles DELETE /users/:id requests.
func handleDeleteUserRequest(c *gin.Context, handler UserHandlerHelper) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := handler.handleDeleteUser(c, id); err != nil {
		if errors.Is(err, domainUser.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// handleGetAllUsersRequest handles GET /users requests.
func handleGetAllUsersRequest(c *gin.Context, handler UserHandlerHelper) {
	users, err := handler.handleGetAllUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// handleUserRequest is a router for different user operations.
func handleUserRequest(c *gin.Context, handler UserHandlerHelper, action string) {
	switch action {
	case "get":
		handleGetUserRequest(c, handler)
	case "create":
		handleCreateUserRequest(c, handler)
	case "update":
		handleUpdateUserRequest(c, handler)
	case "delete":
		handleDeleteUserRequest(c, handler)
	case "getAll":
		handleGetAllUsersRequest(c, handler)
	}
}

// TestUserHandler はテスト用のハンドラー実装です.
type TestUserHandler struct {
	interactor *MockUserInteractor
}

// NewTestUserHandler はテスト用のハンドラーを作成します.
func NewTestUserHandler(mi *MockUserInteractor) *TestUserHandler {
	return &TestUserHandler{
		interactor: mi,
	}
}

// handleGetUser はユーザー取得処理を行います.
func (h *TestUserHandler) handleGetUser(c *gin.Context, id int) (*domainUser.User, error) {
	return h.interactor.GetUser(c, id)
}

// handleCreateUser はユーザー作成処理を行います.
func (h *TestUserHandler) handleCreateUser(c *gin.Context, userDTO *domainUser.User) (*domainUser.User, error) {
	return h.interactor.CreateUser(c, userDTO)
}

// handleUpdateUser はユーザー更新処理を行います.
func (h *TestUserHandler) handleUpdateUser(c *gin.Context, id int, userDTO *domainUser.User) (*domainUser.User, error) {
	return h.interactor.UpdateUser(c, userDTO)
}

// handleDeleteUser はユーザー削除処理を行います.
func (h *TestUserHandler) handleDeleteUser(c *gin.Context, id int) error {
	return h.interactor.DeleteUser(c, id)
}

// handleGetAllUsers は全ユーザー取得処理を行います.
func (h *TestUserHandler) handleGetAllUsers(c *gin.Context) ([]*domainUser.User, error) {
	return h.interactor.GetAllUsers(c)
}

// GetUser はGET /users/:id のハンドラーです.
func (h *TestUserHandler) GetUser(c *gin.Context) {
	handleUserRequest(c, h, "get")
}

// CreateUser はPOST /users のハンドラーです.
func (h *TestUserHandler) CreateUser(c *gin.Context) {
	handleUserRequest(c, h, "create")
}

// UpdateUser はPUT /users/:id のハンドラーです.
func (h *TestUserHandler) UpdateUser(c *gin.Context) {
	handleUserRequest(c, h, "update")
}

// DeleteUser はDELETE /users/:id のハンドラーです.
func (h *TestUserHandler) DeleteUser(c *gin.Context) {
	handleUserRequest(c, h, "delete")
}

// GetAllUsers はGET /users のハンドラーです.
func (h *TestUserHandler) GetAllUsers(c *gin.Context) {
	handleUserRequest(c, h, "getAll")
}

// setupTestRouter はテスト用のルーターを作成します.
func setupTestRouter(mockInteractor *MockUserInteractor) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// テスト用UserHandlerの作成
	userHandler := NewTestUserHandler(mockInteractor)

	// ルートの設定
	v1 := r.Group("/api/v1")
	users := v1.Group("/users")
	users.GET("", userHandler.GetAllUsers)
	users.POST("", userHandler.CreateUser)
	users.GET("/:id", userHandler.GetUser)
	users.PUT("/:id", userHandler.UpdateUser)
	users.DELETE("/:id", userHandler.DeleteUser)

	return r
}

// performRequest はテスト用のHTTPリクエストを実行します.
func performRequest(r http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer

	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			// panicの代わりにエラーを返す
			return httptest.NewRecorder()
		}

		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestGetUser(t *testing.T) {
	// テスト用のモックインタラクター
	mockInteractor := NewMockUserInteractor()

	// テスト用のユーザーデータ
	testUser := &domainUser.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	}
	mockInteractor.users[testUser.ID] = testUser

	// テスト用のルーター
	router := setupTestRouter(mockInteractor)

	t.Run("get existing user", func(t *testing.T) {
		// リクエストの実行
		w := performRequest(router, "GET", "/api/v1/users/"+strconv.Itoa(testUser.ID), nil)

		// レスポンスの検証
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// レスポンスボディをデコード
		var response domainUser.User

		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// ユーザーデータの検証
		if response.ID != testUser.ID {
			t.Errorf("Expected user ID %d, got %d", testUser.ID, response.ID)
		}

		if response.Name != testUser.Name {
			t.Errorf("Expected user name %s, got %s", testUser.Name, response.Name)
		}

		if response.Email != testUser.Email {
			t.Errorf("Expected user email %s, got %s", testUser.Email, response.Email)
		}
	})

	t.Run("get non-existing user", func(t *testing.T) {
		// 存在しないユーザーIDでリクエスト
		w := performRequest(router, "GET", "/api/v1/users/999", nil)

		// 存在しないユーザーの場合はステータスコード404が返る
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("get user with error", func(t *testing.T) {
		// エラーの設定
		mockInteractor.SetError(errors.New("database error"))

		// リクエストの実行
		w := performRequest(router, "GET", "/api/v1/users/1", nil)

		// エラーレスポンスの検証
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}

		// クリーンアップ
		mockInteractor.SetError(nil)
	})

	t.Run("get user with invalid ID", func(t *testing.T) {
		// 無効なIDでリクエスト
		w := performRequest(router, "GET", "/api/v1/users/invalid", nil)

		// 無効なIDでのエラーレスポンスの検証
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestCreateUser(t *testing.T) {
	// テスト用のモックインタラクター
	mockInteractor := NewMockUserInteractor()

	// テスト用のルーター
	router := setupTestRouter(mockInteractor)

	t.Run("create valid user", func(t *testing.T) {
		// 新規ユーザーデータ
		newUser := &domainUser.User{
			Name:  "New User",
			Email: "new@example.com",
		}

		// リクエストの実行
		w := performRequest(router, "POST", "/api/v1/users", newUser)

		// レスポンスの検証
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}

		// レスポンスボディをデコード
		var response domainUser.User

		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// 作成されたユーザーの検証
		if response.ID == 0 {
			t.Error("Expected user ID to be set, got 0")
		}

		if response.Name != newUser.Name {
			t.Errorf("Expected user name %s, got %s", newUser.Name, response.Name)
		}

		if response.Email != newUser.Email {
			t.Errorf("Expected user email %s, got %s", newUser.Email, response.Email)
		}
	})

	t.Run("create user with error", func(t *testing.T) {
		// エラーの設定
		mockInteractor.SetError(errors.New("database error"))

		// 新規ユーザーデータ
		newUser := &domainUser.User{
			Name:  "Error User",
			Email: "error@example.com",
		}

		// リクエストの実行
		w := performRequest(router, "POST", "/api/v1/users", newUser)

		// エラーレスポンスの検証
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}

		// クリーンアップ
		mockInteractor.SetError(nil)
	})

	t.Run("create user with invalid data", func(t *testing.T) {
		// 無効なデータ（JSON形式ではないデータ）でリクエスト
		req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer([]byte("invalid data")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 無効なデータでのエラーレスポンスの検証
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestGetAllUsers(t *testing.T) {
	// テスト用のモックインタラクター
	mockInteractor := NewMockUserInteractor()

	// テスト用のユーザーデータ
	user1 := &domainUser.User{ID: 1, Name: "User 1", Email: "user1@example.com"}
	user2 := &domainUser.User{ID: 2, Name: "User 2", Email: "user2@example.com"}
	mockInteractor.users[user1.ID] = user1
	mockInteractor.users[user2.ID] = user2

	// テスト用のルーター
	router := setupTestRouter(mockInteractor)

	t.Run("get all users", func(t *testing.T) {
		// リクエストの実行
		w := performRequest(router, "GET", "/api/v1/users", nil)

		// レスポンスの検証
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// レスポンスボディをデコード
		var response []*domainUser.User

		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// ユーザーリストの検証
		if len(response) != 2 {
			t.Errorf("Expected 2 users, got %d", len(response))
		}

		// IDでユーザーを確認
		foundIDs := make(map[int]bool)
		for _, u := range response {
			foundIDs[u.ID] = true
		}

		if !foundIDs[1] || !foundIDs[2] {
			t.Errorf("Expected users with IDs 1 and 2, got %v", foundIDs)
		}
	})

	t.Run("get all users with error", func(t *testing.T) {
		// エラーの設定
		mockInteractor.SetError(errors.New("database error"))

		// リクエストの実行
		w := performRequest(router, "GET", "/api/v1/users", nil)

		// エラーレスポンスの検証
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}

		// クリーンアップ
		mockInteractor.SetError(nil)
	})
}

func testUpdateExistingUser(t *testing.T, router *gin.Engine, mockInteractor *MockUserInteractor, existingUser *domainUser.User) {
	t.Helper()
	// 更新ユーザーデータ
	updatedUser := &domainUser.User{
		ID:    existingUser.ID,
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	// リクエストの実行
	w := performRequest(router, "PUT", "/api/v1/users/"+strconv.Itoa(existingUser.ID), updatedUser)

	// レスポンスの検証
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// レスポンスボディをデコード
	var response domainUser.User

	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// 更新されたユーザーの検証
	if response.ID != updatedUser.ID {
		t.Errorf("Expected user ID %d, got %d", updatedUser.ID, response.ID)
	}

	if response.Name != updatedUser.Name {
		t.Errorf("Expected user name %s, got %s", updatedUser.Name, response.Name)
	}

	if response.Email != updatedUser.Email {
		t.Errorf("Expected user email %s, got %s", updatedUser.Email, response.Email)
	}

	// 実際にモックのデータが更新されたか確認
	if mockInteractor.users[existingUser.ID].Name != updatedUser.Name {
		t.Errorf("Expected user name in mock to be %s, got %s", updatedUser.Name, mockInteractor.users[existingUser.ID].Name)
	}
}

func testUpdateNonExistingUser(t *testing.T, router *gin.Engine) {
	t.Helper()
	// 存在しないユーザーデータ
	nonExistingUser := &domainUser.User{
		ID:    999,
		Name:  "Non Existing User",
		Email: "nonexisting@example.com",
	}

	// リクエストの実行
	w := performRequest(router, "PUT", "/api/v1/users/999", nonExistingUser)

	// 存在しないユーザーの更新エラーの検証
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func testUpdateUserWithError(t *testing.T, router *gin.Engine, mockInteractor *MockUserInteractor, existingUser *domainUser.User) {
	t.Helper()
	// エラーの設定
	mockInteractor.SetError(errors.New("database error"))

	// 更新ユーザーデータ
	updatedUser := &domainUser.User{
		ID:    existingUser.ID,
		Name:  "Error User",
		Email: "error@example.com",
	}

	// リクエストの実行
	w := performRequest(router, "PUT", "/api/v1/users/"+strconv.Itoa(existingUser.ID), updatedUser)

	// エラーレスポンスの検証
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// クリーンアップ
	mockInteractor.SetError(nil)
}

func testUpdateUserWithInvalidID(t *testing.T, router *gin.Engine) {
	t.Helper()
	// 無効なIDでリクエスト
	updatedUser := &domainUser.User{
		Name:  "Invalid ID",
		Email: "invalid@example.com",
	}
	w := performRequest(router, "PUT", "/api/v1/users/invalid", updatedUser)

	// 無効なIDでのエラーレスポンスの検証
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func testUpdateUserWithInvalidData(t *testing.T, router *gin.Engine) {
	t.Helper()
	// 無効なデータ（JSON形式ではないデータ）でリクエスト
	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer([]byte("invalid data")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 無効なデータでのエラーレスポンスの検証
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	// テスト用のモックインタラクター
	mockInteractor := NewMockUserInteractor()

	// テスト用のユーザーデータ
	existingUser := &domainUser.User{
		ID:    1,
		Name:  "Existing User",
		Email: "existing@example.com",
	}
	mockInteractor.users[existingUser.ID] = existingUser

	// テスト用のルーター
	router := setupTestRouter(mockInteractor)

	t.Run("update existing user", func(t *testing.T) {
		testUpdateExistingUser(t, router, mockInteractor, existingUser)
	})

	t.Run("update non-existing user", func(t *testing.T) {
		testUpdateNonExistingUser(t, router)
	})

	t.Run("update user with error", func(t *testing.T) {
		testUpdateUserWithError(t, router, mockInteractor, existingUser)
	})

	t.Run("update user with invalid ID", func(t *testing.T) {
		testUpdateUserWithInvalidID(t, router)
	})

	t.Run("update user with invalid data", func(t *testing.T) {
		testUpdateUserWithInvalidData(t, router)
	})
}

func TestDeleteUser(t *testing.T) {
	// テスト用のモックインタラクター
	mockInteractor := NewMockUserInteractor()

	// テスト用のユーザーデータ
	existingUser := &domainUser.User{
		ID:    1,
		Name:  "User to Delete",
		Email: "delete@example.com",
	}
	mockInteractor.users[existingUser.ID] = existingUser

	// テスト用のルーター
	router := setupTestRouter(mockInteractor)

	t.Run("delete existing user", func(t *testing.T) {
		// リクエストの実行
		w := performRequest(router, "DELETE", "/api/v1/users/"+strconv.Itoa(existingUser.ID), nil)

		// レスポンスの検証
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// レスポンスボディをデコード
		var response map[string]string

		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// 削除成功メッセージの検証
		if message, exists := response["message"]; !exists || message != "User deleted successfully" {
			t.Errorf("Expected message 'User deleted successfully', got %v", response)
		}

		// 実際にモックからユーザーが削除されたか確認
		if _, exists := mockInteractor.users[existingUser.ID]; exists {
			t.Errorf("Expected user to be deleted from mock, but it still exists")
		}
	})

	t.Run("delete non-existing user", func(t *testing.T) {
		// 存在しないユーザーIDでリクエスト
		w := performRequest(router, "DELETE", "/api/v1/users/999", nil)

		// 存在しないユーザーの削除エラーの検証
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("delete user with error", func(t *testing.T) {
		// 再度テスト用ユーザーを追加
		mockInteractor.users[existingUser.ID] = existingUser

		// エラーの設定
		mockInteractor.SetError(errors.New("database error"))

		// リクエストの実行
		w := performRequest(router, "DELETE", "/api/v1/users/"+strconv.Itoa(existingUser.ID), nil)

		// エラーレスポンスの検証
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}

		// クリーンアップ
		mockInteractor.SetError(nil)
	})

	t.Run("delete user with invalid ID", func(t *testing.T) {
		// 無効なIDでリクエスト
		w := performRequest(router, "DELETE", "/api/v1/users/invalid", nil)

		// 無効なIDでのエラーレスポンスの検証
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}
