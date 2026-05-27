# createuser CLI

Утиліта для створення користувачів у базі даних для Agro Data Management System.

## Збірка

```bash
go build ./cmd/createuser
```

## Використання

```bash
./createuser --email user@example.com --password mypass123 --role admin
```

- `--email` — email користувача (обов'язково)
- `--password` — пароль користувача (обов'язково)
- `--role` — роль користувача (необов'язково, за замовчуванням `user`)
- `--config` — шлях до конфігураційного файлу (необов'язково, за замовчуванням `config/local.yaml`)

### Приклад

```bash
./createuser --email admin@example.com --password supersecret --role admin
```

Після виконання користувач буде доданий у таблицю `users` вашої БД і зможе авторизуватись через API.
