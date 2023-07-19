package webhook

import (
	"bytes"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/claion-org/claiflow/pkg/server/macro/generic"
	"github.com/claion-org/claiflow/pkg/server/model"
	"github.com/itchyny/gojq"
)

type Config struct {
	URL                string
	Method             string
	Headers            http.Header
	ConditionValidator model.WebhookConditionValidator
	ConditionFilter    string
	Timeout            time.Duration
}

func (config Config) Publish(ctx context.Context, payload interface{}) error {
	pred, err := CheckCondition(config.ConditionValidator, config.ConditionFilter)(ctx, payload)
	if err != nil {
		return err
	}
	if !pred {
		// the condition does not match
		return nil
	}

	if 0 < config.Timeout {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, config.Timeout)
		defer cancel()
	}

	req, err := NewRequest(ctx, config.Method, config.URL,
		ToJson(payload),
		SetHeader(http.Header(config.Headers)),
		Compress_GZip,
	)
	if err != nil {
		return err
	}

	client := DefaultClient()

	var out bytes.Buffer
	if err := Do(client, req,
		CheckStatus,
		GetBody(&out)); err != nil {
		return err
	}

	return nil
}

func CheckConditionFilter(validator model.WebhookConditionValidator, filter string) error {
	switch model.WebhookConditionValidator((validator)) {
	case model.WebhookConditionValidatorJq:
		return generic.Right(gojq.Parse(filter))
	default:
		return nil
	}
}

func CheckCondition(validator model.WebhookConditionValidator, filter string) func(ctx context.Context, payload interface{}) (bool, error) {
	switch model.WebhookConditionValidator((validator)) {
	case model.WebhookConditionValidatorJq:
		return func(ctx context.Context, payload interface{}) (bool, error) {
			// conditional check by jq
			jq, err := gojq.Parse(filter)
			if err != nil {
				return false, err
			}

			var pred = make([]bool, 0, 1)
			iter := jq.RunWithContext(ctx, payload)
			for v, next := iter.Next(); next; v, next = iter.Next() {
				switch v := v.(type) {
				case bool:
					pred = generic.Append(pred, v)
				case string:
					pred = generic.Append(pred, generic.Left(strconv.ParseBool(v)))
				}
			}

			if len(pred) == 0 {
				return false, nil
			}

			ok := generic.Fold(true, pred, func(acc bool, a bool) bool {
				return acc && a
			})

			return ok, nil
		}
	case model.WebhookConditionValidatorNone:
		fallthrough
	default:
		return func(ctx context.Context, payload interface{}) (bool, error) {
			// always pass under any condition
			return true, nil
		}
	}
}
