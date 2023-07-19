package cryptography

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

// CipherString
type CipherString string

func (cs *CipherString) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		return nil
	}

	var b []byte
	switch value := value.(type) {
	case string:
		var i sql.NullString
		if err := i.Scan(value); err != nil {
			return err
		}
		b = []byte(i.String)
	case []byte:
		b = value
	default:
		return errors.New("invalid type")
	}

	bytes, err := EnigmaDecode(b)
	if err != nil {
		return errors.Wrapf(err, "default crypto string: decode")
	}

	*cs = CipherString(bytes)

	return nil
}

func (cs CipherString) Value() (driver.Value, error) {
	out, err := EnigmaEncode([]byte(cs))
	if err != nil {
		return out, errors.Wrapf(err, "default crypto string: encode")
	}

	return out, nil
}

// CipherObject
type CipherObject map[string]interface{}

func (cj *CipherObject) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		return nil
	}

	var b []byte
	switch value := value.(type) {
	case string:
		var i sql.NullString
		if err := i.Scan(value); err != nil {
			return err
		}
		b = []byte(i.String)
	case []byte:
		b = value
	default:
		return errors.New("invalid type")
	}

	bytes, err := EnigmaDecode(b)
	if err != nil {
		return errors.Wrapf(err, "default crypto hashset: decode")
	}

	if err := json.Unmarshal(bytes, &cj); err != nil {
		return errors.Wrapf(err, "default crypto hashset: json unmarshal")
	}

	return nil
}

func (cj CipherObject) Value() (driver.Value, error) {
	out, err := json.Marshal(cj)
	if err != nil {
		return string(out), errors.Wrapf(err, "default crypto hashset: json marshal")
	}
	out, err = EnigmaEncode(out)
	if err != nil {
		return string(out), errors.Wrapf(err, "default crypto hashset: encode")
	}

	return string(out), nil
}

// CipherHeader
type CipherHeader map[string][]string

func (cj *CipherHeader) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		return nil
	}

	var b []byte
	switch value := value.(type) {
	case string:
		var i sql.NullString
		if err := i.Scan(value); err != nil {
			return err
		}
		b = []byte(i.String)
	case []byte:
		b = value
	default:
		return errors.New("invalid type")
	}

	bytes, err := EnigmaDecode(b)
	if err != nil {
		return errors.Wrapf(err, "default crypto hashset: decode")
	}

	if err := json.Unmarshal(bytes, &cj); err != nil {
		return errors.Wrapf(err, "default crypto hashset: json unmarshal")
	}

	return nil
}

func (cj CipherHeader) Value() (driver.Value, error) {
	out, err := json.Marshal(cj)
	if err != nil {
		return string(out), errors.Wrapf(err, "default crypto hashset: json marshal")
	}
	out, err = EnigmaEncode(out)
	if err != nil {
		return string(out), errors.Wrapf(err, "default crypto hashset: encode")
	}

	return string(out), nil
}
