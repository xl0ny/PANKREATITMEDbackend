# Makefile для управления docker-compose и MinIO

# Запуск сервисов
up:
	docker compose up -d

# Остановка сервисов
down:
	docker compose down

# Перезапуск сервисов
restart:
	docker compose down
	docker compose up -d

# Просмотр логов (всех сервисов)
logs:
	docker compose logs -f

# Логи конкретного сервиса: make logs-db / make logs-minio / make logs-redis
logs-%:
	docker compose logs -f $*

# Список работающих контейнеров
ps:
	docker compose ps

# Очистка (контейнеры + сети + тома)
clean:
	docker compose down -v

# --- MinIO ---
# Инициализация MinIO: создание бакета my-bucket и установка публичного доступа
init-minio:
	docker run --rm --network host minio/mc \
		alias set local http://localhost:9000 minioadmin minioadmin123
	docker run --rm --network host minio/mc mb -p local/my-bucket || true
	docker run --rm --network host minio/mc anonymous set download local/my-bucket

swagger:
	swag init -g cmd/pankreatitmed/main.go --parseInternal --parseDependency