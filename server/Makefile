.PHONY: build-server
build-server:
	@go build -o bin/server cmd/server/main.go

.PHONY: generate-query
generate-query:
	@sqlc generate -f internal/user/repositories/postgres/sqlc.yaml
	@sqlc generate -f internal/post/repositories/postgres/sqlc.yaml

.PHONY: diff-query
diff-query:
	@sqlc diff -f internal/user/repositories/postgres/sqlc.yaml
	@sqlc diff -f internal/post/repositories/postgres/sqlc.yaml

.PHONY: help
help:
	@echo available commands:
	@echo   build-server