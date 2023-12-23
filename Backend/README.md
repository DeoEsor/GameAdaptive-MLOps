# Backend
## Env параметры
Все env параметры сохранены в `/internal/config/flags`. Там стоит расположить любой yaml файл и его название указать в env переменной `ENV_FILE_NAME`.

Если эта переменная пустая, то будет по умолчанию выбрано название файла `values_local.yaml`.

Пример содержания yaml файла с env параметрами:
```yaml
env:
  - name: db_dsn
    value: postgres://postgres:postgres@db:5432/ml_ops

  - name: listen_port
    value: 80

  - name: kafka_brokers
    value: "kafka:29092"
```