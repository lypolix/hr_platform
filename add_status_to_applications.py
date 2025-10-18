"""
Скрипт для добавления поля status к откликам
"""
import json
import random

# Читаем текущие отклики
with open('data/applications.json', 'r', encoding='utf-8') as f:
    applications = json.load(f)

# Статусы:
# - new (новое)
# - in_review (на рассмотрении / в работе)
# - rejected (отказ)
# - accepted (принято)

# Распределяем статусы для демо:
# 40% - new
# 30% - in_review  
# 20% - rejected
# 10% - accepted

statuses = ['new', 'in_review', 'rejected', 'accepted']
weights = [40, 30, 20, 10]

for app in applications:
    # Добавляем статус если его нет
    if 'status' not in app:
        # Выбираем случайный статус с весами
        status = random.choices(statuses, weights=weights, k=1)[0]
        app['status'] = status
        
        # Добавляем timestamp обновления статуса
        app['status_updated_at'] = '2025-10-18T12:00:00'

# Сохраняем обновленные данные
with open('data/applications.json', 'w', encoding='utf-8') as f:
    json.dump(applications, f, ensure_ascii=False, indent=2)

# Статистика
status_counts = {}
for app in applications:
    status = app['status']
    status_counts[status] = status_counts.get(status, 0) + 1

print("Статусы добавлены к откликам:")
print(f"Всего откликов: {len(applications)}")
print("\nРаспределение:")
print(f"  Новых (new): {status_counts.get('new', 0)}")
print(f"  На рассмотрении (in_review): {status_counts.get('in_review', 0)}")
print(f"  Отказ (rejected): {status_counts.get('rejected', 0)}")
print(f"  Принято (accepted): {status_counts.get('accepted', 0)}")
