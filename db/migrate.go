package db

import (
	"context"
	"embed"
	"io/fs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func newMigrator(ctx context.Context, db *pgx.Conn) (*migrate.Migrator, error) {
	migrator, err := migrate.NewMigrator(ctx, db, "schema_version")
	if err != nil {
		return nil, err
	}
	migrations, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return nil, err
	}
	if err = migrator.LoadMigrations(migrations); err != nil {
		return nil, err
	}

	return migrator, nil
}

func Migrate(ctx context.Context, conn string) error {
	return MigrateTo(ctx, conn, -1)
}

func MigrateTo(ctx context.Context, conn string, version int32) error {
	db, err := pgx.Connect(ctx, conn)
	if err != nil {
		return err
	}
	defer db.Close(ctx)
	m, err := newMigrator(ctx, db)
	if err != nil {
		return err
	}
	if version >= 0 {
		return m.MigrateTo(ctx, version)
	}
	return m.Migrate(ctx)
}
