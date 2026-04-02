# Booking 

 backend-сервис для бронирования переговорных комнат с JWT-авторизацией, ролями, миграциями, Swagger и защитой от пересекающихся броней на уровне Postgres.


## Возможности

- регистрация и логин пользователей;
- выдача JWT-токена;
- получение данных текущего пользователя через `/me`;
- просмотр списка комнат и отдельной комнаты;
- создание бронирования;
- просмотр своих бронирований;
- отмена бронирования;
- создание и удаление блокировок комнаты администратором;
- просмотр занятых интервалов комнаты в заданном диапазоне;
- Swagger-документация;
- интеграционные тесты на конфликт бронирований и конкурентные запросы.

## Стек

- Go
- Gin
- PostgreSQL
- pgx / pgxpool
- golang-migrate
- JWT
- bcrypt
- Swagger (`swaggo`)
- Testcontainers для интеграционных тестов
- Docker / Docker Compose

## Как устроен проект

Проект разбит на слои, чтобы логика не смешивалась в одном месте.

- `router` — собирает зависимости и настраивает маршруты;
- `middleware` — общие проверки: JWT, роли, request id, логирование;
- `handlers` — HTTP-слой: принимает JSON, валидирует, вызывает сервис, отдаёт ответ;
- `service` — бизнес-логика и правила;
- `repository/postgres` — SQL и работа с базой;
- `domain` — внутренние сущности приложения;
- `dto` — структуры входа и выхода для HTTP;
- `migrations` — схема БД;
- `tests` — интеграционные тесты и вспомогательные утилиты.

## Структура проекта

```text
booking/
├── cmd/
│   └── api/
│       └── main.go                  # точка входа в приложение: загрузка конфига, подключение к БД, запуск миграций, старт HTTP-сервера
|
├── docs/
│   ├── docs.go                      # сгенерированный swagger-код для подключения документации
│   ├── swagger.json                 # swagger-описание API в JSON
│   └── swagger.yaml                 # swagger-описание API в YAML
|
├── internal/                        # внутренняя логика приложения, недоступная для внешнего импорта
│   ├── auth/
│   │   ├── jwt.go                   # создание и проверка JWT-токенов
│   │   └── password.go              # хеширование и проверка паролей через bcrypt
|   |
│   ├── config/
│   │   └── config.go                # чтение и хранение конфигурации приложения из env
|   |
│   ├── domain/
│   │   ├── reservation.go           # домен бронирования
│   │   ├── roomblock.go             # домен блокировки комнаты
│   │   ├── rooms.go                 # домен комнаты
│   │   └── user.go                  # домен пользователя
|   |
│   ├── dto/
│   │   ├── auth.go                  # структуры запросов и ответов для авторизации
│   │   ├── availability.go          # DTO для проверки занятости комнаты
│   │   ├── block.go                 # DTO для создания/удаления блокировок
│   │   ├── reservation.go           # DTO для бронирований
│   │   └── room.go                  # DTO для комнат
|   |
│   ├── errs/
│   │   └── domain.go                # общие доменные ошибки: conflict, not found, validation и т.д.
|   |
│   ├── httpx/
│   │   ├── router.go                # сборка роутера, middleware, handlers и всех зависимостей приложения
|   |   |
│   │   ├── handlers/
│   │   │   ├── auth.go              # HTTP-обработчики регистрации и логина
│   │   │   ├── availability.go      # HTTP-обработчик проверки занятости комнаты
│   │   │   ├── blosk.go             # HTTP-обработчики блокировок комнат (лучше переименовать в block.go)
│   │   │   ├── health.go            # healthcheck-обработчик
│   │   │   ├── me.go                # обработчик получения текущего пользователя
│   │   │   ├── reservations.go      # обработчики бронирований: создание, просмотр, отмена
│   │   │   ├── response.go          # вспомогательные функции для единых HTTP-ответов и ошибок
│   │   │   └── rooms.go             # обработчики комнат: список, получение, создание, обновление, удаление
|   |   |
│   │   └── middleware/
│   │       ├── auth.go              # middleware для проверки JWT и извлечения пользователя
│   │       ├── id.go                # middleware для request id
│   │       ├── logging.go           # middleware для логирования запросов
│   │       └── role.go              # middleware для проверки роли пользователя, например admin
|   |
│   ├── repository/
│   │   └── postgres/
│   │       ├── block.go             # SQL-логика для блокировок комнат
│   │       ├── db.go                # подключение к PostgreSQL, инициализация пула соединений
│   │       ├── migrate.go           # запуск миграций базы данных
│   │       ├── reservations.go      # SQL-логика для бронирований
│   │       ├── rooms.go             # SQL-логика для комнат
│   │       └── users.go             # SQL-логика для пользователей
|   |
│   └── service/
│       ├── auth.go                  # бизнес-логика авторизации и регистрации
│       ├── availability.go          # бизнес-логика проверки доступности комнаты
│       ├── block.go                 # бизнес-логика блокировок
│       ├── reservations.go          # бизнес-логика бронирований
│       └── rooms.go                 # бизнес-логика работы с комнатами
|
├── migrations/
│   ├── 000001_init.up.sql           # создание основных таблиц и начальной схемы БД
│   ├── 000001_init.down.sql         # откат начальной схемы БД
│   ├── 000002_exclude_overlap.up.sql# запрет пересечения бронирований одной комнаты
│   ├── 000002_exclude_overlap.down.sql
│   ├── 000003_blocks_exclude.up.sql # запрет пересечения блокировок одной комнаты
│   └── 000003_blocks_exclude.down.sql
|
├── tests/
│   ├── integrations/
│   │   ├── main_test.go             # общая настройка интеграционных тестов
│   │   ├── reservation_conflict_test.go # тест на конфликт пересекающихся бронирований
│   │   └── reservations_race_test.go    # тест на конкурентные запросы и защиту от гонок
│   |
│   └── testutil/
│       ├── db.go                    # вспомогательные функции для тестовой БД
│       ├── env.go                   # подготовка окружения для тестов
│       ├── http.go                  # вспомогательные HTTP-функции для тестов
│       └── migrate.go               # применение миграций в тестовой среде
|
├── .env.example                     # пример переменных окружения
├── docker-compose.yml               
├── Dockerfile                       
├── go.mod                           # зависимости Go-модуля
└── README.md                       
```

