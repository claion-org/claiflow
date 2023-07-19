package internal

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	commandRegistryLock sync.Mutex
	commandRegistry     map[string]interface{} = make(map[string]interface{})
)

func RegisterCommand(name string, cmdFunc interface{}) error {
	commandRegistryLock.Lock()
	defer commandRegistryLock.Unlock()

	if _, ok := commandRegistry[name]; ok {
		return fmt.Errorf("command's name(%q) has already been registered", name)
	}

	rt := reflect.TypeOf(cmdFunc)
	if rt.Kind() != reflect.Func {
		return fmt.Errorf("cmdFunc must be func type")
	}

	numIn, numOut := rt.NumIn(), rt.NumOut()

	if numIn < 0 || numIn > 1 {
		return fmt.Errorf("cmdFunc's input parameter count must be 0 or 1")
	}

	if numIn == 1 {
		// input parameter must be struct or map or struct pointer or map pointer
		rti0 := rt.In(0)
		if rti0.Kind() != reflect.Struct && rti0.Kind() != reflect.Map {
			if rti0.Kind() != reflect.Pointer {
				return fmt.Errorf("cmdFunc's input parameter must be struct or map or struct pointer or map pointer, not %s", rti0.Kind())
			} else {
				rti0e := rti0.Elem()
				if rti0e.Kind() != reflect.Struct && rti0e.Kind() != reflect.Map {
					return fmt.Errorf("cmdFunc's input parameter must be struct or map or struct pointer or map pointer, not %s pointer", rti0e.Kind())
				}
			}
		}
	}

	if numOut < 1 || numOut > 2 {
		return fmt.Errorf("cmdFunc's output parameter count must be 1 or 2")
	}

	commandRegistry[name] = cmdFunc

	return nil
}

func GetCommand(name string) (interface{}, error) {
	commandRegistryLock.Lock()
	defer commandRegistryLock.Unlock()

	cmdFunc, ok := commandRegistry[name]
	if !ok {
		return nil, fmt.Errorf("command(%q) is not registered", name)
	}

	return cmdFunc, nil
}
