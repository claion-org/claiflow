package flow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/claion-org/claiflow/pkg/client/internal/flow/lexer"
	"github.com/claion-org/claiflow/pkg/client/internal/flow/parser"
	"github.com/itchyny/gojq"
)

const (
	predefinedFlowExprStr_INPUTS = "$inputs"
)

func hasFlowInputsExpr(s string) bool {
	if len(s) <= 0 {
		return false
	}

	length := len(predefinedFlowExprStr_INPUTS)

	if strings.HasPrefix(s, predefinedFlowExprStr_INPUTS) {
		if len(s) > length {
			if s[length] != '.' {
				return false
			}
		}
		return true
	}
	return false
}

type FlowStepInput struct {
	initKey string
	data    map[string]interface{}
}

func findReplaceDeferredInput(key interface{}, specified map[string]interface{}) (bool, interface{}, error) {
	switch kt := key.(type) {
	case string:
		if hasFlowInputsExpr(kt) {
			path := strings.TrimPrefix(kt, predefinedFlowExprStr_INPUTS)
			query, err := gojq.Parse(path)
			if err != nil {
				return true, nil, err
			}

			iter := query.Run(specified)
			var res interface{}
			for {
				v, ok := iter.Next()
				if !ok {
					break
				}
				if err, ok := v.(error); ok {
					return true, nil, err
				}

				if v != nil {
					res = v
					break
				}
			}

			if res == nil {
				return true, nil, fmt.Errorf("not found key %q", kt)
			}

			return true, res, nil
		}
	case map[string]interface{}:
		for k, v := range kt {
			found, value, err := findReplaceDeferredInput(v, specified)
			if err != nil {
				return false, nil, err
			}
			if found {
				kt[k] = value
			}
		}
	case []interface{}:
		for i, v := range kt {
			found, value, err := findReplaceDeferredInput(v, specified)
			if err != nil {
				return false, nil, err
			}
			if found {
				kt[i] = value
			}
		}
	}

	return false, nil, nil
}

func (f *FlowStepInput) FindReplaceDeferredInputsFrom(in map[string]interface{}) error {
	if hasFlowInputsExpr(f.initKey) {
		path := strings.TrimPrefix(f.initKey, predefinedFlowExprStr_INPUTS)
		query, err := gojq.Parse(path)
		if err != nil {
			return err
		}

		iter := query.Run(in)
		var res interface{}
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}
			if err, ok := v.(error); ok {
				return err
			}

			if v != nil {
				res = v
				break
			}
		}
		m, ok := res.(map[string]interface{})
		if !ok {
			return fmt.Errorf("inputs(%s) must be json object", f.initKey)
		}
		f.data = m
		return nil
	}

	_, _, err := findReplaceDeferredInput(f.data, in)
	if err != nil {
		return err
	}

	return nil
}

func (f *FlowStepInput) GetInputs() map[string]interface{} {
	return f.data
}

func splitKeyPathValPath(src string, dataset map[string]interface{}) (string, string, error) {
	p := parser.New(lexer.New(src))
	stmt := p.ParseExpressionStatement()
	if len(p.Errors()) > 0 {
		return "", "", fmt.Errorf("%v", p.Errors())
	}

	stmtStr := stmt.String(dataset)

	var one []rune
	var brace []rune
	var stepInOutKeyParts, keyPathParts []string
	for _, r := range stmtStr {
		if r == '$' && len(brace) == 0 {
			one = []rune{r}
			continue
		} else if r == '[' {
			brace = append(brace, r)
		} else if r == ']' {
			brace = brace[:len(brace)-1]
		} else if r == '.' && len(brace) == 0 {
			if len(one) > 0 {
				if one[0] == '$' {
					stepInOutKeyParts = append(stepInOutKeyParts, string(one[1:]))
				} else {
					keyPathParts = append(keyPathParts, string(one))
				}
				one = nil
			}
			continue
		}
		one = append(one, r)
	}
	if len(one) > 0 {
		if one[0] == '$' {
			stepInOutKeyParts = append(stepInOutKeyParts, string(one[1:]))
		} else {
			keyPathParts = append(keyPathParts, string(one))
		}
	}

	var stepInOutKey, keyPath string
	if len(stepInOutKeyParts) > 0 {
		stepInOutKey = strings.Join(stepInOutKeyParts, ".")
	}
	if len(keyPathParts) > 0 {
		keyPath = "." + strings.Join(keyPathParts, ".")
	}

	if after, found := cutPrefix(keyPath, ".outputs"); found {
		stepInOutKey += ".outputs"
		keyPath = after
		return stepInOutKey, keyPath, nil
	}

	if after, found := cutPrefix(keyPath, ".inputs"); found {
		stepInOutKey += ".inputs"
		keyPath = after
		return stepInOutKey, keyPath, nil
	}

	if after, found := cutPrefix(keyPath, ".key"); found {
		stepInOutKey += ".key"
		keyPath = after
		return stepInOutKey, keyPath, nil
	}

	if after, found := cutPrefix(keyPath, ".val"); found {
		stepInOutKey += ".val"
		keyPath = after
		return stepInOutKey, keyPath, nil
	}

	return stepInOutKey, keyPath, nil
}

