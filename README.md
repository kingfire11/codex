# designapi-instructions

Статический сайт-инструкция (GitHub Pages) для подключения Codex / OpenCode к `https://api.designapi.ink`.

## Локально

```
cd designapi-instructions
python3 -m http.server 8000
# http://localhost:8000
```

## Деплой на GitHub Pages

1. Создать репо `designapi-instructions`, запушить содержимое этой папки в ветку `main`.
2. Settings → Pages → Source: `Deploy from a branch` → `main` / `/ (root)`.
3. Файл `.nojekyll` уже на месте.
4. Опционально: Custom domain (положить файл `CNAME` с доменом).

## Как это работает

- Никакого backend. Поле «токен» — обычный `<input>`.
- Шаблоны в `scripts/*` грузятся через `fetch()`. JS делает `replaceAll('__API_KEY__', token)` и показывает результат / отдаёт через Blob-скачивание.
- Если поле пустое — в скриптах остаётся плейсхолдер `YOUR_API_KEY`.
- Однострочники подсовывают токен через переменную окружения, а не через URL — токен **не уходит** на сервер.

## Связь с инсталлером

Кнопки «Скачать инсталлер» ведут на
`https://github.com/<INSTALLER_REPO>/releases/latest/download/designapi-installer-<os>-<arch>[.exe]`.
Подправьте константу `INSTALLER_REPO` в `assets/app.js` под реальный путь к репозиторию.
