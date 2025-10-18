# Скрипт для проверки Docker-конфигурации (Windows PowerShell)

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Проверка Docker-конфигурации" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

# Проверка Docker
Write-Host "1. Проверка Docker..." -ForegroundColor Yellow
try {
    $dockerVersion = docker --version
    Write-Host "✅ Docker установлен: $dockerVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker не установлен" -ForegroundColor Red
    Write-Host "   Установите Docker Desktop: https://www.docker.com/products/docker-desktop" -ForegroundColor Red
    exit 1
}

Write-Host ""

# Проверка Docker Compose
Write-Host "2. Проверка Docker Compose..." -ForegroundColor Yellow
try {
    $composeVersion = docker-compose --version
    Write-Host "✅ Docker Compose установлен: $composeVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker Compose не установлен" -ForegroundColor Red
    exit 1
}

Write-Host ""

# Проверка файлов
Write-Host "3. Проверка файлов конфигурации..." -ForegroundColor Yellow

if (Test-Path "Dockerfile") {
    Write-Host "✅ Dockerfile найден" -ForegroundColor Green
} else {
    Write-Host "❌ Dockerfile не найден" -ForegroundColor Red
    exit 1
}

if (Test-Path "docker-compose.yml") {
    Write-Host "✅ docker-compose.yml найден" -ForegroundColor Green
} else {
    Write-Host "❌ docker-compose.yml не найден" -ForegroundColor Red
    exit 1
}

if (Test-Path "requirements.txt") {
    Write-Host "✅ requirements.txt найден" -ForegroundColor Green
} else {
    Write-Host "❌ requirements.txt не найден" -ForegroundColor Red
    exit 1
}

Write-Host ""

# Проверка директории данных
Write-Host "4. Проверка директории данных..." -ForegroundColor Yellow
if (Test-Path "data") {
    Write-Host "✅ Директория data\ существует" -ForegroundColor Green
    
    # Подсчет файлов
    $jsonFiles = (Get-ChildItem -Path "data" -Filter "*.json" -Recurse -ErrorAction SilentlyContinue).Count
    $txtFiles = (Get-ChildItem -Path "data\resumes_txt" -Filter "*.txt" -ErrorAction SilentlyContinue).Count
    
    Write-Host "   - JSON файлов: $jsonFiles" -ForegroundColor Gray
    Write-Host "   - Резюме (TXT): $txtFiles" -ForegroundColor Gray
} else {
    Write-Host "⚠️  Директория data\ не найдена. Будет создана при запуске." -ForegroundColor Yellow
}

Write-Host ""

# Проверка порта 5000
Write-Host "5. Проверка доступности порта 5000..." -ForegroundColor Yellow
$port5000 = Get-NetTCPConnection -LocalPort 5000 -ErrorAction SilentlyContinue

if ($port5000) {
    Write-Host "⚠️  Порт 5000 уже занят!" -ForegroundColor Yellow
    Write-Host "   Процесс на порту 5000:" -ForegroundColor Yellow
    $port5000 | Select-Object LocalAddress, LocalPort, State, OwningProcess | Format-Table
    Write-Host ""
    Write-Host "   Остановите процесс или измените порт в docker-compose.yml" -ForegroundColor Yellow
} else {
    Write-Host "✅ Порт 5000 свободен" -ForegroundColor Green
}

Write-Host ""

# Проверка Docker Desktop
Write-Host "6. Проверка состояния Docker Desktop..." -ForegroundColor Yellow
try {
    docker ps > $null 2>&1
    Write-Host "✅ Docker Desktop запущен" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Docker Desktop не запущен или недоступен" -ForegroundColor Yellow
    Write-Host "   Запустите Docker Desktop перед началом работы" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Проверка завершена!" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Для запуска системы выполните:" -ForegroundColor White
Write-Host "  docker-compose up -d" -ForegroundColor Green
Write-Host ""
Write-Host "Для просмотра логов:" -ForegroundColor White
Write-Host "  docker-compose logs -f" -ForegroundColor Green
Write-Host ""
Write-Host "Система будет доступна по адресу:" -ForegroundColor White
Write-Host "  http://localhost:5000" -ForegroundColor Cyan
Write-Host ""
Write-Host "Администратор:" -ForegroundColor White
Write-Host "  http://localhost:5000/" -ForegroundColor Cyan
Write-Host ""
Write-Host "Представитель компании:" -ForegroundColor White
Write-Host "  http://localhost:5000/company?company=Технополис%20Москва" -ForegroundColor Cyan
Write-Host ""
