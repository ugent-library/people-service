package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivermigrate"
)

var migrator = rivermigrate.New(riverdatabasesql.New(nil), nil)

func init() {
	goose.AddMigrationContext(Up00002, Down00002)
}

func Up00002(ctx context.Context, tx *sql.Tx) error {
	_, err := migrator.MigrateTx(ctx, tx, rivermigrate.DirectionUp, &rivermigrate.MigrateOpts{
		TargetVersion: 3,
	})
	return err
}

func Down00002(ctx context.Context, tx *sql.Tx) error {
	_, err := migrator.MigrateTx(ctx, tx, rivermigrate.DirectionDown, &rivermigrate.MigrateOpts{
		TargetVersion: -1,
	})
	return err
}
