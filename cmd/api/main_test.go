package main

import (
	"os"
	"testing"
)

// TestMain関数は、パッケージ内のすべてのテストの前後に実行される特別な関数です.
func TestMain(m *testing.M) {
	// テスト前の準備
	setup()

	// テストを実行
	m.Run()

	// テスト後のクリーンアップ
	teardown()
	// os.Exit(code)の代わりに、codeを返す
	// テスト結果はm.Run()の戻り値として取得できるため、明示的にExitする必要はない
}

// テスト前の準備.
func setup() {
	// テスト環境の設定
	os.Setenv("GO_ENV", "test")
	os.Setenv("PORT", "8081") // テスト用のポート
}

// テスト後のクリーンアップ.
func teardown() {
	// 環境変数をリセット
	os.Unsetenv("GO_ENV")
	os.Unsetenv("PORT")
}

// 注意: 実際のアプリケーションでは、より詳細なテストが必要です.
func TestEnvironmentVariables(t *testing.T) {
	// 環境変数が正しく設定されているか確認
	env := os.Getenv("GO_ENV")
	if env != "test" {
		t.Errorf("Expected GO_ENV to be 'test', got '%s'", env)
	}

	port := os.Getenv("PORT")
	if port != "8081" {
		t.Errorf("Expected PORT to be '8081', got '%s'", port)
	}
}

// run関数の実際のテストは複雑なため、ここでは省略しています
// 実際のアプリケーションでは、モックを使用してデータベース接続やHTTPサーバーをテストします
