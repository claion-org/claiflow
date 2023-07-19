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

func CreateCluster(ctx context.Context, model *model.Cluster) error {
	affected, id, err := Driver().Insert(model.TableName(), model.ColumnNames(), model.ColumnValues())(
		ctx, Database())
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrNoAffected
	}

	model.ID = id

	return err
}

func FindCluster(ctx context.Context, query stmt.ConditionStmt, order stmt.OrderStmt, page stmt.PaginationStmt) ([]model.Cluster, error) {
	out := make([]model.Cluster, 0, INIT_SLICE_CAP)

	var cluster model.Cluster

	err := Driver().QueryRows(cluster.TableName(), cluster.ColumnNames(), query, order, page)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			if err := cluster.Scan(scanner); err != nil {
				return err
			}

			out = generic.Append(out, cluster)

			return nil
		})

	return out, err
}

func GetCluster(ctx context.Context, uuid string) (*model.Cluster, error) {
	var cluster model.Cluster
	cluster.UUID = uuid

	cond := stmt.And(
		stmt.Equal(model.ClusterFieldsUuid.String(), cluster.UUID),
	)

	err := Driver().QueryRow(cluster.TableName(), cluster.ColumnNames(), cond)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			return cluster.Scan(scanner)
		})

	return &cluster, err
}

func UpsertCluster(ctx context.Context, cluster *model.Cluster, updateColumns []string) error {
	_, lastID, err := Driver().Upsert(cluster.TableName(), cluster.ColumnNames(), updateColumns, cluster.ColumnValues())(
		ctx, Database())

	cluster.ID = lastID

	return err
}

func DeleteCluster(ctx context.Context, uuid string) error {
	// Cluster
	var cluster model.Cluster
	cluster.UUID = uuid

	clusterCond := stmt.And(stmt.Equal(
		model.ClusterFieldsUuid.String(), cluster.UUID,
	))

	// ClusterClientToken
	var clientToken model.ClusterClientToken
	clientToken.ClusterUUID = uuid

	clientToeknCond := stmt.And(stmt.Equal(
		model.ClusterClientTokenFieldsClusterUuid.String(), clientToken.ClusterUUID,
	))

	// ClusterClientSession
	var clientSession model.ClusterClientSession
	clientSession.ClusterUUID = uuid

	clientSessionCond := stmt.And(stmt.Equal(
		model.ClusterClientSessionFieldsClusterUuid.String(), clientSession.ClusterUUID,
	))

	err := sqlex.ScopeTx(ctx, Database(), func(tx *sql.Tx) error {
		var err error

		// ClusterClientSession
		_, err = Driver().Delete(clientSession.TableName(), clientSessionCond)(
			ctx, tx)
		if err != nil {
			return err
		}

		// ClusterClientToken
		_, err = Driver().Delete(clientToken.TableName(), clientToeknCond)(
			ctx, tx)
		if err != nil {
			return err
		}

		// Cluster
		_, err = Driver().Delete(cluster.TableName(), clusterCond)(
			ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func IsExistsCluster(ctx context.Context, uuid string) (bool, error) {
	var cluster model.Cluster
	cluster.UUID = uuid

	cond := stmt.And(
		stmt.Equal(model.ClusterFieldsUuid.String(), cluster.UUID),
	)

	exists, err := Driver().Exist(cluster.TableName(), cond)(
		ctx, Database())

	return exists, err
}
