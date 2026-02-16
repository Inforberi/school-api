# Путь к миграциям
MIGRATIONS_DIR := migrations

# Строка подключения к базе (можно переопределить в окружении)
DATABASE_URL ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_NAME_DB)?sslmode=disable

# Флаг verbose для подробного вывода
VERBOSE ?= -verbose

# Команды
.PHONY: migrate-create migrate_up migrate_down migrate_force migrate_version

help:
	@echo "Доступные команды:"
	@echo "  make migrate-create     — создать новую миграцию"
	@echo "  make migrate_up         — применить миграции"
	@echo "  make migrate_down       — откатить миграции"
	@echo "  make migrate-force      — force установить версию"
	@echo "  make migrate_version    — показать текущую версию"

# Создать новую миграцию с интерактивным вводом имени
migrate-create:
	@read -p "Введите имя миграции (example: create_users_table): " name; \
	if [ -z "$$name" ]; then \
		echo "❌ Имя миграции не должно быть пустым"; \
		exit 1; \
	fi; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

# Применить все миграции вверх
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" $(VERBOSE) up

# Откатить все миграции вниз
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" $(VERBOSE) down

# Форсировать версию (когда база стала dirty)
migrate-force:
	@read -p "Введите версию для force (например 000003): " ver; \
	if [ -z "$$ver" ]; then \
		echo "❌ Версия не должна быть пустой"; \
		exit 1; \
	fi; \
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" $(VERBOSE) force $$ver

# Показать текущую версию миграций
migrate_version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version