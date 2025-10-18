# Файлы для копирования
$source = "c:\Users\Елена\Documents\Хакатоны\Моспром\nlp_analytics"
$dest = "c:\Users\Елена\Documents\Хакатоны\Моспром\technopolis-hr-system"

Write-Host "Копирование файлов в новую структуру проекта..." -ForegroundColor Green

# Python файлы
Copy-Item "$source\app.py" $dest
Copy-Item "$source\resume_parser.py" $dest
Copy-Item "$source\resume_scorer.py" $dest
Copy-Item "$source\industrial_skills.py" $dest
Copy-Item "$source\requirements.txt" $dest

# Директории
Copy-Item "$source\templates" $dest -Recurse -Force
Copy-Item "$source\static" $dest -Recurse -Force
Copy-Item "$source\data" $dest -Recurse -Force

Write-Host "✓ Копирование завершено!" -ForegroundColor Green
Write-Host ""
Write-Host "Для запуска Docker контейнера:" -ForegroundColor Yellow
Write-Host "  cd $dest" -ForegroundColor Cyan
Write-Host "  docker-compose up -d" -ForegroundColor Cyan
Write-Host ""
Write-Host "Для просмотра логов:" -ForegroundColor Yellow
Write-Host "  docker-compose logs -f" -ForegroundColor Cyan
