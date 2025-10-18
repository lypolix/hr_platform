# Makefile для Technopolis HR System

.PHONY: help build up down logs restart clean prune test

help:
	@echo "Technopolis HR System - Docker Commands"
	@echo ""
	@echo "Основные команды:"
	@echo "  make build     - Собрать Docker образ"
	@echo "  make up        - Запустить контейнер"
	@echo "  make down      - Остановить контейнер"
	@echo "  make restart   - Перезапустить контейнер"
	@echo "  make logs      - Показать логи"
	@echo "  make clean     - Остановить и удалить контейнер"
	@echo "  make prune     - Очистить неиспользуемые образы"
	@echo "  make rebuild   - Пересобрать с нуля"
	@echo "  make status    - Показать статус контейнера"
	@echo ""

build:
	@echo "🔨 Сборка Docker образа..."
	docker-compose build

up:
	@echo "🚀 Запуск контейнера..."
	docker-compose up -d
	@echo "✅ Система запущена на http://localhost:5000"

down:
	@echo "🛑 Остановка контейнера..."
	docker-compose down

logs:
	@echo "📋 Логи контейнера..."
	docker-compose logs -f

restart:
	@echo "🔄 Перезапуск контейнера..."
	docker-compose restart
	@echo "✅ Контейнер перезапущен"

clean:
	@echo "🧹 Остановка и удаление контейнера..."
	docker-compose down -v
	@echo "✅ Контейнер удален"

prune:
	@echo "🧹 Очистка неиспользуемых образов..."
	docker image prune -f
	@echo "✅ Очистка завершена"

rebuild:
	@echo "🔨 Пересборка с нуля..."
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d
	@echo "✅ Пересборка завершена. Система на http://localhost:5000"

status:
	@echo "📊 Статус контейнера:"
	docker ps -a | grep technopolis-hr-system || echo "Контейнер не запущен"
	@echo ""
	@echo "📊 Использование ресурсов:"
	docker stats --no-stream technopolis-hr-system 2>/dev/null || echo "Контейнер не запущен"

shell:
	@echo "🐚 Вход в контейнер..."
	docker exec -it technopolis-hr-system /bin/bash

backup:
	@echo "💾 Создание бэкапа данных..."
	@mkdir -p backups
	@cp -r data backups/data_$$(date +%Y%m%d_%H%M%S)
	@echo "✅ Бэкап создан в backups/"

dev:
	@echo "🔧 Запуск в режиме разработки..."
	@export FLASK_ENV=development && docker-compose up

prod:
	@echo "🚀 Запуск в production режиме..."
	@export FLASK_ENV=production && docker-compose up -d
	@echo "✅ Production режим запущен на http://localhost:5000"
