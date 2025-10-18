"""
Создаем общий банк резюме (резюме без привязки к вакансиям)
"""
import json
import random
from datetime import datetime, timedelta

# Читаем все резюме
with open('data/parsed_resumes.json', 'r', encoding='utf-8') as f:
    all_resumes = json.load(f)

# Выбираем 10 случайных резюме для общего банка
# (как будто они пришли на почту без отклика на вакансию)
general_pool = random.sample(all_resumes, min(10, len(all_resumes)))

# Добавляем метаданные
for resume in general_pool:
    # Дата поступления (последние 2 недели)
    days_ago = random.randint(0, 14)
    received_date = datetime.now() - timedelta(days=days_ago)
    resume['received_at'] = received_date.strftime('%Y-%m-%dT%H:%M:%S')
    resume['source'] = 'email'  # источник: email, site, etc
    resume['status'] = 'available'  # доступно для всех

# Сохраняем общий банк
with open('data/general_resume_pool.json', 'w', encoding='utf-8') as f:
    json.dump(general_pool, f, ensure_ascii=False, indent=2)

print(f"Общий банк резюме создан: {len(general_pool)} резюме")
print("\nПримеры специализаций:")
for resume in general_pool[:5]:
    specs = ', '.join(resume['specializations'])
    print(f"  - {resume['filename']}: {specs}")
