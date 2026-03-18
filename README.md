


booking-service/
├── cmd/
│   └── api/
│       └── main.go                 # вход: config → db → deps → router → run
│
├── internal/
│   ├── app/
│   │   └── app.go                  # сборка зависимостей (wire-up), старт/stop
│   │
│   ├── config/
│   │   ├── config.go               # структура конфига
│   │   └── env.go                  # чтение env + дефолты
│   │
│   ├── domain/
│   │   ├── user.go
│   │   ├── room.go
│   │   ├── reservation.go
│   │   └── room_block.go           # доменные сущности (не DTO)
│   │
│   ├── dto/
│   │   ├── auth.go                 # Register/Login request/response
│   │   ├── room.go                 # Create/Update room DTO
│   │   ├── reservation.go          # Create/Cancel reservation DTO
│   │   ├── availability.go
│   │   └── common.go               # pagination, filters, error response
│   │
│   ├── http/
│   │   ├── router.go               # все маршруты + группы (/admin)
│   │   ├── handlers/
│   │   │   ├── auth.go
│   │   │   ├── me.go
│   │   │   ├── rooms.go
│   │   │   ├── reservations.go
│   │   │   ├── availability.go
│   │   │   └── blocks.go
│   │   ├── middleware/
│   │   │   ├── auth.go             # JWT auth
│   │   │   ├── role.go             # admin-only
│   │   │   ├── request_id.go
│   │   │   ├── logging.go
│   │   │   └── recover.go
│   │   └── response.go             # единый JSON-формат ошибок/ответов
│   │
│   ├── service/
│   │   ├── auth.go
│   │   ├── rooms.go
│   │   ├── reservations.go         # бизнес-логика + правила
│   │   ├── availability.go
│   │   └── blocks.go
│   │
│   ├── repository/
│   │   ├── postgres/
│   │   │   ├── db.go               # подключение, ping, pool
│   │   │   ├── tx.go               # helper для транзакций (если надо)
│   │   │   ├── users.go
│   │   │   ├── rooms.go
│   │   │   ├── reservations.go     # INSERT ловит EXCLUDE → domain ErrConflict
│   │   │   └── blocks.go
│   │   └── errors.go               # маппинг ошибок БД → доменные
│   │
│   ├── auth/
│   │   ├── jwt.go                  # генерация/проверка токенов
│   │   ├── password.go             # bcrypt hash/compare
│   │   └── context.go              # положить/достать user из ctx
│   │
│   ├── validation/
│   │   ├── auth.go
│   │   ├── room.go
│   │   └── reservation.go          # доп. валидация интервалов
│   │
│   ├── errs/
│   │   ├── domain.go               # ErrNotFound/ErrConflict/ErrForbidden...
│   │   └── http.go                 # маппинг domain errors → HTTP codes
│   │
│   └── observability/
│       └── logger.go               # slog setup (можно и без этого пакета)
│
├── migrations/
│   ├── 000001_init.up.sql
│   ├── 000001_init.down.sql
│   ├── 000002_exclude_overlap.up.sql   # тут EXCLUDE + btree_gist
│   └── 000002_exclude_overlap.down.sql
│
├── deployments/
│   └── docker/
│       ├── Dockerfile
│       └── docker-compose.yml
│
├── tests/
│   ├── integration/
│   │   ├── auth_test.go
│   │   ├── rooms_test.go
│   │   ├── reservations_test.go
│   │   └── conflicts_test.go       # параллельные брони → одна проходит
│   └── testutil/
│       ├── http_client.go
│       ├── fixtures.go
│       └── db.go
│
├── docs/
│   └── openapi.yaml                # если будешь делать OpenAPI (опционально)
│
├── scripts/
│   ├── migrate-up.sh
│   ├── migrate-down.sh
│   └── test.sh
│
├── .env.example
├── Makefile
├── go.mod
├── go.sum
└── README.md