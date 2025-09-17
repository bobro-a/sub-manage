# Subscription Manager Service

Сервис для управления подписками пользователей: создание, чтение, обновление и удаление (CRUDL).  
Использует PostgreSQL (в Docker), конфигурации YAML и миграции

---

## Структура проекта
```
sub-manage/
├── cmd/sub-manage/            # main.go (точка входа)
│       └── main.go
├── configs/                   # Конфиги
│       └── config.yaml
│       └── config-docker.yaml
├── docs/                        # swagger документация
├── docker-compose.yml           # Поднятие БД + сервиса
├── Dockerfile                   # Билд образа приложения
├── go.mod / go.sum              # Go модули
├── samples/                     # Примеры вызовов curl
└── internal/                    # Внутренние пакеты
│   ├── app/
│   │   └── app.go
│   ├── handler/
│   │   └── handler.go
│   ├── migrations/
│   │   └── subscriptions.up.sql
│   ├── repo/
│   │   └── repo.go
│   └── usecase/
│       ├── usecase.go
├── pkg/                         # Пакеты для конфигураций и моделей
│   ├── config/
│   │   └── config.go
│   └── models/
│       └── sub.go
│       └── filter.go
```
---

## Запуск
```bash
docker-compose up --build -d
```
---

## Проверка
Пример вызова через curl [`samples/command.md`](samples/command.md)

##  Миграции

SQL-файл в `internal/migrations/`.  
Применяется при запуске `main.go` или внутри контейнера.

---

## Очистка
```bash
docker-compose down -v
```
