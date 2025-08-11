# プロジェクト内のテスト対象パッケージ一覧（cmd配下は除外）
COVERPKGS := $(shell go list ./... | grep -v '/cmd')

# デフォルトターゲット（単純なテスト）
test:
	go test $(COVERPKGS)

# カバレッジ出力（テキスト + HTML）
cover:
	go test -coverprofile=coverage.out $(COVERPKGS)
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
