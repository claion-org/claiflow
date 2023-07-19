package globvar

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/claion-org/claiflow/pkg/logger"
	"github.com/claion-org/claiflow/pkg/server/database/vanilla/excute"
	"github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt"
	"github.com/claion-org/claiflow/pkg/server/macro/generic"
	"github.com/claion-org/claiflow/pkg/server/model"
	"github.com/pkg/errors"
)

const (
	SLICE_CAPACITY = 10
)

type GlobalVariantUpdate struct {
	*sql.DB
	dialect excute.SqlExcutor
	offset  time.Time //updated column
}

func NewGlobalVariablesUpdate(db *sql.DB, dialect excute.SqlExcutor) *GlobalVariantUpdate {
	return &GlobalVariantUpdate{
		DB:      db,
		dialect: dialect,
	}
}

// func (worker *GlobalVariantUpdate) Dialect() string {
// 	return worker.dialect
// }

// Update
//
//	Update = read db -> global_variables
func (worker *GlobalVariantUpdate) Update() error {
	records := make([]model.GlobalVariable, 0, SLICE_CAPACITY)
	globvar := model.GlobalVariable{}
	globvar.Updated = worker.offset

	globvar_cond := stmt.GT("updated", globvar.Updated)

	err := worker.dialect.QueryRows(globvar.TableName(), globvar.ColumnNames(), globvar_cond, nil, nil)(
		context.Background(), worker)(
		func(scan excute.Scanner) error {
			if err := globvar.Scan(scan); err != nil {
				return errors.WithStack(err)
			}

			records = generic.Append(records, globvar)

			return nil
		})
	if err != nil {
		return err
	}
	for _, record := range records {
		for _, gv := range GlobVars {
			if record.UUID == gv.UUID() {
				if err := gv.SetValue(record.Value); err != nil {
					logger.Logger().Error(err, "store global_variables",
						"uuid", record.UUID,
						"key", record.Name,
						"value", record.Value,
					)

					return err
				}
			}
		}
	}

	//update offset
	worker.offset = time.Now()

	return nil
}

// WhiteListCheck
// 리스트 체크
func (worker *GlobalVariantUpdate) WhiteListCheck() (err error) {
	records := make([]model.GlobalVariable, 0, SLICE_CAPACITY)

	globvar := model.GlobalVariable{}
	globvar.Updated = worker.offset

	err = worker.dialect.QueryRows(globvar.TableName(), globvar.ColumnNames(), nil, nil, nil)(
		context.Background(), worker)(
		func(scan excute.Scanner) error {
			if err := globvar.Scan(scan); err != nil {
				return errors.WithStack(err)
			}

			records = generic.Append(records, globvar)

			return nil
		})
	if err != nil {
		return
	}

	missingNames := make([]string, 0, len(GlobVars))
	count := 0

	for _, gv := range GlobVars {
		var found bool = false
		for _, record := range records {
			found = found || record.UUID == gv.UUID()
		}
		if !found {
			count++
			missingNames = append(missingNames, gv.Name())
		}

	}

	// push, pop := macro.StringBuilder()
	// for _, key := range KeyNames() {
	// 	var found bool = false
	// LOOP:
	// 	for i := range records {
	// 		if key == records[i].Name {
	// 			found = true
	// 			break LOOP
	// 		}
	// 	}
	// 	if !found {
	// 		count++
	// 		push(key)
	// 	}
	// }

	missingNames = generic.Map(missingNames, func(s string) string { return strconv.Quote(s) })

	if 0 < count {
		return errors.Errorf("not exists global_variables keys=%s", strings.Join(missingNames, ","))
	}

	return nil
}

func (worker *GlobalVariantUpdate) Merge() (err error) {
	records := make([]model.GlobalVariable, 0, SLICE_CAPACITY)

	globvar := model.GlobalVariable{}
	globvar.Updated = worker.offset

	err = worker.dialect.QueryRows(globvar.TableName(), globvar.ColumnNames(), nil, nil, nil)(
		context.Background(), worker)(
		func(scan excute.Scanner) error {
			if err := globvar.Scan(scan); err != nil {
				return errors.WithStack(err)
			}

			records = generic.Append(records, globvar)

			return nil
		})
	if err != nil {
		return
	}

	for _, gv := range GlobVars {
		var found bool = false
		for _, record := range records {
			found = found || record.UUID == gv.UUID()
		}

		if !found {
			updateColumns := []string{
				"summary",
				"value",
				"updated",
			}

			globvar := ToModel(gv, time.Now())

			_, _, err = worker.dialect.Upsert(globvar.TableName(), globvar.ColumnNames(), updateColumns, globvar.ColumnValues())(
				context.Background(), worker)
			if err != nil {
				return errors.Wrapf(err, "failed to create or update global_variables")
			}
		}
	}

	return nil
}

// func foreach_environment(elems []envv1.Environment, fn func(elem envv1.Environment) bool) {
// 	for n := range elems {
// 		ok := fn(elems[n])
// 		if !ok {
// 			return
// 		}
// 	}
// }
