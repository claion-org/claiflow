package stmt

import "github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt"

const (
	__DIALECT__ = "mysql"
)

func Dialect() string {
	return __DIALECT__
}

func init() {
	stmt.SetConditionStmtBuilder(__DIALECT__, NewMysqlCondition())
	stmt.SetOrderStmtBuilder(__DIALECT__, NewMysqlOrder())
	stmt.SetPaginationStmtBuilder(__DIALECT__, NewMysqlPagination())
}
