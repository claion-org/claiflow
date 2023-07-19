package control

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/claion-org/claiflow/pkg/server/database/vanilla/excute"
	"github.com/claion-org/claiflow/pkg/server/database/vanilla/sqlex"
	"github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt"
	"github.com/claion-org/claiflow/pkg/server/macro/generic"
	"github.com/claion-org/claiflow/pkg/server/model"
)

func CreateClusterClientToken(ctx context.Context, token *model.ClusterClientToken) error {
	// is cluster exists?
	exists, err := IsExistsCluster(ctx, token.ClusterUUID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%w: cluster was not exists uuid=%q", ErrNoRows, token.ClusterUUID)
	}

	// insert record
	affected, id, err := Driver().Insert(token.TableName(), token.ColumnNames(), token.ColumnValues())(
		ctx, Database())
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrNoAffected
	}

	token.ID = id

	return err
}

func FindClusterClientToken(ctx context.Context, query stmt.ConditionStmt, order stmt.OrderStmt, page stmt.PaginationStmt) ([]model.ClusterClientToken, error) {
	out := make([]model.ClusterClientToken, 0, INIT_SLICE_CAP)

	var token model.ClusterClientToken

	err := Driver().QueryRows(token.TableName(), token.ColumnNames(), query, order, page)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			if err := token.Scan(scanner); err != nil {
				return err
			}

			out = generic.Append(out, token)

			return nil
		})

	return out, err
}

func GetClusterClientToken(ctx context.Context, uuid string) (*model.ClusterClientToken, error) {
	var token model.ClusterClientToken
	token.UUID = uuid

	cond := stmt.And(
		stmt.Equal(model.ClusterClientTokenFieldsUuid.String(), token.UUID),
	)

	err := Driver().QueryRow(token.TableName(), token.ColumnNames(), cond)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			return token.Scan(scanner)
		})

	return &token, err
}

func GetClusterClientTokenByAssertion(ctx context.Context, clusterUUID string, token_ string) (*model.ClusterClientToken, error) {
	var token model.ClusterClientToken
	token.ClusterUUID = clusterUUID
	token.Token = token_

	cond := stmt.And(
		stmt.Equal(model.ClusterClientTokenFieldsClusterUuid.String(), token.ClusterUUID),
		stmt.Equal(model.ClusterClientTokenFieldsToken.String(), token.Token),
	)

	err := Driver().QueryRow(token.TableName(), token.ColumnNames(), cond)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			return token.Scan(scanner)
		})

	return &token, err
}

// GetClusterClientTokenColumns
func GetClusterClientTokenColumns(ctx context.Context, uuid string, columns map[string]interface{}) error {
	var token model.ClusterClientToken
	token.UUID = uuid

	cond := stmt.And(
		stmt.Equal(model.ClusterClientTokenFieldsUuid.String(), token.UUID),
	)

	columnNames := make([]string, 0, len(columns))
	columnPtrs := make([]interface{}, 0, len(columns))

	for k, v := range columns {
		columnNames = append(columnNames, k)
		columnPtrs = append(columnPtrs, v)
	}

	err := Driver().QueryRow(token.TableName(), columnNames, cond)(
		ctx, Database())(func(scanner excute.Scanner) error {
		return scanner.Scan(columnPtrs...)
	})

	return err
}

func UpsertClusterClientToken(ctx context.Context, token *model.ClusterClientToken, updateColumns []string) error {
	_, lastID, err := Driver().Upsert(token.TableName(), token.ColumnNames(), updateColumns, token.ColumnValues())(
		ctx, Database())

	token.ID = lastID

	return err
}

func DeleteClusterClientToken(ctx context.Context, uuid string) error {
	var token model.ClusterClientToken
	token.UUID = uuid

	tokenCond := stmt.And(stmt.Equal(
		model.ClusterClientTokenFieldsUuid.String(), token.UUID,
	))

	var session model.ClusterClientSession
	session.ClusterClientTokenUUID = uuid

	sessionCond := stmt.And(stmt.Equal(
		model.ClusterClientSessionFieldsClusterClientTokenUuid.String(), session.ClusterClientTokenUUID,
	))

	err := sqlex.ScopeTx(ctx, Database(), func(tx *sql.Tx) error {
		_, err := Driver().Delete(session.TableName(), sessionCond)(
			ctx, tx)
		if err != nil {
			return err
		}

		_, err = Driver().Delete(token.TableName(), tokenCond)(
			ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