## Как проходит запрос по коду 

### Пример 1. Бронирование

Маршрут: `POST /reservations`

Путь запроса:

1. `internal/httpx/router.go`
   - маршрут подключается в protected-группу;
   - перед обработчиком срабатывает JWT middleware.

2. `internal/httpx/middleware/auth.go`
   - берёт заголовок `Authorization: Bearer <token>`;
   - валидирует токен;
   - кладёт `user_id` и `role` в `gin.Context`.

3. `internal/httpx/handlers/reservations.go`
   - читает JSON тела запроса;
   - валидирует базовые поля;
   - достаёт `user_id` из middleware;
   - вызывает сервис `ReservationService.Create(...)`.

4. `internal/service/reservations.go`
   - проверяет корректность интервала;
   - переводит время в UTC;
   - проверяет ограничения по длительности;
   - не даёт бронировать в прошлом;
   - получает комнату через `rooms.GetByID(...)`;
   - запрещает бронь неактивной комнаты;
   - собирает доменную сущность `domain.Reservation`.

5. `internal/repository/postgres/reservations.go`
   - делает `INSERT INTO reservations ...`;
   - если Postgres ловит пересечение по `EXCLUDE`, репозиторий возвращает `errs.ErrConflict`.

6. `handlers/reservations.go`
   - превращает доменную сущность в `dto.ReservationResponse`;
   - отдаёт `201 Created`.

То есть слой HTTP не знает SQL, а слой репозитория не знает ничего про Gin.

### Пример 2. Логин

- `handlers/auth.go` принимает email и пароль;
- `service/auth.go` получает пользователя из БД;
- `internal/auth/password.go` сверяет пароль через bcrypt;
- `internal/auth/jwt.go` подписывает JWT;
- клиент получает `access_token`.


## База данных и миграции

Схема базы создаётся миграциями. Две активные брони одной комнаты не могут пересекаться по времени. Защита от гонок лежит не только в коде, но и в самой БД.

### Основные таблицы

