"""
Добавление уникальных ID к откликам
"""
import json
import uuid
from pathlib import Path

# Путь к файлу
applications_file = Path(__file__).parent / "data" / "applications.json"

# Загружаем отклики
with open(applications_file, 'r', encoding='utf-8') as f:
    applications = json.load(f)

# Добавляем ID к каждому отклику
for i, app in enumerate(applications, 1):
    if 'id' not in app:
        # Генерируем уникальный ID в формате app_001, app_002 и т.д.
        app['id'] = f"app_{i:03d}"

# Сохраняем обратно
with open(applications_file, 'w', encoding='utf-8') as f:
    json.dump(applications, f, ensure_ascii=False, indent=2)

print(f"✅ ID добавлены к {len(applications)} откликам")
print(f"Примеры:")
for app in applications[:3]:
    print(f"  - {app['id']}: {app['candidate_name']} → {app['vacancy_id']}")
