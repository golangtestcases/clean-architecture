# Clean Architecture Template

Этот проект представляет собой шаблон для создания приложений на Go с использованием принципов Clean Architecture.

## Структура проекта

```
├── cmd/
│   └── server/
│       └── main.go                    # Точка входа в приложение
├── internal/
│   ├── app/
│   │   ├── handlers/                  # HTTP хендлеры (слой представления)
│   │   │   ├── create_entity_handler/ # Отдельная папка для каждого маршрута
│   │   │   │   ├── create_entity_handler.go  # Хендлер + интерфейс + ServeHTTP
│   │   │   │   ├── request.go         # DTO для входящих данных
│   │   │   │   └── response.go        # DTO для исходящих данных
│   │   │   └── get_entity_handler/
│   │   │       ├── get_entity_handler.go
│   │   │       └── response.go
│   │   └── app.go                     # Инициализация приложения + bootstrapHandler
│   ├── domain/
│   │   ├── model/                     # Доменные модели (бизнес-сущности)
│   │   │   ├── entity.go
│   │   │   └── user.go
│   │   └── entities/                  # Бизнес-логика для сущности
│   │       ├── repository/            # Интерфейсы и реализации репозиториев
│   │       │   └── repository.go
│   │       └── service/               # Бизнес-сервисы
│   │           └── service.go
│   └── infra/                         # Инфраструктурный слой
│       └── config/
│           ├── config.go              # Конфигурация + LoadConfig
│           └── http/
│               └── middlewares/
│                   └── timer_middleware.go
└── go.mod
```

## Принципы архитектуры

### 1. Слои архитектуры

- **cmd/** - Точка входа в приложение
- **internal/app/** - Слой приложения (хендлеры, роутинг)
- **internal/domain/** - Доменный слой (бизнес-логика, модели)
- **internal/infra/** - Инфраструктурный слой (конфигурация, внешние зависимости)

### 2. Dependency Inversion

Зависимости направлены внутрь к доменному слою. Внешние слои зависят от внутренних, но не наоборот.

### 3. Разделение моделей

- **Доменные модели** (`internal/domain/model`) - бизнес-сущности, не зависят от внешних слоев
- **DTO модели** (`handlers/*/request.go`, `response.go`) - для HTTP API, содержат теги JSON и валидацию

### 4. Паттерн хендлеров

Каждый маршрут в отдельной папке с:
- Интерфейсом сервиса (только нужные методы)
- Структурой хендлера
- Конструктором `NewXxxHandler()`
- Методом `ServeHTTP()`

## Использование шаблона

1. Замените `github.com/IT-Nick/clean-architecture-template` на ваш модуль
2. Переименуйте `Entity` на вашу доменную сущность
3. Адаптируйте модели под ваши требования
4. Добавьте необходимые зависимости в go.mod

## Запуск

```bash
go run cmd/server/main.go
```

## Расширение

Для добавления нового маршрута:

1. Создайте папку `internal/app/handlers/{action}_{entity}_handler/`
2. Добавьте файлы:
   - `{action}_{entity}_handler.go` - хендлер с интерфейсом и ServeHTTP
   - `request.go` - DTO для входящих данных (если нужно)
   - `response.go` - DTO для исходящих данных
3. Зарегистрируйте в `internal/app/app.go` в функции `bootstrapHandler()`

Для добавления новой сущности:

1. Создайте модель в `internal/domain/model/`
2. Создайте папку `internal/domain/{entities}/`
3. Добавьте `repository/repository.go` и `service/service.go`
4. Создайте хендлеры для нужных операций