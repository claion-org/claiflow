// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package model

import (
	"fmt"
	"strings"
)

const (
	// ClusterFieldsId is a ClusterFields of type Id.
	ClusterFieldsId ClusterFields = iota
	// ClusterFieldsUuid is a ClusterFields of type Uuid.
	ClusterFieldsUuid
	// ClusterFieldsName is a ClusterFields of type Name.
	ClusterFieldsName
	// ClusterFieldsSummary is a ClusterFields of type Summary.
	ClusterFieldsSummary
	// ClusterFieldsCreated is a ClusterFields of type Created.
	ClusterFieldsCreated
	// ClusterFieldsUpdated is a ClusterFields of type Updated.
	ClusterFieldsUpdated
)

var ErrInvalidClusterFields = fmt.Errorf("not a valid ClusterFields, try [%s]", strings.Join(_ClusterFieldsNames, ", "))

const _ClusterFieldsName = "iduuidnamesummarycreatedupdated"

var _ClusterFieldsNames = []string{
	_ClusterFieldsName[0:2],
	_ClusterFieldsName[2:6],
	_ClusterFieldsName[6:10],
	_ClusterFieldsName[10:17],
	_ClusterFieldsName[17:24],
	_ClusterFieldsName[24:31],
}

// ClusterFieldsNames returns a list of possible string values of ClusterFields.
func ClusterFieldsNames() []string {
	tmp := make([]string, len(_ClusterFieldsNames))
	copy(tmp, _ClusterFieldsNames)
	return tmp
}

var _ClusterFieldsMap = map[ClusterFields]string{
	ClusterFieldsId:      _ClusterFieldsName[0:2],
	ClusterFieldsUuid:    _ClusterFieldsName[2:6],
	ClusterFieldsName:    _ClusterFieldsName[6:10],
	ClusterFieldsSummary: _ClusterFieldsName[10:17],
	ClusterFieldsCreated: _ClusterFieldsName[17:24],
	ClusterFieldsUpdated: _ClusterFieldsName[24:31],
}

// String implements the Stringer interface.
func (x ClusterFields) String() string {
	if str, ok := _ClusterFieldsMap[x]; ok {
		return str
	}
	return fmt.Sprintf("ClusterFields(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ClusterFields) IsValid() bool {
	_, ok := _ClusterFieldsMap[x]
	return ok
}

var _ClusterFieldsValue = map[string]ClusterFields{
	_ClusterFieldsName[0:2]:                    ClusterFieldsId,
	strings.ToLower(_ClusterFieldsName[0:2]):   ClusterFieldsId,
	_ClusterFieldsName[2:6]:                    ClusterFieldsUuid,
	strings.ToLower(_ClusterFieldsName[2:6]):   ClusterFieldsUuid,
	_ClusterFieldsName[6:10]:                   ClusterFieldsName,
	strings.ToLower(_ClusterFieldsName[6:10]):  ClusterFieldsName,
	_ClusterFieldsName[10:17]:                  ClusterFieldsSummary,
	strings.ToLower(_ClusterFieldsName[10:17]): ClusterFieldsSummary,
	_ClusterFieldsName[17:24]:                  ClusterFieldsCreated,
	strings.ToLower(_ClusterFieldsName[17:24]): ClusterFieldsCreated,
	_ClusterFieldsName[24:31]:                  ClusterFieldsUpdated,
	strings.ToLower(_ClusterFieldsName[24:31]): ClusterFieldsUpdated,
}

// ParseClusterFields attempts to convert a string to a ClusterFields.
func ParseClusterFields(name string) (ClusterFields, error) {
	if x, ok := _ClusterFieldsValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _ClusterFieldsValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return ClusterFields(0), fmt.Errorf("%s is %w", name, ErrInvalidClusterFields)
}