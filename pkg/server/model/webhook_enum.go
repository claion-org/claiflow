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
	// WebhookConditionValidatorNone is a WebhookConditionValidator of type None.
	WebhookConditionValidatorNone WebhookConditionValidator = iota
	// WebhookConditionValidatorJq is a WebhookConditionValidator of type Jq.
	WebhookConditionValidatorJq
)

var ErrInvalidWebhookConditionValidator = fmt.Errorf("not a valid WebhookConditionValidator, try [%s]", strings.Join(_WebhookConditionValidatorNames, ", "))

const _WebhookConditionValidatorName = "nonejq"

var _WebhookConditionValidatorNames = []string{
	_WebhookConditionValidatorName[0:4],
	_WebhookConditionValidatorName[4:6],
}

// WebhookConditionValidatorNames returns a list of possible string values of WebhookConditionValidator.
func WebhookConditionValidatorNames() []string {
	tmp := make([]string, len(_WebhookConditionValidatorNames))
	copy(tmp, _WebhookConditionValidatorNames)
	return tmp
}

var _WebhookConditionValidatorMap = map[WebhookConditionValidator]string{
	WebhookConditionValidatorNone: _WebhookConditionValidatorName[0:4],
	WebhookConditionValidatorJq:   _WebhookConditionValidatorName[4:6],
}

