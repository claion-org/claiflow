package control

import (
	"context"
	"database/sql"

	"github.com/claion-org/claiflow/pkg/server/database/vanilla/excute"
	"github.com/claion-org/claiflow/pkg/server/database/vanilla/sqlex"
	"github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt"
	"github.com/claion-org/claiflow/pkg/server/macro/generic"
	"github.com/claion-org/claiflow/pkg/server/model"
)

func CreateClusterClientSession(ctx context.Context, session *model.ClusterClientSession) error {
	// insert record
	affected, id, err := Driver().Insert(session.TableName(), session.ColumnNames(), session.ColumnValues())(
		ctx, Database())
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrNoAffected
	}

	session.ID = id

	return err
}

func FindClusterClientSession(ctx context.Context, query stmt.ConditionStmt, order stmt.OrderStmt, page stmt.PaginationStmt) ([]model.ClusterClientSession, error) {
	out := make([]model.ClusterClientSession, 0, INIT_SLICE_CAP)

	var session model.ClusterClientSession

	err := Driver().QueryRows(session.TableName(), session.ColumnNames(), query, order, page)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			if err := session.Scan(scanner); err != nil {
				return err
			}

			out = generic.Append(out, session)

			return nil
		})

	return out, err
}

func GetClusterClientSession(ctx context.Context, uuid string) (*model.ClusterClientSession, error) {
	var session model.ClusterClientSession
	session.UUID = uuid

	cond := stmt.And(
		stmt.Equal(model.ClusterClientSessionFieldsUuid.String(), session.UUID),
	)

	err := Driver().QueryRow(session.TableName(), session.ColumnNames(), cond)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			return session.Scan(scanner)
		})

	return &session, err
}

func UpsertClusterClientSession(ctx context.Context, session *model.ClusterClientSession, updateColumns []string) error {
	_, id, err := Driver().Upsert(session.TableName(), session.ColumnNames(), updateColumns, session.ColumnValues())(
		ctx, Database())

	session.ID = id

	return err
}

func DeleteClusterClientSession(ctx context.Context, uuid string) error {
	var session model.ClusterClientSession
	session.UUID = uuid

	sessionCond := stmt.And(stmt.Equal(
		model.ClusterClientSessionFieldsUuid.String(), session.UUID,
	))

	err := sqlex.ScopeTx(ctx, Database(), func(tx *sql.Tx) error {
		_, err := Driver().Delete(session.TableName(), sessionCond)(
			ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func IsExistsClusterClientSession(ctx context.Context, uuid string) (bool, error) {
	var session model.ClusterClientSession
	session.UUID = uuid

	cond := stmt.And(
		stmt.Equal(model.ClusterClientSessionFieldsUuid.String(), session.UUID),
	)

	exists, err := Driver().Exist(session.TableName(), cond)(
		ctx, Database())

	return exists, err
}
