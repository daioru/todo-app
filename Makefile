PHONY: .run
run:
	go run cmd/todo-app/main.go

PHONY: .migrate
migrate:
	go run cmd/migration/main.go

.PHONY: .build
build: swag
	CGO_ENABLED=0  go build \
		-tags='no_mysql no_sqlite3' \
		-o ./bin/todo-app$(shell go env GOEXE) ./cmd/todo-app/main.go

.PHONY: .swag
swag:
	swag init -d cmd/todo-app,internal/handlers,internal/models

.PHONE: .migrate_build
migrate_build:
	CGO_ENABLED=0  go build \
		-tags='no_mysql no_sqlite3' \
		-o ./bin/migration$(shell go env GOEXE) ./cmd/migration/main.go