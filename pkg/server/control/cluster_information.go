package control

import (
	"context"
	"database/sql"
	"time"

	"github.com/claion-org/claiflow/pkg/server/database/vanilla/excute"
	"github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt"
	"github.com/claion-org/claiflow/pkg/server/model"
)

func GetClusterPollingOffset(ctx context.Context, clusterUUID string) (time.Time, error) {
	var clusterInfo model.ClusterInformation
	clusterInfo.ClusterUUID = clusterUUID

	cond := stmt.And(
		stmt.Equal(model.ClusterInformationFieldsClusterUuid.String(), clusterInfo.ClusterUUID),
	)

	columns := []string{
		model.ClusterInformationFieldsPollingOffset.String(),
	}

	var pollingOffest sql.NullTime

	err := Driver().QueryRow(clusterInfo.TableName(), columns, cond)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			return scanner.Scan(&pollingOffest)
		})

	return pollingOffest.Time, err
}

func UpsertClusterPollingOffset(ctx context.Context, clusterUUID string, pollingOffset time.Time, updated time.Time) error {
	var clusterInfo model.ClusterInformation
	clusterInfo.ClusterUUID = clusterUUID
	clusterInfo.PollingOffset = pollingOffset
	clusterInfo.Updated = updated

	updateColumns := []string{
		model.ClusterInformationFieldsPollingOffset.String(),
		model.ClusterInformationFieldsUpdated.String(),
	}

	affected, _, err := Driver().Upsert(clusterInfo.TableName(), clusterInfo.ColumnNames(), updateColumns, clusterInfo.ColumnValues())(
		ctx, Database())
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
