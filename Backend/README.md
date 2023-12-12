# Backend
## Env параметры
Все env параметры сохранены в `/internal/config/flags`. Там стоит расположить любой yaml файл и его название указать в env переменной `ENV_FILE_NAME`.

Если эта переменная пустая, то будет по умолчанию выбрано название файла `values_local.yaml`.

Пример содержания yaml файла с env параметрами:
```yaml
env:
  - name: db_dsn
    value: postgres://postgres:postgres@localhost:5432/manga_parser

  - name: listen_port
    value: 80
```