# Docker - Запуск системы в контейнере

## Быстрый старт

### 1. Запуск с Docker Compose (рекомендуется)

```bash
# Сборка и запуск
docker-compose up -d

# Просмотр логов
docker-compose logs -f

# Остановка
docker-compose down
```

Приложение будет доступно по адресу: `http://localhost:5000`

### 2. Запуск с Docker (без compose)

```bash
# Сборка образа
docker build -t technopolis-hr-system .

# Запуск контейнера
docker run -d \
  -p 5000:5000 \
  -v $(pwd)/data:/app/data \
  --name technopolis-hr \
  technopolis-hr-system

# Просмотр логов
docker logs -f technopolis-hr

# Остановка контейнера
docker stop technopolis-hr

# Удаление контейнера
docker rm technopolis-hr
```

---

## Конфигурация

### Переменные окружения

Вы можете настроить приложение через переменные окружения:

```yaml
environment:
  - FLASK_APP=app.py
  - FLASK_ENV=production      # development или production
  - PYTHONUNBUFFERED=1         # Для вывода логов в реальном времени
  - HOST=0.0.0.0               # Хост для прослушивания
  - PORT=5000                  # Порт приложения
```

### Порты

По умолчанию приложение работает на порту **5000**.

Изменить порт в `docker-compose.yml`:
```yaml
ports:
  - "8080:5000"  # Внешний порт 8080 → внутренний 5000
```

### Volumes (тома)

Директория `data/` монтируется как том для персистентности данных:

```yaml
volumes:
  - ./data:/app/data
```

Это означает, что:
- Изменения в `data/` сохраняются между перезапусками
- Вы можете редактировать файлы данных на хосте
- Резюме и статистика не теряются при пересоздании контейнера

---

## Проверка работоспособности

### Health Check

Контейнер автоматически проверяет свое состояние каждые 30 секунд:

```bash
# Проверить статус контейнера
docker ps

# Если контейнер healthy, вы увидите:
# STATUS: Up X minutes (healthy)
```

### Ручная проверка

```bash
# Тест API
curl http://localhost:5000/api/statistics

# Проверка логов
docker logs technopolis-hr-system

# Вход в контейнер
docker exec -it technopolis-hr-system /bin/bash
```

---

## Режимы работы

### Development (разработка)

```bash
# В docker-compose.yml
environment:
  - FLASK_ENV=development

# Или через docker run
docker run -d \
  -p 5000:5000 \
  -v $(pwd)/data:/app/data \
  -e FLASK_ENV=development \
  --name technopolis-hr \
  technopolis-hr-system
```

**Особенности:**
- Debug режим включен
- Автоматическая перезагрузка при изменении кода (если смонтирован код)
- Подробные логи ошибок

### Production (продакшн)

```bash
# В docker-compose.yml (по умолчанию)
environment:
  - FLASK_ENV=production

# Или через docker run
docker run -d \
  -p 5000:5000 \
  -v $(pwd)/data:/app/data \
  -e FLASK_ENV=production \
  --name technopolis-hr \
  technopolis-hr-system
```

**Особенности:**
- Debug режим выключен
- Оптимизированная производительность
- Безопасные настройки

---

## Управление данными

### Резервное копирование

```bash
# Копирование данных из контейнера
docker cp technopolis-hr-system:/app/data ./backup

# Или через volume
cp -r ./data ./backup_$(date +%Y%m%d)
```

### Восстановление данных

```bash
# Восстановление из бэкапа
docker cp ./backup/. technopolis-hr-system:/app/data

# Перезапуск для применения изменений
docker restart technopolis-hr-system
```

### Очистка данных

```bash
# Остановить контейнер
docker-compose down

# Удалить данные
rm -rf ./data/*.json

# Запустить заново
docker-compose up -d
```

---

## Мониторинг

### Просмотр логов

```bash
# Последние 100 строк
docker logs --tail 100 technopolis-hr-system

# С отслеживанием (follow)
docker logs -f technopolis-hr-system

# С временными метками
docker logs -t technopolis-hr-system
```

### Использование ресурсов

```bash
# Статистика ресурсов
docker stats technopolis-hr-system

# Вывод:
# CONTAINER ID   NAME                       CPU %     MEM USAGE / LIMIT
# abc123         technopolis-hr-system      0.50%     150MiB / 4GiB
```

---

## Обновление

### Пересборка образа

```bash
# Остановить и удалить контейнер
docker-compose down

# Пересобрать с нуля (без кэша)
docker-compose build --no-cache

# Запустить обновленную версию
docker-compose up -d
```

### Обновление зависимостей

```bash
# Обновить requirements.txt
# Затем пересобрать образ
docker-compose build
docker-compose up -d
```

---

## Troubleshooting

### Контейнер не запускается

```bash
# Проверить логи
docker logs technopolis-hr-system

# Проверить состояние
docker ps -a

# Попробовать запустить интерактивно
docker run -it --rm technopolis-hr-system /bin/bash
```

### Порт уже занят

```bash
# Найти процесс на порту 5000
# Windows PowerShell:
netstat -ano | findstr :5000

# Linux/Mac:
lsof -i :5000

# Изменить порт в docker-compose.yml
ports:
  - "5001:5000"
```

### Проблемы с NLTK данными

```bash
# Войти в контейнер
docker exec -it technopolis-hr-system /bin/bash

# Скачать данные вручную
python -c "import nltk; nltk.download('punkt'); nltk.download('stopwords')"

# Выйти
exit
```

### Проблемы с томами (volumes)

```bash
# Удалить все тома
docker-compose down -v

# Пересоздать
docker-compose up -d
```

---

## Полезные команды

```bash
# Список всех контейнеров
docker ps -a

# Список образов
docker images

# Удалить неиспользуемые образы
docker image prune

# Удалить всё неиспользуемое
docker system prune -a

# Перезапустить контейнер
docker restart technopolis-hr-system

# Остановить все контейнеры
docker stop $(docker ps -q)

# Удалить все контейнеры
docker rm $(docker ps -aq)
```

---

## Docker на Windows

### Использование PowerShell

```powershell
# Сборка
docker-compose up -d

# Просмотр логов
docker-compose logs -f

# Остановка
docker-compose down

# Монтирование папки в Windows
volumes:
  - C:\Users\User\data:/app/data
```

### Docker Desktop

1. Установить [Docker Desktop for Windows](https://www.docker.com/products/docker-desktop)
2. Включить WSL 2 integration (в настройках)
3. Запустить через PowerShell или Docker Desktop UI

---

## Безопасность

### Рекомендации для production

1. **Не запускайте от root:**
   ```dockerfile
   RUN useradd -m -u 1000 appuser
   USER appuser
   ```

2. **Ограничьте ресурсы:**
   ```yaml
   deploy:
     resources:
       limits:
         cpus: '1'
         memory: 1G
   ```

3. **Используйте secrets для паролей:**
   ```yaml
   secrets:
     - api_key
   ```

4. **Обновляйте базовый образ:**
   ```bash
   docker pull python:3.11-slim
   docker-compose build --pull
   ```

---

## Документация

- [Официальная документация Docker](https://docs.docker.com/)
- [Docker Compose документация](https://docs.docker.com/compose/)
- [Best practices для Dockerfile](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)

---

## Поддержка

При возникновении проблем:

1. Проверьте логи: `docker logs technopolis-hr-system`
2. Убедитесь, что порт 5000 свободен
3. Проверьте наличие файлов данных в `./data/`
4. Пересоберите образ: `docker-compose build --no-cache`
