# Статические файлы

## Как заменить логотип

1. Положите файл логотипа в эту папку (`static/logo.png` или `static/logo.svg`)
2. Откройте `templates/index.html`
3. Найдите блок `.logo` в CSS (строка ~27)
4. Замените строку:
   ```css
   background: url('data:image/svg+xml;utf8,...') center/contain no-repeat;
   ```
   на:
   ```css
   background: url('/static/logo.png') center/contain no-repeat;
   ```

## Настройка Flask для статических файлов

В `app.py` уже настроена поддержка статических файлов:
```python
app = Flask(__name__, static_folder='static', static_url_path='/static')
```

После добавления логотипа перезапустите Flask сервер.
