# Makefile for goose v3
.PHONY: db up down redo version force create

export DATABASE_URL=postgresql://postgres:123456@localhost:5432/test_db_02?sslmode=disable
export MIGRATIONS_DIR=./db/postgresql/migrations

# Create migration: make db name=create_users
db:
	@echo "Creating migration: $(name)"
	goose create $(name) sql -dir $(MIGRATIONS_DIR)

# Run all up migrations
up:
	goose -dir $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

# Rollback one step
down:
	goose -dir $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

# Redo last migration
redo:
	goose -dir $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1
	goose -dir $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

# Show current DB version
version:
	goose -dir $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

# Force migration version (interactive)
force:
	@read -p "Enter version to force to: " v; \
	goose -dir $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $$v

# Interactive migration creation
create:
	@read -p "Migration name: " name; \
	goose create $$name sql -dir $(MIGRATIONS_DIR)
