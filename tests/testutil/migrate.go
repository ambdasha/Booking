package testutil
//чтобы тестовая БД получала ту же структуру, что и обычная база приложения.
import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func RunMigrations(dsn string) error {
	root, err := findRepoRoot()
	if err != nil {
		return err
	}

	//строится путь к папке migrations, а потом превращает его в URL вида
	migrationsPath := filepath.ToSlash(filepath.Join(root, "migrations"))
	sourceURL := "file://" + migrationsPath

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("sql open: %w", err)
	}
	defer db.Close()

	driver, err := migratepg.WithInstance(db, &migratepg.Config{})
	if err != nil {
		return fmt.Errorf("migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

// Ищем корень репозитория по go.mod, чтобы тесты работали из любой папки
func findRepoRoot() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("runtime.Caller failed")
	}
	dir := filepath.Dir(file)

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("go.mod not found вверх по дереву")
		}
		dir = parent
	}
}