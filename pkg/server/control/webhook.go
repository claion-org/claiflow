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

func CreateWebhook(ctx context.Context, webhook *model.Webhook) error {
	affected, id, err := Driver().Insert(webhook.TableName(), webhook.ColumnNames(), webhook.ColumnValues())(
		ctx, Database())
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrNoAffected
	}

	webhook.ID = id

	return err
}

func FindWebhook(ctx context.Context, query stmt.ConditionStmt, order stmt.OrderStmt, page stmt.PaginationStmt) ([]model.Webhook, error) {
	out := make([]model.Webhook, 0, INIT_SLICE_CAP)

	var webhook model.Webhook

	err := Driver().QueryRows(webhook.TableName(), webhook.ColumnNames(), query, order, page)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			if err := webhook.Scan(scanner); err != nil {
				return err
			}

			out = generic.Append(out, webhook)

			return nil
		})

	return out, err
}

func GetWebhook(ctx context.Context, uuid string) (*model.Webhook, error) {
	var webhook model.Webhook
	webhook.UUID = uuid

	cond := stmt.And(
		stmt.Equal(model.WebhookFieldsUuid.String(), webhook.UUID),
	)

	err := Driver().QueryRow(webhook.TableName(), webhook.ColumnNames(), cond)(
		ctx, Database())(
		func(scanner excute.Scanner) error {
			return webhook.Scan(scanner)
		})

	return &webhook, err
}

func UpsertWebhook(ctx context.Context, webhook *model.Webhook, updateColumns []string) error {
	_, lastID, err := Driver().Upsert(webhook.TableName(), webhook.ColumnNames(), updateColumns, webhook.ColumnValues())(
		ctx, Database())

	webhook.ID = lastID

	return err
}

func DeleteWebhook(ctx context.Context, uuid string) error {
	// Webhook
	var webhook model.Webhook
	webhook.UUID = uuid

	clusterCond := stmt.And(stmt.Equal(
		model.WebhookFieldsUuid.String(), webhook.UUID,
	))

	err := sqlex.ScopeTx(ctx, Database(), func(tx *sql.Tx) error {
		var err error

		// Webhook
		_, err = Driver().Delete(webhook.TableName(), clusterCond)(
			ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
