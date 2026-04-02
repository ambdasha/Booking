package integrations

import (
	"context"
	"os"
	"testing"

	repopg "booking/internal/repository/postgres"
	"booking/tests/testutil"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	env *testutil.Env
	db  *pgxpool.Pool
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	//создание среды
	e, err := testutil.NewEnv(ctx)
	if err != nil {
		panic(err) //если тестовая среда не поднялась, дальше тесты бессмысленны
	}
	env = e
	defer env.Terminate(ctx)

	// отдельный pool только для ResetDB/PromoteAdmin
	dbpool, err := repopg.NewPool(ctx, env.DSN)
	if err != nil {
		panic(err)
	}
	db = dbpool
	defer db.Close()

	//запуск тестов
	code := m.Run()
	os.Exit(code)
}