package sql

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	_ "github.com/jackc/pgx/v5/stdlib" // SQL driver

	"github.com/TBD54566975/ftl/backend/common/log"
)

//go:embed schema
var migrationSchema embed.FS

// Migrate the database.
func Migrate(ctx context.Context, dsn string) error {
	u, err := url.Parse(dsn)
	if err != nil {
		return fmt.Errorf("%s: %w", "invalid DSN", err)
	}
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("%s: %w", "failed to connect to database", err)
	}
	defer conn.Close()

	db := dbmate.New(u)
	db.FS = migrationSchema
	db.Log = log.FromContext(ctx).Scope("migrate").WriterAt(log.Debug)
	db.MigrationsDir = []string{"schema"}
	err = db.CreateAndMigrate()
	if err != nil {
		return fmt.Errorf("%s: %w", "failed to create and migrate database", err)
	}
	return nil
}
