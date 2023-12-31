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
	// ClusterServiceStatusFieldsPdate is a ClusterServiceStatusFields of type Pdate.
	ClusterServiceStatusFieldsPdate ClusterServiceStatusFields = iota
	// ClusterServiceStatusFieldsClusterUuid is a ClusterServiceStatusFields of type Cluster_uuid.
	ClusterServiceStatusFieldsClusterUuid
	// ClusterServiceStatusFieldsUuid is a ClusterServiceStatusFields of type Uuid.
	ClusterServiceStatusFieldsUuid
	// ClusterServiceStatusFieldsCreated is a ClusterServiceStatusFields of type Created.
	ClusterServiceStatusFieldsCreated
	// ClusterServiceStatusFieldsStepMax is a ClusterServiceStatusFields of type Step_max.
	ClusterServiceStatusFieldsStepMax
	// ClusterServiceStatusFieldsStepSeq is a ClusterServiceStatusFields of type Step_seq.
	ClusterServiceStatusFieldsStepSeq
	// ClusterServiceStatusFieldsStatus is a ClusterServiceStatusFields of type Status.
	ClusterServiceStatusFieldsStatus
	// ClusterServiceStatusFieldsStarted is a ClusterServiceStatusFields of type Started.
	ClusterServiceStatusFieldsStarted
	// ClusterServiceStatusFieldsEnded is a ClusterServiceStatusFields of type Ended.
	ClusterServiceStatusFieldsEnded
	// ClusterServiceStatusFieldsMessage is a ClusterServiceStatusFields of type Message.
	ClusterServiceStatusFieldsMessage
)

var ErrInvalidClusterServiceStatusFields = fmt.Errorf("not a valid ClusterServiceStatusFields, try [%s]", strings.Join(_ClusterServiceStatusFieldsNames, ", "))

const _ClusterServiceStatusFieldsName = "pdatecluster_uuiduuidcreatedstep_maxstep_seqstatusstartedendedmessage"

var _ClusterServiceStatusFieldsNames = []string{
	_ClusterServiceStatusFieldsName[0:5],
	_ClusterServiceStatusFieldsName[5:17],
	_ClusterServiceStatusFieldsName[17:21],
	_ClusterServiceStatusFieldsName[21:28],
	_ClusterServiceStatusFieldsName[28:36],
	_ClusterServiceStatusFieldsName[36:44],
	_ClusterServiceStatusFieldsName[44:50],
	_ClusterServiceStatusFieldsName[50:57],
	_ClusterServiceStatusFieldsName[57:62],
	_ClusterServiceStatusFieldsName[62:69],
}

// ClusterServiceStatusFieldsNames returns a list of possible string values of ClusterServiceStatusFields.
func ClusterServiceStatusFieldsNames() []string {
	tmp := make([]string, len(_ClusterServiceStatusFieldsNames))
	copy(tmp, _ClusterServiceStatusFieldsNames)
	return tmp
}

var _ClusterServiceStatusFieldsMap = map[ClusterServiceStatusFields]string{
	ClusterServiceStatusFieldsPdate:       _ClusterServiceStatusFieldsName[0:5],
	ClusterServiceStatusFieldsClusterUuid: _ClusterServiceStatusFieldsName[5:17],
	ClusterServiceStatusFieldsUuid:        _ClusterServiceStatusFieldsName[17:21],
	ClusterServiceStatusFieldsCreated:     _ClusterServiceStatusFieldsName[21:28],
	ClusterServiceStatusFieldsStepMax:     _ClusterServiceStatusFieldsName[28:36],
	ClusterServiceStatusFieldsStepSeq:     _ClusterServiceStatusFieldsName[36:44],
	ClusterServiceStatusFieldsStatus:      _ClusterServiceStatusFieldsName[44:50],
	ClusterServiceStatusFieldsStarted:     _ClusterServiceStatusFieldsName[50:57],
	ClusterServiceStatusFieldsEnded:       _ClusterServiceStatusFieldsName[57:62],
	ClusterServiceStatusFieldsMessage:     _ClusterServiceStatusFieldsName[62:69],
}

// String implements the Stringer interface.
func (x ClusterServiceStatusFields) String() string {
	if str, ok := _ClusterServiceStatusFieldsMap[x]; ok {
		return str
	}
	return fmt.Sprintf("ClusterServiceStatusFields(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ClusterServiceStatusFields) IsValid() bool {
	_, ok := _ClusterServiceStatusFieldsMap[x]
	return ok
}

var _ClusterServiceStatusFieldsValue = map[string]ClusterServiceStatusFields{
	_ClusterServiceStatusFieldsName[0:5]:                    ClusterServiceStatusFieldsPdate,
	strings.ToLower(_ClusterServiceStatusFieldsName[0:5]):   ClusterServiceStatusFieldsPdate,
	_ClusterServiceStatusFieldsName[5:17]:                   ClusterServiceStatusFieldsClusterUuid,
	strings.ToLower(_ClusterServiceStatusFieldsName[5:17]):  ClusterServiceStatusFieldsClusterUuid,
	_ClusterServiceStatusFieldsName[17:21]:                  ClusterServiceStatusFieldsUuid,
	strings.ToLower(_ClusterServiceStatusFieldsName[17:21]): ClusterServiceStatusFieldsUuid,
	_ClusterServiceStatusFieldsName[21:28]:                  ClusterServiceStatusFieldsCreated,
	strings.ToLower(_ClusterServiceStatusFieldsName[21:28]): ClusterServiceStatusFieldsCreated,
	_ClusterServiceStatusFieldsName[28:36]:                  ClusterServiceStatusFieldsStepMax,
	strings.ToLower(_ClusterServiceStatusFieldsName[28:36]): ClusterServiceStatusFieldsStepMax,
	_ClusterServiceStatusFieldsName[36:44]:                  ClusterServiceStatusFieldsStepSeq,
	strings.ToLower(_ClusterServiceStatusFieldsName[36:44]): ClusterServiceStatusFieldsStepSeq,
	_ClusterServiceStatusFieldsName[44:50]:                  ClusterServiceStatusFieldsStatus,
	strings.ToLower(_ClusterServiceStatusFieldsName[44:50]): ClusterServiceStatusFieldsStatus,
	_ClusterServiceStatusFieldsName[50:57]:                  ClusterServiceStatusFieldsStarted,
	strings.ToLower(_ClusterServiceStatusFieldsName[50:57]): ClusterServiceStatusFieldsStarted,
	_ClusterServiceStatusFieldsName[57:62]:                  ClusterServiceStatusFieldsEnded,
	strings.ToLower(_ClusterServiceStatusFieldsName[57:62]): ClusterServiceStatusFieldsEnded,
	_ClusterServiceStatusFieldsName[62:69]:                  ClusterServiceStatusFieldsMessage,
	strings.ToLower(_ClusterServiceStatusFieldsName[62:69]): ClusterServiceStatusFieldsMessage,
}

// ParseClusterServiceStatusFields attempts to convert a string to a ClusterServiceStatusFields.
func ParseClusterServiceStatusFields(name string) (ClusterServiceStatusFields, error) {
	if x, ok := _ClusterServiceStatusFieldsValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _ClusterServiceStatusFieldsValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return ClusterServiceStatusFields(0), fmt.Errorf("%s is %w", name, ErrInvalidClusterServiceStatusFields)
}
