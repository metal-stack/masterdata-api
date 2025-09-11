package datastore

import (
	"context"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func RunQuery[E any](ctx context.Context, log *slog.Logger, db *sqlx.DB, builder squirrel.SelectBuilder, input map[string]any, resultFn func(e E) error) error {
	query, vals, err := builder.ToSql()
	if err != nil {
		return err
	}

	if log.Enabled(ctx, slog.LevelDebug) {
		log.Debug("query", "sql", query, "values", vals, "input", input)
	}

	rows, err := db.NamedQueryContext(ctx, query, input)
	if err != nil {
		return err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Error("error closing result rows", "error", err)
		}
	}()

	for rows.Next() {
		var e E

		err := rows.StructScan(&e)
		if err != nil {
			return err
		}

		err = resultFn(e)
		if err != nil {
			return err
		}
	}

	return nil
}
