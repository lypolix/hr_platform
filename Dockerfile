FROM python:3.11-slim

# Установка системных зависимостей
RUN apt-get update && apt-get install -y \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Рабочая директория
WORKDIR /app

# Копирование requirements
COPY requirements.txt .

# Установка Python зависимостей
RUN pip install --no-cache-dir -r requirements.txt

# Загрузка NLTK данных
RUN python -c "import nltk; nltk.download('punkt'); nltk.download('stopwords'); nltk.download('averaged_perceptron_tagger')"

# Копирование кода приложения
COPY app.py .
COPY resume_parser.py .
COPY resume_scorer.py .
COPY industrial_skills.py .
COPY templates/ templates/
COPY static/ static/

# Создание директории для данных
RUN mkdir -p data/resumes_txt

# Копирование данных (если существуют)
COPY data/ data/

# Открытие порта
EXPOSE 5000

# Переменные окружения
ENV FLASK_APP=app.py
ENV PYTHONUNBUFFERED=1
ENV HOST=0.0.0.0
ENV PORT=5000

# Healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD python -c "import requests; requests.get('http://localhost:5000/api/statistics')" || exit 1

# Запуск приложения
CMD ["python", "app.py"]