func cutPrefix(s, prefix string) (after string, found bool) {
	if !strings.HasPrefix(s, prefix) {
		return s, false
	}
	return s[len(prefix):], true
}

func FindReplacePassedInput(key interface{}, dataset map[string]interface{}) (bool, interface{}, error) {
	switch kt := key.(type) {
	case string:
		if strings.HasPrefix(kt, "$") {
			stepInOutKey, keyPath, err := splitKeyPathValPath(kt, dataset)
			if err != nil {
				return false, nil, err
			}

			val, ok := dataset[stepInOutKey]
			if !ok {
				return true, nil, fmt.Errorf("not found key %q. src=%s", stepInOutKey, kt)
			}

			var jv interface{}
			switch vv := val.(type) {
			case string:
				if !checkJsonObjectOrArray([]byte(vv)) {
					return true, vv, nil
				}
				if err := json.Unmarshal([]byte(vv), &jv); err != nil {
					return true, nil, err
				}
			case []byte:
				if !checkJsonObjectOrArray(vv) {
					return true, vv, nil
				}
				if err := json.Unmarshal(vv, &jv); err != nil {
					return true, nil, err
				}
			case json.RawMessage:
				if !checkJsonObjectOrArray(vv) {
					return true, vv, nil
				}
				if err := json.Unmarshal(vv, &jv); err != nil {
					return true, nil, err
				}
			case map[string]interface{}, []interface{}:
				jv = vv
			default:
				b, err := json.Marshal(vv)
				if err != nil {
					return true, nil, err
				}

				if err := json.Unmarshal(b, &jv); err != nil {
					return true, nil, err
				}
			}

			query, err := gojq.Parse(keyPath)
			if err != nil {
				return true, nil, err
			}

			iter := query.Run(jv)
			var res interface{}
			for {
				v, ok := iter.Next()
				if !ok {
					break
				}
				if err, ok := v.(error); ok {
					return true, nil, err
				}

				if v != nil {
					res = v
					break
				}
			}

			if res == nil {
				return true, nil, fmt.Errorf("not found key %q", kt)
			}

			return true, res, nil
		}
	case map[string]interface{}:
		for k, v := range kt {
			found, value, err := FindReplacePassedInput(v, dataset)
			if err != nil {
				return false, nil, err
			}
			if found {
				kt[k] = value
			}
		}
	case []interface{}:
		for i, v := range kt {
			found, value, err := FindReplacePassedInput(v, dataset)
			if err != nil {
				return false, nil, err
			}
			if found {
				kt[i] = value
			}
		}
	}

	return false, nil, nil
}

func checkJsonObjectOrArray(b []byte) bool {
	x := bytes.TrimLeft(b, " \t\r\n")

	if ok := json.Valid([]byte(x)); !ok {
		return false
	}

	if len(x) > 0 {
		if x[0] == '{' || x[0] == '[' {
			return true
		}
	}

	return false
}