- `users`
- `rooms`
- `reservations`
- `room_blocks`

### Что важно в миграциях

#### `000001_init.up.sql`
Создаёт таблицы и базовые ограничения.

#### `000002_exclude_overlap.up.sql`
Добавляет `EXCLUDE` constraint для таблицы `reservations`.

#### `000003_blocks_exclude.up.sql`
Добавляет такой же запрет пересечения для `room_blocks`.

## Переменные окружения

Основные переменные:

- `HTTP_ADDR` — адрес HTTP-сервера, по умолчанию `:8080`
- `DB_DSN` — DSN для подключения к Postgres
- `JWT_SECRET` — секрет для подписи токенов
- `LOG_LEVEL` — уровень логирования

Для Docker Compose ещё используются:

- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_PORT`
- `API_PORT`
- `DB_DSN_DOCKER`

Пример уже есть в `.env.example`.

## Запуск локально

### Что нужно

- Go
- PostgreSQL
- Docker, если прогонять интеграционные тесты

Сервис в `cmd/api/main.go` сам запускает миграции через `postgres.RunMigrations(...)`.

**приложение нужно запускать из корня репозитория**, иначе миграции могут не найтись.

### Шаги

1. Скопировать переменные окружения:

```bash
cp .env.example .env
```

На Windows можно просто создать `.env` вручную по образцу.

2. Поднять Postgres любым способом.

Например, если база уже установлена локально, то достаточно создать БД `booking` и проверить `DB_DSN`.

3. Выставить переменные окружения.

Пример для PowerShell:

```powershell
$env:HTTP_ADDR=':8080'
$env:DB_DSN='postgres://postgres:postgres@localhost:5432/booking?sslmode=disable'
$env:JWT_SECRET='dev_secret_change_me'
$env:LOG_LEVEL='info'
```

4. Запустить приложение **из корня проекта**:

```bash
go run ./cmd/api
```

После запуска будут автоматически применены миграции, откроется HTTP-сервер и можно будет работать с API.

### Что проверить после старта

- healthcheck: `GET /health`
- swagger: `GET /swagger/index.html`

- `http://localhost:8080/health`
- `http://localhost:8080/swagger/index.html`

## Запуск через Docker Compose

1. Создать `.env` по образцу `.env.example`.
2. Из корня проекта выполнить:

```bash
docker compose up --build
```

Что произойдёт:

- поднимется контейнер с PostgreSQL;
- соберётся и запустится API;
- API получит `DB_DSN_DOCKER`, где хост базы — `postgres`;
- после старта приложение само применит миграции.

Остановка:

```bash
docker compose down
```

## Основные маршруты

### Публичные

- `GET /health`
- `GET /swagger/*any`
- `POST /auth/register`
- `POST /auth/login`
- `GET /rooms`
- `GET /rooms/:id`
- `GET /rooms/:id/availability`

### Для любого авторизованного пользователя

- `GET /me`
- `POST /reservations`
- `GET /reservations/my`
- `GET /reservations/:id`
- `POST /reservations/:id/cancel`

### Только для администратора

- `POST /admin/rooms`
- `PUT /admin/rooms/:id`
- `DELETE /admin/rooms/:id`
- `POST /admin/rooms/:id/blocks`
- `DELETE /admin/blocks/:block_id`

## Минимальный сценарий работы с API


## Тесты

Интеграционные тесты лежат в `tests/integrations`.

Они проверяют:

- `reservation_conflict_test.go` — вторая пересекающаяся бронь получает `409 Conflict`;
- `reservations_race_test.go` — при конкурентных запросах одну и ту же бронь сможет создать только один запрос.

- тесты проверяют не только handler или service отдельно;
- поднимается тестовая среда с Postgres через Testcontainers;

Запуск:

```bash
go test ./tests/integrations -v
```

Для этого нужен Docker.


## Swagger

В проекте уже есть Swagger-документация в папке `docs/`.

Если нужно пересобрать swagger-файлы после изменения комментариев над handler-методами:

```bash
swag init -g cmd/api/main.go -o docs
```

После запуска приложения Swagger доступен по адресу:

```text
/http://localhost:8080/swagger/index.html
```

