package api_test

import (
	"encoding/json"
	"reflect"
	"testing"
)

type Object struct {
	String  string  `json:"string,omitempty"`
	StringP *string `json:"stringP,omitempty"`
	Int     int     `json:"int,omitempty"`
	IntP    *int    `json:"intP,omitempty"`
}

func TestJsonMarshal(t *testing.T) {

	ss := []any{
		nil,
		Pointer("foo"),
		"bar",
		Object{
			String: "", StringP: nil,
			Int:  0,
			IntP: nil,
		},
		Object{
			String: "abc", StringP: Pointer("abc"),
			Int:  123,
			IntP: Pointer(123),
		},
	}

	for _, it := range ss {
		j, err := json.MarshalIndent(it, "", "  ")
		if err != nil {
			t.Error(err, "failed to JSON marshal")
			continue
		}

		t.Log(string(j))
	}
}

func TestJsonUnmarshal(t *testing.T) {

	type Vaild struct {
		Raw json.RawMessage
		Ref any
	}

	var obj Object
	var objEmpty Object

	ss := []Vaild{
		{[]byte(`"foo"`), reflect.ValueOf("").Interface()},
		{[]byte(`123`), reflect.ValueOf(int(0)).Interface()},
		{[]byte(`{
			"string": "abc",
			"stringP": "abc",
			"int": 123,
			"intP": 123
		  }`), &obj},
		{[]byte(`{
			"string": "",
			"stringP": "",
			"int": 0,
			"intP": 0
		  }`), &obj},
		{[]byte(`{}`), &objEmpty},
	}

	for _, it := range ss {
		if err := json.Unmarshal([]byte(it.Raw), &it.Ref); err != nil {
			t.Error(err, "failed to JSON marshal")
			continue
		}

		t.Log(it.Ref)
	}
}

func Pointer[A any](a A) *A {
	return &a
}
