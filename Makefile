# Стандартная цель по умолчанию
.PHONY: all
all: build

# Сборка проекта
.PHONY: build
build:
	docker-compose build

# Запуск контейнеров в фоновом режиме
.PHONY: up
up:
	docker-compose up -d

# Остановка и удаление контейнеров
.PHONY: down
down:
	docker-compose down

# Запуск тестов
.PHONY: test
test:
	docker-compose run --rm app go test ./... -v

# Просмотр логов
.PHONY: logs
logs:
	docker-compose logs

# Вход внутрь контейнера (замените `app` на имя вашего сервиса)
.PHONY: shell
shell:
	docker-compose exec app sh

# Очистка неиспользуемых Docker объектов (волюмы, сети, контейнеры и образы)
.PHONY: clean
clean:
	docker system prune -a
	docker volume prune

# Остановка всех контейнеров
.PHONY: stop
stop:
	docker-compose stop

# Перезапуск всех контейнеров
.PHONY: restart
restart:
	docker-compose restart

# Пересборка и перезапуск контейнеров
.PHONY: rebuild
rebuild: down build up

