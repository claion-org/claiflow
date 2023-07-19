package control

import (
	"context"
	"database/sql"
	"math"
	"sort"
	"time"

	"github.com/claion-org/claiflow/pkg/server/database/vanilla/excute"
	"github.com/claion-org/claiflow/pkg/server/database/vanilla/sqlex"
	"github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt"
	"github.com/claion-org/claiflow/pkg/server/macro/generic"
	"github.com/claion-org/claiflow/pkg/server/model"
)

func CreateService(ctx context.Context, newServicies []model.ServiceWithStatuses) error {
	err := sqlex.ScopeTx(ctx, Database(), func(tx *sql.Tx) error {
		for i := range newServicies {
			var err error

			service := newServicies[i].ClusterService
			statusies := newServicies[i].ClusterServiceStatuses

			_, _, err = Driver().Insert(service.TableName(), service.ColumnNames(), service.ColumnValues())(
				ctx, tx)
			if err != nil {
				return err
			}

			vv := generic.Map(statusies, func(a model.ClusterServiceStatus) []any { return a.ColumnValues() })
			var status model.ClusterServiceStatus
			_, _, err = Driver().Insert(status.TableName(), status.ColumnNames(), vv...)(
				ctx, tx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func FindClusterService(ctx context.Context, query stmt.ConditionStmt, order stmt.OrderStmt, page stmt.PaginationStmt) ([]model.ServiceWithStatuses, error) {
	// default pagination
	if page == nil {
		page = stmt.Limit(20)
	}

	services := make([]model.ClusterService, 0, INIT_SLICE_CAP)

	var service model.ClusterService

	err := Driver().QueryRows(service.TableName(), service.ColumnNames(), query, order, page)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			if err := service.Scan(scanner); err != nil {
				return err
			}

			services = generic.Append(services, service)

			return nil
		})

	out := make([]model.ServiceWithStatuses, 0, len(services))
	for i := range services {
		status, err := GetClusterServiceStatuses(ctx, services[i].ClusterUUID, services[i].UUID)
		if err != nil {
			return nil, err
		}

		out = generic.Append(out, model.ServiceWithStatuses{
			ClusterService:         services[i],
			ClusterServiceStatuses: status,
		})
	}

	return out, err
}

func GetClusterService(ctx context.Context, clusterUUID, uuid string) (*model.ClusterService, error) {
	var service model.ClusterService
	service.ClusterUUID = clusterUUID
	service.UUID = uuid

	cond := stmt.And(
		stmt.Equal(model.ClusterServiceFieldsClusterUuid.String(), service.ClusterUUID),
		stmt.Equal(model.ClusterServiceFieldsUuid.String(), service.UUID),
	)

	err := Driver().QueryRow(service.TableName(), service.ColumnNames(), cond)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			return service.Scan(scanner)
		})

	return &service, err
}

func GetClusterServiceResult(ctx context.Context, clusterUUID, uuid string) (*model.ClusterServiceResult, error) {
	var result model.ClusterServiceResult
	result.ClusterUUID = clusterUUID
	result.UUID = uuid

	cond := stmt.And(
		stmt.Equal(model.ClusterServiceFieldsClusterUuid.String(), result.ClusterUUID),
		stmt.Equal(model.ClusterServiceResultFieldsUuid.String(), result.UUID),
	)

	err := Driver().QueryRow(result.TableName(), result.ColumnNames(), cond)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			return result.Scan(scanner)
		})

	return &result, err
}

func GetClusterServiceResults(ctx context.Context, clusterUUID, uuid string) ([]model.ClusterServiceResult, error) {
	var out = make([]model.ClusterServiceResult, 0, INIT_SLICE_CAP)
	var result model.ClusterServiceResult
	result.ClusterUUID = clusterUUID
	result.UUID = uuid

	cond := stmt.And(
		stmt.Equal(model.ClusterServiceFieldsClusterUuid.String(), result.ClusterUUID),
		stmt.Equal(model.ClusterServiceResultFieldsUuid.String(), result.UUID),
	)

	err := Driver().QueryRows(result.TableName(), result.ColumnNames(), cond, nil, nil)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			if err := result.Scan(scanner); err != nil {
				return err
			}

			out = generic.Append(out, result)

			return nil
		})

	return out, err
}

func GetClusterServiceStatuses(ctx context.Context, clusterUUID, uuid string) ([]model.ClusterServiceStatus, error) {
	var out = make([]model.ClusterServiceStatus, 0, INIT_SLICE_CAP)
	var status model.ClusterServiceStatus
	status.ClusterUUID = clusterUUID
	status.Uuid = uuid

	cond := stmt.And(
		stmt.Equal(model.ClusterServiceFieldsClusterUuid.String(), status.ClusterUUID),
		stmt.Equal(model.ClusterServiceStatusFieldsUuid.String(), status.Uuid),
	)

	err := Driver().QueryRows(status.TableName(), status.ColumnNames(), cond, nil, nil)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			if err := status.Scan(scanner); err != nil {
				return err
			}

			out = generic.Append(out, status)

			return nil
		})

	// sort
	sort.Slice(out, func(i, j int) bool { return out[i].Created.Before(out[j].Created) })

	return out, err
}

func PollClusterService(ctx context.Context, clusterUUID string, offset time.Time, filter func(*model.ClusterService) bool) ([]model.ClusterService, error) {
	out := make([]model.ClusterService, 0, INIT_SLICE_CAP)

	var service model.ClusterService
	service.ClusterUUID = clusterUUID
	service.Created = offset

	cond := stmt.And(
		stmt.Equal(model.ClusterServiceFieldsClusterUuid.String(), service.ClusterUUID),
		stmt.GT(model.ClusterServiceFieldsCreated.String(), service.Created),
	)

	order := stmt.Asc(model.ClusterServiceFieldsCreated.String())

	page := stmt.Limit(math.MaxInt8)

	err := Driver().QueryRows(service.TableName(), service.ColumnNames(), cond, order, page)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			var service model.ClusterService
			if err := service.Scan(scanner); err != nil {
				return err
			}

			if filter(&service) {
				out = generic.Append(out, service)
			}

			return nil
		})

	return out, err
}

func CreateClusterServiceStatus(ctx context.Context, status *model.ClusterServiceStatus) error {
	_, _, err := Driver().Insert(status.TableName(), status.ColumnNames(), status.ColumnValues())(
		ctx, Database())

	return err
}

func CreateClusterServiceStatuses(ctx context.Context, statuses []model.ClusterServiceStatus) error {
	vv := generic.Map(statuses, func(a model.ClusterServiceStatus) []any {
		return a.ColumnValues()
	})

	_, _, err := Driver().Insert(model.ClusterServiceStatus{}.TableName(), model.ClusterServiceStatus{}.ColumnNames(), vv...)(
		ctx, Database())

	return err
}

func UpsertClusterServiceResult(ctx context.Context, result *model.ClusterServiceResult) error {

	updateColumns := []string{
		model.ClusterServiceResultFieldsCreated.String(),
		model.ClusterServiceResultFieldsResultType.String(),
		model.ClusterServiceResultFieldsResult.String(),
	}

	_, _, err := Driver().Upsert(result.TableName(), result.ColumnNames(), updateColumns, result.ColumnValues())(
		ctx, Database())

	return err
}
