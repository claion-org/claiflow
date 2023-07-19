package executor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/claion-org/claiflow/pkg/client/internal"
	"github.com/claion-org/claiflow/pkg/client/internal/flow"
	"github.com/claion-org/claiflow/pkg/client/internal/service"
	"github.com/claion-org/claiflow/pkg/client/log"
)

type ServiceExecutor struct {
	service       service.Service
	updateChannel chan<- *service.UpdateService
}

func NewServiceExecutor(service service.Service, updateChannel chan<- *service.UpdateService) *ServiceExecutor {
	return &ServiceExecutor{service: service, updateChannel: updateChannel}
}

func (se *ServiceExecutor) Execute() (err error) {
	// var result service.Result
	var result string
	flowStore := make(map[string]interface{})

	defer func() {
		se.service.EndTime = time.Now()
		if err != nil {
			se.service.Status = service.ServiceStatusFailed
			se.service.Result.Err = err
		} else {
			se.service.Status = service.ServiceStatusSuccess
			se.service.Result = service.Result{Body: result}
		}
	}()

	se.service.StartTime = time.Now()
	se.service.Status = service.ServiceStatusProcessing

	for i, step := range se.service.Flow {
		log.Infof("step start: cluster_uuid=%s, service_uuid=%s, step_position=%d\n", se.service.ClusterId, se.service.Id, i)
		stepStartTime := time.Now()

		// update step_status_processing to service scheduler through returnChannel.
		se.SendServiceStatusUpdate(i, service.StepStatusProcessing, service.Result{}, stepStartTime, time.Time{})
		if err = execFlowStep(step, flowStore); err != nil {
			se.SendServiceStatusUpdate(i, service.StepStatusFail, service.Result{Err: err}, stepStartTime, time.Now())
			log.Errorf("failed to execute step: cluster_uuid=%s, service_uuid=%s, step_position=%d, error=%v\n", se.service.ClusterId, se.service.Id, i, err)
			return err
		}
		log.Debugf("executed step: cluster_uuid=%s, service_uuid=%s, step_position=%d\n", se.service.ClusterId, se.service.Id, i)

		// update step_status_success to service scheduler through returnChannel.
		result, err := json.Marshal(step.GetOutputs())
		if err != nil {
			se.SendServiceStatusUpdate(i, service.StepStatusFail, service.Result{Err: err}, stepStartTime, time.Now())
			log.Errorf("failed to execute step: cluster_uuid=%s, service_uuid=%s, step_position=%d, error=%v\n", se.service.ClusterId, se.service.Id, i, err)
			return err
		}
		se.SendServiceStatusUpdate(i, service.StepStatusSuccess, service.Result{Body: string(result)}, stepStartTime, time.Now())
		log.Infof("step end: cluster_uuid=%s, service_uuid=%s, step_position=%d\n", se.service.ClusterId, se.service.Id, i)
	}

	return nil
}

