# Инструкция по развёртыванию

## Вариант 1: Docker (Рекомендуется)

### Требования
- Docker Desktop для Windows
- Docker Compose

### Шаги

1. **Перейдите в директорию проекта:**
```powershell
cd "c:\Users\Елена\Documents\Хакатоны\Моспром\technopolis-hr-system"
```

2. **Соберите и запустите контейнер:**
```powershell
docker-compose up -d --build
```

3. **Проверьте статус:**
```powershell
docker-compose ps
```

4. **Откройте в браузере:**
http://localhost:5000

5. **Просмотр логов:**
```powershell
docker-compose logs -f
```

6. **Остановка:**
```powershell
docker-compose down
```

## Вариант 2: Без Docker

### Требования
- Python 3.11+
- pip

### Шаги

1. **Создайте виртуальное окружение:**
```powershell
cd "c:\Users\Елена\Documents\Хакатоны\Моспром\technopolis-hr-system"
python -m venv .venv
.venv\Scripts\activate
```

2. **Установите зависимости:**
```powershell
pip install -r requirements.txt
```

3. **Запустите сервер:**
```powershell
python app.py
```

4. **Откройте в браузере:**
http://localhost:5000

## Первоначальная настройка

После запуска системы:

1. Нажмите кнопку "Распарсить все резюме"
2. Дождитесь завершения парсинга
3. Нажмите "Обновить статистику"
4. Система готова к работе!

## Структура данных

```
technopolis-hr-system/
├── data/
│   ├── resumes_txt/          # Папка с резюме в формате .txt
│   ├── parsed_resumes.json   # Обработанные резюме
│   ├── vacancies.json        # Вакансии
│   └── statistics.json       # Статистика
```

## Добавление новых резюме

1. Поместите файлы резюме (.txt) в папку `data/resumes_txt/`
2. На главной странице нажмите "Распарсить все резюме"
3. Система автоматически обработает новые файлы

## Обновление кода

### Docker:
```powershell
docker-compose down
docker-compose up -d --build
```

### Без Docker:
```powershell
# Просто перезапустите app.py
```

## Порты

- **5000** - Web-интерфейс и API

## Переменные окружения

В `docker-compose.yml` можно настроить:
- `FLASK_ENV` - режим работы (production/development)
- `RESUME_DIR` - путь к папке с резюме

## Troubleshooting

### Проблема: Порт 5000 занят
**Решение:** Измените порт в `docker-compose.yml`:
```yaml
ports:
  - "8000:5000"  # Внешний порт 8000
```

### Проблема: Контейнер не запускается
**Решение:** Проверьте логи:
```powershell
docker-compose logs
```

### Проблема: Не загружаются данные
**Решение:** Убедитесь, что папка `data` существует и содержит файлы

## Бэкап данных

Создайте копию папки `data`:
```powershell
Copy-Item data data_backup -Recurse
```

## Production deployment

Для production используйте:
- Nginx как reverse proxy
- SSL сертификаты
- Переменную окружения `FLASK_ENV=production`
