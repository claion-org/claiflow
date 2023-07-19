package excute

import (
	"context"
	"database/sql"

	"github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt"
)

var ErrorCompose = stmt.ErrorCompose
var CauseIter = stmt.CauseIter

func ExecContext(ctx context.Context, tx Preparer, query string, args []interface{}) (affected int64, lastid int64, err error) {
	VANILLA_DEBUG_PRINT(query, args)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return 0, 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, 0, err
	}

	affected, err = result.RowsAffected()
	if err != nil {
		return 0, 0, err
	}

	lastid, err = result.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	return affected, lastid, nil
}

func QueryRowContext(ctx context.Context, tx Preparer, query string, args []interface{}) func(CallbackScanner) error {
	VANILLA_DEBUG_PRINT(query, args)

	return func(scan CallbackScanner) error {
		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return err
		}

		defer stmt.Close()

		row := stmt.QueryRowContext(ctx, args...)
		if err := row.Err(); err != nil {
			return err
		}

		return scan(row)
	}
}

func QueryRowsContext(ctx context.Context, tx Preparer, query string, args []interface{}) func(CallbackScanner) error {
	VANILLA_DEBUG_PRINT(query, args)

	return func(scan CallbackScanner) error {
		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		var rows *sql.Rows
		rows, err = stmt.QueryContext(ctx, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		if err := rows.Err(); err != nil {
			return err
		}

		for rows.Next() {
			if err := scan(rows); err != nil {
				return err
			}
		}

		return nil
	}
}

func Repeat(n int, s string) []string {
	ss := make([]string, n)
	for i := 0; i < n; i++ {
		ss[i] = s
	}
	return ss
}

func Count(dialect SqlExcutor, tableName string, cond stmt.ConditionStmt) func(ctx context.Context, tx Preparer) (count int, err error) {
	var (
		columns = []string{"COUNT(1)"}
	)

	return func(ctx context.Context, tx Preparer) (int, error) {
		var count int
		err := dialect.QueryRow(tableName, columns, cond)(
			ctx, tx)(
			func(scan Scanner) error {
				return scan.Scan(&count)
			})

		return count, err
	}
}

func Exist(dialect SqlExcutor, tableName string, cond stmt.ConditionStmt) func(ctx context.Context, tx Preparer) (exist bool, err error) {
	var (
		columns = []string{"COUNT(1)"}
	)

	return func(ctx context.Context, tx Preparer) (bool, error) {
		var count int
		err := dialect.QueryRow(tableName, columns, cond)(
			ctx, tx)(
			func(scan Scanner) error {
				return scan.Scan(&count)
			})

		return 0 < count, err
	}
}