// String implements the Stringer interface.
func (x WebhookConditionValidator) String() string {
	if str, ok := _WebhookConditionValidatorMap[x]; ok {
		return str
	}
	return fmt.Sprintf("WebhookConditionValidator(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x WebhookConditionValidator) IsValid() bool {
	_, ok := _WebhookConditionValidatorMap[x]
	return ok
}

var _WebhookConditionValidatorValue = map[string]WebhookConditionValidator{
	_WebhookConditionValidatorName[0:4]:                  WebhookConditionValidatorNone,
	strings.ToLower(_WebhookConditionValidatorName[0:4]): WebhookConditionValidatorNone,
	_WebhookConditionValidatorName[4:6]:                  WebhookConditionValidatorJq,
	strings.ToLower(_WebhookConditionValidatorName[4:6]): WebhookConditionValidatorJq,
}

// ParseWebhookConditionValidator attempts to convert a string to a WebhookConditionValidator.
func ParseWebhookConditionValidator(name string) (WebhookConditionValidator, error) {
	if x, ok := _WebhookConditionValidatorValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _WebhookConditionValidatorValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return WebhookConditionValidator(0), fmt.Errorf("%s is %w", name, ErrInvalidWebhookConditionValidator)
}

const (
	// WebhookFieldsId is a WebhookFields of type Id.
	WebhookFieldsId WebhookFields = iota
	// WebhookFieldsUuid is a WebhookFields of type Uuid.
	WebhookFieldsUuid
	// WebhookFieldsName is a WebhookFields of type Name.
	WebhookFieldsName
	// WebhookFieldsSummary is a WebhookFields of type Summary.
	WebhookFieldsSummary
	// WebhookFieldsUrl is a WebhookFields of type Url.
	WebhookFieldsUrl
	// WebhookFieldsMethod is a WebhookFields of type Method.
	WebhookFieldsMethod
	// WebhookFieldsHeaders is a WebhookFields of type Headers.
	WebhookFieldsHeaders
	// WebhookFieldsTimeout is a WebhookFields of type Timeout.
	WebhookFieldsTimeout
	// WebhookFieldsConditionValidator is a WebhookFields of type Condition_validator.
	WebhookFieldsConditionValidator
	// WebhookFieldsConditionFilter is a WebhookFields of type Condition_filter.
	WebhookFieldsConditionFilter
	// WebhookFieldsCreated is a WebhookFields of type Created.
	WebhookFieldsCreated
	// WebhookFieldsUpdated is a WebhookFields of type Updated.
	WebhookFieldsUpdated
)

var ErrInvalidWebhookFields = fmt.Errorf("not a valid WebhookFields, try [%s]", strings.Join(_WebhookFieldsNames, ", "))

const _WebhookFieldsName = "iduuidnamesummaryurlmethodheaderstimeoutcondition_validatorcondition_filtercreatedupdated"

var _WebhookFieldsNames = []string{
	_WebhookFieldsName[0:2],
	_WebhookFieldsName[2:6],
	_WebhookFieldsName[6:10],
	_WebhookFieldsName[10:17],
	_WebhookFieldsName[17:20],
	_WebhookFieldsName[20:26],
	_WebhookFieldsName[26:33],
	_WebhookFieldsName[33:40],
	_WebhookFieldsName[40:59],
	_WebhookFieldsName[59:75],
	_WebhookFieldsName[75:82],
	_WebhookFieldsName[82:89],
}

// WebhookFieldsNames returns a list of possible string values of WebhookFields.
func WebhookFieldsNames() []string {
	tmp := make([]string, len(_WebhookFieldsNames))
	copy(tmp, _WebhookFieldsNames)
	return tmp
}

var _WebhookFieldsMap = map[WebhookFields]string{
	WebhookFieldsId:                 _WebhookFieldsName[0:2],
	WebhookFieldsUuid:               _WebhookFieldsName[2:6],
	WebhookFieldsName:               _WebhookFieldsName[6:10],
	WebhookFieldsSummary:            _WebhookFieldsName[10:17],
	WebhookFieldsUrl:                _WebhookFieldsName[17:20],
	WebhookFieldsMethod:             _WebhookFieldsName[20:26],
	WebhookFieldsHeaders:            _WebhookFieldsName[26:33],
	WebhookFieldsTimeout:            _WebhookFieldsName[33:40],
	WebhookFieldsConditionValidator: _WebhookFieldsName[40:59],
	WebhookFieldsConditionFilter:    _WebhookFieldsName[59:75],
	WebhookFieldsCreated:            _WebhookFieldsName[75:82],
	WebhookFieldsUpdated:            _WebhookFieldsName[82:89],
}

// String implements the Stringer interface.
func (x WebhookFields) String() string {
	if str, ok := _WebhookFieldsMap[x]; ok {
		return str
	}
	return fmt.Sprintf("WebhookFields(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x WebhookFields) IsValid() bool {
	_, ok := _WebhookFieldsMap[x]
	return ok
}

var _WebhookFieldsValue = map[string]WebhookFields{
	_WebhookFieldsName[0:2]:                    WebhookFieldsId,
	strings.ToLower(_WebhookFieldsName[0:2]):   WebhookFieldsId,
	_WebhookFieldsName[2:6]:                    WebhookFieldsUuid,
	strings.ToLower(_WebhookFieldsName[2:6]):   WebhookFieldsUuid,
	_WebhookFieldsName[6:10]:                   WebhookFieldsName,
	strings.ToLower(_WebhookFieldsName[6:10]):  WebhookFieldsName,
	_WebhookFieldsName[10:17]:                  WebhookFieldsSummary,
	strings.ToLower(_WebhookFieldsName[10:17]): WebhookFieldsSummary,
	_WebhookFieldsName[17:20]:                  WebhookFieldsUrl,
	strings.ToLower(_WebhookFieldsName[17:20]): WebhookFieldsUrl,
	_WebhookFieldsName[20:26]:                  WebhookFieldsMethod,
	strings.ToLower(_WebhookFieldsName[20:26]): WebhookFieldsMethod,
	_WebhookFieldsName[26:33]:                  WebhookFieldsHeaders,
	strings.ToLower(_WebhookFieldsName[26:33]): WebhookFieldsHeaders,
	_WebhookFieldsName[33:40]:                  WebhookFieldsTimeout,
	strings.ToLower(_WebhookFieldsName[33:40]): WebhookFieldsTimeout,
	_WebhookFieldsName[40:59]:                  WebhookFieldsConditionValidator,
	strings.ToLower(_WebhookFieldsName[40:59]): WebhookFieldsConditionValidator,
	_WebhookFieldsName[59:75]:                  WebhookFieldsConditionFilter,
	strings.ToLower(_WebhookFieldsName[59:75]): WebhookFieldsConditionFilter,
	_WebhookFieldsName[75:82]:                  WebhookFieldsCreated,
	strings.ToLower(_WebhookFieldsName[75:82]): WebhookFieldsCreated,
	_WebhookFieldsName[82:89]:                  WebhookFieldsUpdated,
	strings.ToLower(_WebhookFieldsName[82:89]): WebhookFieldsUpdated,
}

// ParseWebhookFields attempts to convert a string to a WebhookFields.
func ParseWebhookFields(name string) (WebhookFields, error) {
	if x, ok := _WebhookFieldsValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _WebhookFieldsValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return WebhookFields(0), fmt.Errorf("%s is %w", name, ErrInvalidWebhookFields)
}