func execFlowStep(step flow.FlowStep, flowStore map[string]interface{}) error {
	switch styp := step.(type) {
	case *flow.CommandStep:
		cmdStep := styp
		var te *StepExecutor

		if err := cmdStep.Inputs.FindReplacePassedInputsFrom(flowStore); err != nil {
			return fmt.Errorf("failed to find/replace service. error=%v", err)
		}

		if stepInputs := cmdStep.Inputs.GetInputs(); stepInputs != nil {
			flowStore[cmdStep.Id+".inputs"] = cmdStep.Inputs.GetInputs()
		}

		te, err := NewStepExecutor(cmdStep)
		if err != nil {
			return err
		}

		log.Debugf("prepare to execute step: step_command=%s, step_input=%v\n", cmdStep.Command, cmdStep.Inputs.GetInputs())
		result, err := te.Execute()
		if err != nil {
			log.Errorf("failed to execute step: step_command=%s, step_input=%v, error=%v\n", cmdStep.Command, cmdStep.Inputs.GetInputs(), err)
			return err
		}
		log.Debugf("executed step: step_command=%s\n", cmdStep.Command)
		flowStore[cmdStep.Id+".outputs"] = result
		cmdStep.Outputs = result
	case *flow.IterationStep:
		iterationStep := styp
		found, val, err := flow.FindReplacePassedInput(iterationStep.Range, flowStore)
		if err != nil {
			return fmt.Errorf("failed to find/replace service. error=%v", err)
		} else if found {
			iterationStep.Range = val
		}

		switch irt := iterationStep.Range.(type) {
		case map[string]interface{}:
			return fmt.Errorf("unsupported map type in range")
		case []interface{}:
			stepId := step.GetId()
			flowStore[fmt.Sprintf("%s.len", stepId)] = len(irt)
			for key, val := range irt {
				flowStore[fmt.Sprintf("%s.%s", step.GetId(), "key")] = key
				flowStore[fmt.Sprintf("%s.%s", step.GetId(), "val")] = val
				for _, st := range iterationStep.Steps {
					switch st.(type) {
					case *flow.CommandStep:
						cmdStep := flow.CopyStep(st)
						cmdStep.(*flow.CommandStep).Id = fmt.Sprintf("%s[%d].%s", step.GetId(), key, cmdStep.GetId())

						flow.ReplaceStepIdAll(cmdStep, fmt.Sprintf("$%s.%s", step.GetId(), "key"), strconv.Itoa(key))
						if err := execFlowStep(cmdStep, flowStore); err != nil {
							return err
						}
						if v := iterationStep.Outputs; v == nil {
							// flowStore[stepId+".outputs"] = []interface{}{cmdStep.(*CommandStep).Outputs}
							iterationStep.Outputs = []interface{}{cmdStep.(*flow.CommandStep).Outputs}
						} else {
							// flowStore[stepId+".outputs"] = append(v.([]interface{}), cmdStep.(*CommandStep).Outputs)
							iterationStep.Outputs = append(v.([]interface{}), cmdStep.(*flow.CommandStep).Outputs)
						}
					case *flow.IterationStep:
						iterStep := flow.CopyStep(st)
						iterStep.(*flow.IterationStep).Id = fmt.Sprintf("%s[%d].%s", step.GetId(), key, iterStep.GetId())

						flow.ReplaceStepIdAll(iterStep, fmt.Sprintf("$%s.%s", step.GetId(), "key"), strconv.Itoa(key))
						iterStep.(*flow.IterationStep).ParentStepId = fmt.Sprintf("%s[%d]", step.GetId(), key)
						if err := execFlowStep(iterStep, flowStore); err != nil {
							return err
						}
						if v := iterationStep.Outputs; v == nil {
							// flowStore[stepId+".outputs"] = []interface{}{cmdStep.(*CommandStep).Outputs}
							iterationStep.Outputs = []interface{}{iterStep.(*flow.IterationStep).Outputs}
						} else {
							// flowStore[stepId+".outputs"] = append(v.([]interface{}), cmdStep.(*CommandStep).Outputs)
							iterationStep.Outputs = append(v.([]interface{}), iterStep.(*flow.IterationStep).Outputs)
						}
					default:
						return fmt.Errorf("unknown step type. step_type=%T", step)
					}
				}
			}
		}
	case *flow.PrintStep:
		printStep := styp

		var flow flow.Flow
		if err := json.Unmarshal(printStep.FlowSpec, &flow); err != nil {
			return err
		}

		m := make(map[string]bool)
		switch ht := printStep.Print.Hide.(type) {
		case []interface{}:
			for _, v := range ht {
				m[fmt.Sprintf("%v", v)] = true
			}
		case []string:
			for _, v := range ht {
				m[v] = true
			}
		default:
			return fmt.Errorf("'hide' type is not []interface{} or []string. hide_type=%T", printStep.Print.Hide)
		}

		flowBytes, err := printSelectedFlow(m, flow, flowStore)
		if err != nil {
			return err
		}

		flowStore[printStep.Id+".outputs"] = string(flowBytes)
		printStep.Outputs = json.RawMessage(flowBytes)
	default:
		return fmt.Errorf("unknown step type. step_type=%T", step)
	}

	return nil
}

