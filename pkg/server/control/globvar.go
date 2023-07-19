package control

import (
	"context"

	"github.com/claion-org/claiflow/pkg/server/database/vanilla/excute"
	"github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt"
	"github.com/claion-org/claiflow/pkg/server/macro/generic"
	"github.com/claion-org/claiflow/pkg/server/model"
)

func FindGlobalVariables(ctx context.Context, query stmt.ConditionStmt, order stmt.OrderStmt) ([]model.GlobalVariable, error) {
	out := make([]model.GlobalVariable, 0, INIT_SLICE_CAP)

	var globvar model.GlobalVariable

	err := Driver().QueryRows(globvar.TableName(), globvar.ColumnNames(), query, order, nil)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			if err := globvar.Scan(scanner); err != nil {
				return err
			}

			out = generic.Append(out, globvar)

			return nil
		})

	return out, err
}

func GetGlobalVariable(ctx context.Context, uuid string) (*model.GlobalVariable, error) {
	var globvar model.GlobalVariable
	globvar.UUID = uuid

	cond := stmt.And(
		stmt.Equal(model.GlobalVariableFieldsUuid.String(), globvar.UUID),
	)

	err := Driver().QueryRow(globvar.TableName(), globvar.ColumnNames(), cond)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			return globvar.Scan(scanner)
		})

	return &globvar, err
}

func UpsertGlobalVariable(ctx context.Context, globvar *model.GlobalVariable, updateColumns []string) error {
	_, lastID, err := Driver().Upsert(globvar.TableName(), globvar.ColumnNames(), updateColumns, globvar.ColumnValues())(
		ctx, Database())

	globvar.ID = lastID

	return err
}

// func DeleteGlobalVariable(ctx context.Context, key string) error {
// 	globvar, ok := defaultGlobalVariables()[key]

// 	if !ok {
// 		return fmt.Errorf("invalid global variable key")
// 	}

// 	cond := stmt.And(
// 		stmt.Equal("key", globvar.Key),
// 	)

// 	_, err := Driver().Delete(globvar.TableName(), cond)(
// 		ctx, Database())

// 	return err
// }
