package webhook

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/itchyny/gojq"
)

func TestJq_Map(t *testing.T) {
	data := `
	{
		"outcome": "UP",
			"checks": [
			{
				"name": "core-isAlive",
				"state": "UP",
				"data": {
					"buildtime": "2020-01-29T12:46:05.183Z",
					"version": "0.0.1-SNAPSHOT",
					"api-versions": "2"
				}
			},
			{
				"name": "core-db2-isAlive",
				"state": "UP",
				"data": {}
			},
			{
				"name": "core-postgres-isAlive",
				"state": "UP",
				"data": {}
			}]
	}`

	src := ".checks[] | .name"

	var v = map[string]any{}
	err := json.Unmarshal([]byte(data), &v)
	if err != nil {
		t.Fatal(err)
	}

	jq, err := gojq.Parse(src)
	if err != nil {
		t.Fatal(err)
	}

	iter := jq.Run(v)

	for v, ok := iter.Next(); ok; v, ok = iter.Next() {
		t.Log(v)
	}
}

func TestJq_Slice(t *testing.T) {
	data := `
	[{
		"outcome": "UP",
			"checks": [
			{
				"name": "core-isAlive",
				"state": "UP",
				"data": {
					"buildtime": "2020-01-29T12:46:05.183Z",
					"version": "0.0.1-SNAPSHOT",
					"api-versions": "2"
				}
			},
			{
				"name": "core-db2-isAlive",
				"state": "UP",
				"data": {}
			},
			{
				"name": "core-postgres-isAlive",
				"state": "UP",
				"data": {}
			}]
	}]`

	src := ".[].checks[] | .name"

	var v any
	err := json.Unmarshal([]byte(data), &v)
	if err != nil {
		t.Fatal(err)
	}

	jq, err := gojq.Parse(src)
	if err != nil {
		t.Fatal(err)
	}

	iter := jq.Run(v)

	for v, ok := iter.Next(); ok; v, ok = iter.Next() {
		t.Log(v)
	}
}

func TestJq_Bytes(t *testing.T) {
	data := []string{
		`"hello"`,
		`123`,
		`true`,
		`false`,
		`0.5`,
		`{"foo": "bar"}`,
		`["foo", "bar"]`,
	}

	for _, data := range data {

		var v any
		err := json.Unmarshal([]byte(data), &v)
		if err != nil {
			t.Error(data, err)
			continue
		}

		b, err := json.Marshal(v)
		if err != nil {
			t.Error(data, err)
			continue
		}

		t.Logf("in=%v out=%q type=%q\n", data, string(b), fmt.Sprintf("%T", v))
	}

}

func TestJq_Validator_status_eq_4(t *testing.T) {
	data := `{"status": 4}`
	src := ".status == 4"

	var v any
	err := json.Unmarshal([]byte(data), &v)
	if err != nil {
		t.Fatal(err)
	}

	jq, err := gojq.Parse(src)
	if err != nil {
		t.Fatal(err)
	}

	iter := jq.Run(v)

	for v, ok := iter.Next(); ok; v, ok = iter.Next() {
		switch v := v.(type) {
		case bool:
			t.Log(v)
		}
	}
}