func (se *ServiceExecutor) SendServiceStatusUpdate(seq int, status service.StepStatus, result service.Result, st, et time.Time) {
	if se.updateChannel != nil {
		update := service.UpdateService{
			Id:        se.service.Id,
			StepCount: len(se.service.Flow),
			Sequence:  seq,
			Status:    status,
			Result:    result,
			Started:   st,
			Ended:     et,
		}

		se.updateChannel <- &update
	}
}

type StepExecutor struct {
	commandFunc interface{}
	step        *flow.CommandStep
}

func NewStepExecutor(step *flow.CommandStep) (*StepExecutor, error) {
	cmdFunc, err := internal.GetCommand(step.Command)
	if err != nil {
		return nil, err
	}

	return &StepExecutor{commandFunc: cmdFunc, step: step}, nil
}

func (se *StepExecutor) Execute() (interface{}, error) {
	res, err := ExecuteCommandFuncFromInputsMap(se.commandFunc, se.step.Inputs.GetInputs())
	if err != nil {
		return nil, err
	}

	switch res := res.(type) {
	case string:
		if len(res) > 0 && (res[0] == '{' || res[0] == '[') {
			return json.RawMessage(res), nil
		}
		return res, nil
	case []byte:
		if len(res) > 0 && (res[0] == '{' || res[0] == '[') {
			return json.RawMessage(res), nil
		}
		return res, nil
	case nil:
		return nil, nil
	}

	return res, nil
}

func printSelectedFlow(hideSteps map[string]bool, flow flow.Flow, flowStore map[string]interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("[")
	firstVisited := false
	for _, step := range flow {
		if _, ok := hideSteps[step.GetId()]; !ok {
			b, err := printSelectedFlowStep(hideSteps, step, flowStore)
			if err != nil {
				return nil, err
			}
			if !firstVisited {
				firstVisited = true
			} else if len(b) > 0 {
				buf.WriteString(",")
			}
			buf.Write(b)
		}
	}
	buf.WriteString("]")
	return buf.Bytes(), nil
}

func printSelectedFlowStep(hideSteps map[string]bool, step flow.FlowStep, flowStore map[string]interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	switch sstep := step.(type) {
	case *flow.CommandStep:
		buf.WriteString("{")
		buf.WriteString(fmt.Sprintf("\"id\":%q", sstep.GetId()))
		buf.WriteString(fmt.Sprintf(",\"command\":%q", sstep.Command))
		if in := flowStore[sstep.GetId()+".inputs"]; in != nil {
			b, _ := json.Marshal(in)
			buf.WriteString(",\"inputs\":")
			buf.Write(b)
		}
		if out := flowStore[sstep.GetId()+".outputs"]; out != nil {
			switch ot := out.(type) {
			case string:
				if len(ot) > 0 && (ot[0] == '{' || ot[0] == '[') {
					out = json.RawMessage(ot)
				}
			}
			b, _ := json.Marshal(out)
			buf.WriteString(",\"outputs\":")
			buf.Write(b)
		}
		buf.WriteString("}")
	case *flow.IterationStep:
		buf.WriteString("[")
		rangeLenInf, ok := flowStore[sstep.Id+".len"]
		if !ok {
			return nil, fmt.Errorf("not found key. key=%s", sstep.Id+".len")
		}
		rangeLen, ok := rangeLenInf.(int)
		if !ok {
			return nil, fmt.Errorf("failed type assertion to int. got=%T", rangeLenInf)
		}

		firstVisited1 := false
		for i := 0; i < rangeLen; i++ {
			if firstVisited1 {
				buf.WriteString(",")
			}
			firstVisited2 := false
			for _, st := range sstep.Steps {
				if _, ok := hideSteps[st.GetId()]; !ok {
					if !firstVisited2 {
						firstVisited1 = true
						firstVisited2 = true
					} else {
						buf.WriteString(",")
					}
					copyStep := flow.CopyStep(st)
					copyStep.SetId(fmt.Sprintf("%s[%d].%s", step.GetId(), i, copyStep.GetId()))
					b, err := printSelectedFlowStep(hideSteps, copyStep, flowStore)
					if err != nil {
						return nil, err
					}
					buf.Write(b)
				}
			}
		}
		buf.WriteString("]")
	}
	return buf.Bytes(), nil
}

