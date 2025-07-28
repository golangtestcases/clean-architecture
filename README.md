# Subscription Service

REST-сервис для агрегации данных об онлайн-подписках пользователей, построенный на принципах Чистой Архитектуры.

## Возможности

- CRUD операции над подписками
- Подсчет суммарной стоимости подписок с фильтрацией
- PostgreSQL с миграциями
- Структурированное логирование
- Docker Compose для развертывания
- Swagger документация

## API Endpoints

- `POST /api/subscriptions` - создание подписки
- `GET /api/subscriptions/{id}` - получение подписки по ID
- `PUT /api/subscriptions/{id}` - обновление подписки
- `DELETE /api/subscriptions/{id}` - удаление подписки
- `GET /api/subscriptions` - список подписок с пагинацией
- `GET /api/subscriptions/cost` - получение общей стоимости с фильтрами
- `GET /swagger/` - Swagger документация

### Фильтры для /api/subscriptions/cost:
- `user_id` - UUID пользователя
- `service_name` - название сервиса (частичное совпадение)
- `start_date` - дата начала периода (MM-YYYY)
- `end_date` - дата окончания периода (MM-YYYY)

## Модель данных

```json
{
  "id": "uuid",
  "service_name": "string",
  "price": "integer",
  "user_id": "uuid", 
  "start_date": "MM-YYYY",
  "end_date": "MM-YYYY"
}
```

## Запуск

### С Docker Compose
```bash
docker-compose up --build
```

### Локально
1. Создайте `.env` файл из `.env.example`
2. Запустите PostgreSQL
3. Выполните миграции
4. Запустите приложение:
```bash
go run cmd/server/main.go
```

## Конфигурация

Настройки через переменные окружения или `.env` файл:

- `SERVER_HOST` - хост сервера (по умолчанию: localhost)
- `SERVER_PORT` - порт сервера (по умолчанию: 8080)
- `DB_HOST` - хост PostgreSQL
- `DB_PORT` - порт PostgreSQL
- `DB_USER` - пользователь БД
- `DB_PASSWORD` - пароль БД
- `DB_NAME` - имя БД
- `LOG_LEVEL` - уровень логирования (debug, info, warn, error)

## Разработка

### Генерация Swagger документации
```bash
make swagger
```

### Сборка проекта
```bash
make build
```

### Запуск локально
```bash
make run
```

### Docker команды
```bash
make docker-up    # Запуск с Docker Compose
make docker-down  # Остановка контейнеров
```