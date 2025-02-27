PHONY: .run
run:
	go run cmd/todo-app/main.go

PHONY: .migrate
migrate:
	go run cmd/migration/main.go

.PHONY: .swag
swag:
	swag init -d cmd/todo-app,internal/handlers,internal/models