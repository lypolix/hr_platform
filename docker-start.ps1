# Быстрый запуск Technopolis HR System через Docker
# Windows PowerShell

Write-Host ""
Write-Host "╔═══════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║                                                           ║" -ForegroundColor Cyan
Write-Host "║     TECHNOPOLIS HR SYSTEM - DOCKER QUICK START            ║" -ForegroundColor Cyan
Write-Host "║                                                           ║" -ForegroundColor Cyan
Write-Host "╚═══════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# Проверка Docker
Write-Host "🔍 Проверка Docker..." -ForegroundColor Yellow
try {
    docker --version | Out-Null
    Write-Host "✅ Docker найден" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker не установлен!" -ForegroundColor Red
    Write-Host "Установите Docker Desktop: https://www.docker.com/products/docker-desktop" -ForegroundColor Yellow
    Read-Host "Нажмите Enter для выхода"
    exit 1
}

Write-Host ""

# Остановка старых контейнеров
Write-Host "🛑 Остановка старых контейнеров..." -ForegroundColor Yellow
docker-compose down 2>$null
Write-Host "✅ Готово" -ForegroundColor Green

Write-Host ""

# Сборка образа
Write-Host "🔨 Сборка Docker образа..." -ForegroundColor Yellow
docker-compose build

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Ошибка сборки!" -ForegroundColor Red
    Read-Host "Нажмите Enter для выхода"
    exit 1
}

Write-Host "✅ Образ собран" -ForegroundColor Green

Write-Host ""

# Запуск контейнера
Write-Host "🚀 Запуск контейнера..." -ForegroundColor Yellow
docker-compose up -d

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Ошибка запуска!" -ForegroundColor Red
    Read-Host "Нажмите Enter для выхода"
    exit 1
}

Write-Host "✅ Контейнер запущен" -ForegroundColor Green

Write-Host ""

# Ожидание запуска
Write-Host "⏳ Ожидание инициализации (10 сек)..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

Write-Host ""
Write-Host "╔═══════════════════════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║                                                           ║" -ForegroundColor Green
Write-Host "║              🎉 СИСТЕМА УСПЕШНО ЗАПУЩЕНА! 🎉              ║" -ForegroundColor Green
Write-Host "║                                                           ║" -ForegroundColor Green
Write-Host "╚═══════════════════════════════════════════════════════════╝" -ForegroundColor Green
Write-Host ""

Write-Host "📍 ДОСТУП К СИСТЕМЕ:" -ForegroundColor Cyan
Write-Host ""
Write-Host "   1️⃣  АДМИНИСТРАТОР (HR Технополиса)" -ForegroundColor White
Write-Host "      http://localhost:5000/" -ForegroundColor Yellow
Write-Host "      • Полный доступ ко всем данным" -ForegroundColor Gray
Write-Host "      • 13 вакансий от всех компаний" -ForegroundColor Gray
Write-Host "      • 126 откликов" -ForegroundColor Gray
Write-Host "      • Статистика и аналитика" -ForegroundColor Gray
Write-Host ""
Write-Host "   2️⃣  ПРЕДСТАВИТЕЛЬ КОМПАНИИ" -ForegroundColor White
Write-Host "      http://localhost:5000/company?company=Технополис%%20Москва" -ForegroundColor Yellow
Write-Host "      • Ограниченный доступ" -ForegroundColor Gray
Write-Host "      • Только свои вакансии и отклики" -ForegroundColor Gray
Write-Host "      • Управление статусами откликов" -ForegroundColor Gray
Write-Host "      • Общий банк резюме" -ForegroundColor Gray
Write-Host ""

Write-Host "💡 ПОЛЕЗНЫЕ КОМАНДЫ:" -ForegroundColor Cyan
Write-Host ""
Write-Host "   Просмотр логов:" -ForegroundColor White
Write-Host "   docker-compose logs -f" -ForegroundColor Yellow
Write-Host ""
Write-Host "   Остановка:" -ForegroundColor White
Write-Host "   docker-compose down" -ForegroundColor Yellow
Write-Host ""
Write-Host "   Перезапуск:" -ForegroundColor White
Write-Host "   docker-compose restart" -ForegroundColor Yellow
Write-Host ""
Write-Host "   Статус контейнера:" -ForegroundColor White
Write-Host "   docker ps" -ForegroundColor Yellow
Write-Host ""

Write-Host "📖 Документация: README.md, DOCKER.md" -ForegroundColor Cyan
Write-Host ""

# Предложение открыть браузер
$openBrowser = Read-Host "Открыть систему в браузере? (y/n)"
if ($openBrowser -eq 'y' -or $openBrowser -eq 'Y') {
    Start-Process "http://localhost:5000"
}

Write-Host ""
Write-Host "Готово! 🚀" -ForegroundColor Green
Write-Host ""
