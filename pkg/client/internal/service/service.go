package service

import (
	"time"

	"github.com/claion-org/claiflow/pkg/client/internal/flow"
)

type ServiceStatus int32

const (
	ServiceStatusPreparing ServiceStatus = iota + 1
	ServiceStatusStart
	ServiceStatusProcessing
	ServiceStatusSuccess
	ServiceStatusFailed
)

func (s ServiceStatus) String() string {
	switch s {
	case ServiceStatusPreparing:
		return "preparing"
	case ServiceStatusStart:
		return "start"
	case ServiceStatusProcessing:
		return "processing"
	case ServiceStatusSuccess:
		return "success"
	case ServiceStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

type StepStatus int32

const (
	StepStatusPreparing = iota + 1
	StepStatusProcessing
	StepStatusSuccess
	StepStatusFail
)

func (s StepStatus) String() string {
	switch s {
	case StepStatusPreparing:
		return "preparing"
	case StepStatusProcessing:
		return "processing"
	case StepStatusSuccess:
		return "success"
	case StepStatusFail:
		return "fail"
	default:
		return "unknown"
	}
}

type Result struct {
	Body string
	Err  error
}

type StepCommand struct {
	Method string
	Args   map[string]interface{}
}

type Service struct {
	Id          string
	Name        string
	ClusterId   string
	Priority    int
	CreatedTime time.Time
	StartTime   time.Time
	UpdateTime  time.Time
	EndTime     time.Time
	Status      ServiceStatus
	Flow        flow.Flow
	// Inputs      map[string]interface{}
	Result Result
}

func (s *Service) GetId() string {
	return s.Id
}

func (s *Service) GetPriority() int {
	return s.Priority
}

func (s *Service) GetCreatedTime() time.Time {
	return s.CreatedTime
}

type UpdateService struct {
	Id        string
	StepCount int
	Sequence  int
	Status    StepStatus
	Result    Result
	Started   time.Time
	Ended     time.Time
}

func (s *UpdateService) GetId() string {
	return s.Id
}

func (s *UpdateService) GetStepCount() int {
	return s.StepCount
}

func (s *UpdateService) GetSequence() int {
	return s.Sequence
}

func (s *UpdateService) GetStatus() StepStatus {
	return s.Status
}

func (s *UpdateService) GetResult() Result {
	return s.Result
}

func (s *UpdateService) GetStarted() time.Time {
	return s.Started
}

func (s *UpdateService) GetEnded() time.Time {
	return s.Ended
}
