MIGRATE=migrate
MIGRATIONS_DIR=./db/postgresql/migrations
DB_URL=postgres://postgres:122002@172.31.10.86:5432/db_mini_shop

up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

redo:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

version:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

force:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force

create:
	@read -p "Migration name: " name; \
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) $$name