func ExecuteCommandFuncFromInputsMap(cmdFunc interface{}, inputs map[string]interface{}) (interface{}, error) {
	rvCmdFunc := reflect.ValueOf(cmdFunc)
	if rvCmdFunc.Type().NumIn() != 1 {
		return executeCommandFunc(cmdFunc, nil)
	}

	rti0 := reflect.ValueOf(cmdFunc).Type().In(0)
	if rti0.Kind() == reflect.Struct || (rti0.Kind() == reflect.Pointer && rti0.Elem().Kind() == reflect.Struct) {
		if inputs == nil {
			return executeCommandFunc(cmdFunc, nil)
		}

		b, err := json.Marshal(inputs)
		if err != nil {
			return nil, err
		}

		var rvInputs reflect.Value
		if rti0.Kind() == reflect.Pointer {
			rvInputs = reflect.New(rti0.Elem())
		} else {
			rvInputs = reflect.New(rti0)
		}

		if err := json.Unmarshal(b, rvInputs.Interface()); err != nil {
			return nil, err
		}

		if rti0.Kind() != reflect.Pointer {
			return executeCommandFunc(cmdFunc, rvInputs.Elem().Interface())
		}

		return executeCommandFunc(cmdFunc, rvInputs.Interface())
	} else if rti0.Kind() == reflect.Map || (rti0.Kind() == reflect.Pointer && rti0.Elem().Kind() == reflect.Map) {
		return executeCommandFunc(cmdFunc, inputs)
	}

	return nil, fmt.Errorf("input parameter must be struct or map or struct pointer or map pointer")
}

func executeCommandFunc(cmdFunc interface{}, inputs interface{}) (interface{}, error) {
	rvCmdFunc := reflect.ValueOf(cmdFunc)
	var rvInputs []reflect.Value

	if rvCmdFunc.Type().NumIn() == 1 {
		if inputs == nil {
			rvInputs = append(rvInputs, reflect.New(rvCmdFunc.Type().In(0)).Elem())
		} else {
			rvIn := reflect.ValueOf(inputs)
			if rvCmdFunc.Type().In(0).Kind() != rvIn.Kind() {
				return nil, fmt.Errorf("cmdFunc's input parameter type is not matched. wont=%s, got=%s", rvCmdFunc.Type().In(0).Kind(), rvIn.Kind())
			}
			rvInputs = append(rvInputs, reflect.ValueOf(inputs))
		}
	} else if rvCmdFunc.Type().NumIn() != 0 {
		return nil, fmt.Errorf("cmdFunc's input parameter count must be 0 or 1")
	}

	rvOutputs := rvCmdFunc.Call(rvInputs)

	if len(rvOutputs) < 1 || len(rvOutputs) > 2 {
		return nil, fmt.Errorf("function should return either error or (result, error), got %d return values", len(rvOutputs))
	}

	var err error
	if errInterface := rvOutputs[len(rvOutputs)-1].Interface(); errInterface != nil {
		var ok bool
		if err, ok = errInterface.(error); !ok {
			return nil, fmt.Errorf("failed to serialize error result as it is not of error interface: %v", errInterface)
		}
	}

	var res interface{}
	if len(rvOutputs) > 1 && (rvOutputs[0].Kind() != reflect.Ptr || !rvOutputs[0].IsNil()) {
		res = rvOutputs[0].Interface()
	}
	return res, err
}
