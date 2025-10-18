#!/bin/bash
# Скрипт для проверки Docker-конфигурации

echo "========================================="
echo "Проверка Docker-конфигурации"
echo "========================================="
echo ""

# Проверка Docker
echo "1. Проверка Docker..."
if command -v docker &> /dev/null; then
    echo "✅ Docker установлен: $(docker --version)"
else
    echo "❌ Docker не установлен"
    echo "   Установите Docker: https://www.docker.com/products/docker-desktop"
    exit 1
fi

echo ""

# Проверка Docker Compose
echo "2. Проверка Docker Compose..."
if command -v docker-compose &> /dev/null; then
    echo "✅ Docker Compose установлен: $(docker-compose --version)"
else
    echo "❌ Docker Compose не установлен"
    exit 1
fi

echo ""

# Проверка файлов
echo "3. Проверка файлов конфигурации..."

if [ -f "Dockerfile" ]; then
    echo "✅ Dockerfile найден"
else
    echo "❌ Dockerfile не найден"
    exit 1
fi

if [ -f "docker-compose.yml" ]; then
    echo "✅ docker-compose.yml найден"
else
    echo "❌ docker-compose.yml не найден"
    exit 1
fi

if [ -f "requirements.txt" ]; then
    echo "✅ requirements.txt найден"
else
    echo "❌ requirements.txt не найден"
    exit 1
fi

echo ""

# Проверка директории данных
echo "4. Проверка директории данных..."
if [ -d "data" ]; then
    echo "✅ Директория data/ существует"
    
    # Подсчет файлов
    json_files=$(find data -name "*.json" 2>/dev/null | wc -l)
    txt_files=$(find data/resumes_txt -name "*.txt" 2>/dev/null | wc -l)
    
    echo "   - JSON файлов: $json_files"
    echo "   - Резюме (TXT): $txt_files"
else
    echo "⚠️  Директория data/ не найдена. Будет создана при запуске."
fi

echo ""

# Проверка порта 5000
echo "5. Проверка доступности порта 5000..."
if lsof -Pi :5000 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "⚠️  Порт 5000 уже занят!"
    echo "   Процесс на порту 5000:"
    lsof -i :5000
    echo ""
    echo "   Остановите процесс или измените порт в docker-compose.yml"
else
    echo "✅ Порт 5000 свободен"
fi

echo ""
echo "========================================="
echo "Проверка завершена!"
echo "========================================="
echo ""
echo "Для запуска системы выполните:"
echo "  docker-compose up -d"
echo ""
echo "Для просмотра логов:"
echo "  docker-compose logs -f"
echo ""
echo "Система будет доступна по адресу:"
echo "  http://localhost:5000"
echo ""
