package testutil
//поднимает полное тестовое окружение
import (
	"context"
	"log/slog"
	"net/http/httptest"
	"os"
	"time"

	"booking/internal/config"
	"booking/internal/httpx"
	repopg "booking/internal/repository/postgres"

	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type Env struct {
	BaseURL string
	DSN string

	DBClose func()  //функция закрытия пула соединений с БД
	Close   func()  //функция остановки HTTP-сервера

	pg *tcpostgres.PostgresContainer
}


//функция создания всей среды
func NewEnv(ctx context.Context) (*Env, error) {
	//поднятие  Postgres контейнера
	ctxStart, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	pg, err := tcpostgres.Run(
		ctxStart,
		"postgres:16",
		tcpostgres.WithDatabase("booking"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, err
	}

	//получение DSN 
	dsn, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = pg.Terminate(ctx)
		return nil, err
	}

	//прогнать миграции по только что созданной пустой тестовой БД
	if err := RunMigrations(dsn); err != nil {
		_ = pg.Terminate(ctx)
		return nil, err
	}

	//поднять pgxpool тем же кодом, что в приложении
	dbpool, err := repopg.NewPool(ctx, dsn)
	if err != nil {
		_ = pg.Terminate(ctx)
		return nil, err
	}

	//запустить роутер в памяти
	var cfg config.Config
	cfg.HTTP.Addr = ":0"
	cfg.DB.DSN = dsn
	cfg.Auth.JWTSecret = "test_secret"
	cfg.Log.Level = "error"

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	router := httpx.NewRouter(cfg, dbpool, logger)
	ts := httptest.NewServer(router)

	env := &Env{
		BaseURL: ts.URL,
		DSN:     dsn,
		DBClose: dbpool.Close,
		Close:   ts.Close,
		pg:      pg,
	}
	return env, nil
}

//выключает всё, что было запущено
func (e *Env) Terminate(ctx context.Context) {
	if e.Close != nil {
		e.Close()
	}
	if e.DBClose != nil {
		e.DBClose()
	}
	if e.pg != nil {
		_ = e.pg.Terminate(ctx)
	}
}