func (f *FlowStepInput) FindReplacePassedInputsFrom(in map[string]interface{}) error {
	_, _, err := FindReplacePassedInput(f.data, in)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------

// type FlowStep struct {
// 	Id      string        `json:"$id"`
// 	Command string        `json:"$command"`
// 	Inputs  FlowStepInput `json:"inputs,omitempty"`
// 	Outputs interface{}   `json:"outputs,omitempty"`
// 	Error   error         `json:"error,omitempty"`
// }

type FlowStep interface {
	GetId() string
	SetId(id string)
	step()
	GetOutputs() interface{}
}

var _ FlowStep = (*CommandStep)(nil)

type CommandStep struct {
	Id      string        `json:"$id"`
	Command string        `json:"$command"`
	Inputs  FlowStepInput `json:"inputs"`
	Outputs interface{}   `json:"outputs,omitempty"`
	Error   error         `json:"error,omitempty"`
}

func (cs *CommandStep) GetId() string           { return cs.Id }
func (cs *CommandStep) SetId(id string)         { cs.Id = id }
func (cs *CommandStep) step()                   {}
func (cs *CommandStep) GetOutputs() interface{} { return cs.Outputs }

var _ FlowStep = (*IfStep)(nil)

type IfStep struct {
	Id        string      `json:"$id"`
	Condition string      `json:"$condition"`
	Then      []FlowStep  `json:"$then"`
	Else      []FlowStep  `json:"$else"`
	Outputs   interface{} `json:"outputs,omitempty"`
}

func (is *IfStep) GetId() string           { return is.Id }
func (is *IfStep) SetId(id string)         { is.Id = id }
func (is *IfStep) step()                   {}
func (is *IfStep) GetOutputs() interface{} { return is.Outputs }

var _ FlowStep = (*IterationStep)(nil)

type IterationStep struct {
	Id           string      `json:"$id"`
	Range        interface{} `json:"$range"`
	Steps        []FlowStep  `json:"$steps"`
	Outputs      interface{} `json:"outputs,omitempty"`
	Error        error       `json:"error,omitempty"`
	ParentStepId string
}

func (is *IterationStep) GetId() string           { return is.Id }
func (is *IterationStep) SetId(id string)         { is.Id = id }
func (is *IterationStep) step()                   {}
func (is *IterationStep) GetOutputs() interface{} { return is.Outputs }

var _ FlowStep = (*PrintStep)(nil)

type PrintSpec struct {
	Hide interface{} `json:"hide"`
}

type PrintStep struct {
	Id       string    `json:"$id"`
	Print    PrintSpec `json:"$print"`
	FlowSpec []byte
	Outputs  interface{} `json:"outputs,omitempty"`
	Error    error       `json:"error,omitempty"`
}

func (ps *PrintStep) GetId() string           { return ps.Id }
func (ps *PrintStep) SetId(id string)         { ps.Id = id }
func (ps *PrintStep) step()                   {}
func (ps *PrintStep) GetOutputs() interface{} { return ps.Outputs }

func CopyStep(step FlowStep) FlowStep {
	switch st := step.(type) {
	case *CommandStep:
		cdata, _ := copyMap(st.Inputs.data)
		return &CommandStep{
			Id:      st.Id,
			Command: st.Command,
			Inputs: FlowStepInput{
				initKey: st.Inputs.initKey,
				data:    cdata,
			},
			Outputs: st.Outputs,
			Error:   st.Error,
		}
	case *IterationStep:
		iterStep := &IterationStep{
			Id:      st.Id,
			Range:   st.Range,
			Outputs: st.Outputs,
		}
		for _, sst := range st.Steps {
			iterStep.Steps = append(iterStep.Steps, CopyStep(sst))
		}
		return iterStep
	}
	return nil
}

func copyMap(m map[string]interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func SwitchFlowStep(flowstep FlowStep, in map[string]interface{}) error {
	switch sstep := flowstep.(type) {
	case *CommandStep:
		if err := sstep.Inputs.FindReplaceDeferredInputsFrom(in); err != nil {
			return err
		}
	case *IterationStep:
		found, val, err := findReplaceDeferredInput(sstep.Range, in)
		if err != nil {
			return err
		} else if found {
			sstep.Range = val
		}
		for _, st := range sstep.Steps {
			SwitchFlowStep(st, in)
		}
	case *PrintStep:
		found, val, err := findReplaceDeferredInput(sstep.Print.Hide, in)
		if err != nil {
			return err
		}

		if found && val != nil {
			sstep.Print.Hide = val
		}
	default:
		return fmt.Errorf("unknown step type. step_type=%T", flowstep)
	}

	return nil
}

func ReplaceStepIdAll(step FlowStep, old, new string) {
	switch sttyp := step.(type) {
	case *CommandStep:
		sttyp.Id = strings.ReplaceAll(sttyp.Id, old, new)
	case *IterationStep:
		if rstr, ok := sttyp.Range.(string); ok {
			sttyp.Range = strings.ReplaceAll(rstr, old, new)
		}
		for _, ss := range sttyp.Steps {
			ReplaceStepIdAll(ss, old, new)
		}
	}
}

//---------------------------------------------------

type Flow []FlowStep

func (f *Flow) UnmarshalJSON(b []byte) error {
	var l []json.RawMessage

	if err := json.Unmarshal(b, &l); err != nil {
		return err
	}

	ids := make(map[string]bool)

	for _, e := range l {
		m := make(map[string]interface{})
		if err := json.Unmarshal(e, &m); err != nil {
			return err
		}

		var id string
		if idInf, ok := m["$id"]; !ok {
			return fmt.Errorf("step's $id is required")
		} else {
			id = fmt.Sprintf("%v", idInf)
			if id == "" {
				return fmt.Errorf("step's $id value is required")
			}
		}

		var step FlowStep
		if cmdInf, ok := m["$command"]; ok {
			cmdStep := &CommandStep{
				Id:      id,
				Command: fmt.Sprintf("%v", cmdInf),
			}

			inputsInf, found := m["inputs"]
			if found {
				switch t := inputsInf.(type) {
				case string:
					if !hasFlowInputsExpr(t) {
						return fmt.Errorf("'string' type of 'inputs' must start with %q. not %q", predefinedFlowExprStr_INPUTS, t)
					}
					cmdStep.Inputs = FlowStepInput{initKey: t}
				case map[string]interface{}:
					cmdStep.Inputs = FlowStepInput{data: t}
				default:
					return fmt.Errorf("'inputs' type must be '$string' or 'map[string]interface{}'. not '%T'", t)
				}
			}

			step = cmdStep

		} else if rangeInf, ok := m["$range"]; ok {
			iterationStep := &IterationStep{
				Id: id,
			}
			switch t := rangeInf.(type) {
			case string:
				// if !hasFlowInputsExpr(t) {
				// 	return fmt.Errorf("'string' type of 'inputs' must start with %q. not %q", predefinedFlowExprStr_INPUTS, t)
				// }
				iterationStep.Range = t
			case map[string]interface{}:
				iterationStep.Range = t
			case []interface{}:
				iterationStep.Range = t
			default:
				return fmt.Errorf("'inputs' type must be '$string' or 'map[string]interface{}'. not '%T'", t)
			}

			if stepsInf, found := m["$steps"]; found {
				b, err := json.Marshal(stepsInf)
				if err != nil {
					return err
				}
				var f Flow
				if err := json.Unmarshal(b, &f); err != nil {
					return err
				}
				iterationStep.Steps = f
			}

			step = iterationStep

		} else if printInf, ok := m["$print"]; ok {
			printStep := &PrintStep{
				Id:       id,
				FlowSpec: b,
			}

			printSpecMap, ok := printInf.(map[string]interface{})
			if !ok {
				return fmt.Errorf("'$print' type is map[string]interface{}. not %T", printInf)
			}

			hideInf, ok := printSpecMap["hide"]
			if !ok {
				return fmt.Errorf("'hide' is required")
			}

			switch t := hideInf.(type) {
			case string:
				if !hasFlowInputsExpr(t) {
					return fmt.Errorf("'string' type of 'inputs' must start with %q. not %q", predefinedFlowExprStr_INPUTS, t)
				}
				printStep.Print.Hide = t
			case []interface{}:
				var hideList []string
				for _, v := range t {
					hideList = append(hideList, fmt.Sprintf("%v", v))
				}
				printStep.Print.Hide = hideList
			default:
				return fmt.Errorf("'hide' type must be '$string' or '[]interface{}'. not '%T'", hideInf)
			}

			step = printStep

		} else {
			return fmt.Errorf("unknown step. step=%#v", m)
		}

		if _, ok := ids[step.GetId()]; ok {
			return fmt.Errorf("step's $id(%q) already exists", step.GetId())
		}
		ids[step.GetId()] = true
		*f = append(*f, step)
	}

	return nil
}
