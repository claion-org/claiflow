package control

import (
	"context"

	"github.com/claion-org/claiflow/pkg/server/database/vanilla/excute"
	"github.com/claion-org/claiflow/pkg/server/database/vanilla/stmt"
	"github.com/claion-org/claiflow/pkg/server/model"
)

func GetTemplate(ctx context.Context, uuid string) (*model.Template, error) {
	var template model.Template
	template.UUID = uuid

	cond := stmt.And(
		stmt.Equal(model.TemplateFieldsUuid.String(), template.UUID),
	)

	err := Driver().QueryRow(template.TableName(), template.ColumnNames(), cond)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			return template.Scan(scanner)
		})

	return &template, err
}
