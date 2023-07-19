package service

import (
	"encoding/json"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/claion-org/claiflow/pkg/client/internal/flow"
	"github.com/claion-org/claiflow/pkg/client/log"
	apiclient "github.com/claion-org/claiflow/pkg/server/api/client"
	"github.com/claion-org/claiflow/pkg/server/model"
)

type Object = map[string]interface{}

func ConvertServiceListServerToClient(server []*apiclient.ServicePollingResponseV1_Data) ([]*Service, []*UpdateService) {
	var client []*Service
	var failed []*UpdateService
	for _, data := range server {
		inputs := Object{}
		if err := json.Unmarshal(data.GetInputs(), &inputs); err != nil {
			log.Warnf("failed to convert service: service_uuid=%s, error=%v\n", data.GetUuid(), err)
			t := time.Now()
			failed = append(failed, CreateUpdateService(data.GetUuid(), int(data.GetStepMax()), 0, StepStatusFail, Result{Err: err}, t, t))
			continue
		}

		serv := &Service{
			Id:          data.GetUuid(),
			Name:        data.GetName(),
			ClusterId:   data.GetClusterUuid(),
			Priority:    int(data.GetPriority()),
			CreatedTime: data.Created.AsTime(),
		}

		var fl flow.Flow
		if err := json.Unmarshal([]byte(data.GetFlow()), &fl); err != nil {
			log.Warnf("failed to convert service: service_uuid=%s, error=%v\n", data.GetUuid(), err)
			t := time.Now()
			failed = append(failed, CreateUpdateService(data.GetUuid(), int(data.GetStepMax()), 0, StepStatusFail, Result{Err: err}, t, t))
			continue
		}

		if len(fl) <= 0 {
			log.Warnf("service steps is empty: service_uuid=%s\n", data.GetUuid())
			continue
		}

		isFailed := false
		for _, flowstep := range fl {
			if err := flow.SwitchFlowStep(flowstep, inputs); err != nil {
				log.Warnf("failed to convert service: service_uuid=%s, error=%v\n", data.GetUuid(), err)
				t := time.Now()
				failed = append(failed, CreateUpdateService(data.GetUuid(), int(data.GetStepMax()), 0, StepStatusFail, Result{Err: err}, t, t))
				isFailed = true
				break
			}
		}
		if isFailed {
			continue
		}

		serv.Flow = fl

		client = append(client, serv)
	}

	return client, failed
}

func ConvertServiceStepUpdateClientToServer(client *UpdateService) *apiclient.UpdateServiceStatusRequestV1 {
	server := &apiclient.UpdateServiceStatusRequestV1{
		Uuid:     client.GetId(),
		Sequence: int32(client.GetSequence()),
		Result:   client.GetResult().Body,
		Started:  timestamppb.New(client.GetStarted()),
		Ended:    timestamppb.New(client.GetEnded()),
	}

	if client.GetResult().Err != nil {
		server.Error = client.GetResult().Err.Error()
	}

	switch client.GetStatus() {
	case StepStatusPreparing, StepStatusProcessing:
		server.Status = int32(model.StepStatusProcessing)
	case StepStatusSuccess:
		server.Status = int32(model.StepStatusSucceeded)
	case StepStatusFail:
		server.Status = int32(model.StepStatusFailed)
	}

	return server
}

func CreateUpdateService(id string, seqCount, seq int, status StepStatus, result Result, start, end time.Time) *UpdateService {
	return &UpdateService{
		Id:        id,
		StepCount: seqCount,
		Sequence:  seq,
		Status:    status,
		Result:    result,
		Started:   start,
		Ended:     end,
	}
}
