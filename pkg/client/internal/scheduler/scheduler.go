package scheduler

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/claion-org/claiflow/pkg/client/internal/service"
)

const minWokersCount = 5

type Scheduler struct {
	servicesStatusMap map[string]service.ServiceStatus
	lock              sync.RWMutex
	updateChan        chan *service.UpdateService // this channel receives service's status
	notifyUpdateChan  chan *service.UpdateService
	swp               *serviceWorkerPool
}

func NewScheduler(optWorkersCount int) *Scheduler {
	workersCount := minWokersCount
	if maxCpus := runtime.GOMAXPROCS(0); maxCpus > minWokersCount {
		workersCount = maxCpus
	}

	if optWorkersCount > minWokersCount {
		workersCount = optWorkersCount
	}

	return &Scheduler{
		servicesStatusMap: make(map[string]service.ServiceStatus),
		updateChan:        make(chan *service.UpdateService),
		notifyUpdateChan:  make(chan *service.UpdateService),
		swp:               NewServiceWorkerPool(workersCount)}
}

func (s *Scheduler) Start() error {
	if s.updateChan == nil || s.servicesStatusMap == nil {
		return fmt.Errorf("scheduler don't have channel")
	}

	go s.RecvNotifyServiceStatus()

	go s.swp.Run(context.Background(), s.updateChan)

	return nil
}

func (s *Scheduler) GetWorkersCount() int {
	return s.swp.limit
}

func (s *Scheduler) GetAvailableWorkerChannel() chan *service.Service {
	return <-s.swp.services
}

func (s *Scheduler) GetAvailableWorkerChannelList() []chan *service.Service {
	var chans []chan *service.Service
	cnt := 0
	for {
		select {
		case req := <-s.swp.services:
			chans = append(chans, req)
			continue
		default:
			if len(chans) <= 0 {
				req := <-s.swp.services
				chans = append(chans, req)
				continue
			} else if cnt == 0 && len(chans) == 1 {
				cnt++
				<-time.After(time.Millisecond * 100)
				continue
			}
			return chans
		}
	}
}

func (s *Scheduler) PushService(reqCh chan *service.Service, serv *service.Service) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.servicesStatusMap[serv.GetId()]; ok {
		return
	}

	s.servicesStatusMap[serv.GetId()] = service.ServiceStatusPreparing
	reqCh <- serv
}

func (s *Scheduler) PushServices(reqChs []chan *service.Service, servs []*service.Service) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, serv := range servs {
		if _, ok := s.servicesStatusMap[serv.GetId()]; ok {
			return
		}

		s.servicesStatusMap[serv.GetId()] = service.ServiceStatusPreparing

		if len(reqChs) > 0 {
			reqChs[0] <- serv
			reqChs = reqChs[1:]
		} else {
			reqCh := <-s.swp.services
			reqCh <- serv
		}
	}

	// return remain chan
	for _, ch := range reqChs {
		ch <- nil
	}
}

func (s *Scheduler) RecvNotifyServiceStatus() {
	// If you want to stop. close(s.ch).
	for update := range s.updateChan {
		s.notifyUpdateChan <- update
	}
}

func (s *Scheduler) UpdateServiceStatus(update *service.UpdateService) {
	serviceStatus := service.ServiceStatusProcessing
	if update.GetStatus() == service.StepStatusFail {
		serviceStatus = service.ServiceStatusFailed
	} else {
		if update.GetStepCount() == update.GetSequence()+1 {
			if update.GetStatus() == service.StepStatusSuccess {
				serviceStatus = service.ServiceStatusSuccess
			}
		}
	}
	s.lock.Lock()
	prevStatus, ok := s.servicesStatusMap[update.GetId()]
	if ok {
		if prevStatus < serviceStatus {
			s.servicesStatusMap[update.GetId()] = service.ServiceStatus(serviceStatus)
		}
	}
	s.lock.Unlock()
}

func (s *Scheduler) NotifyServiceUpdate() <-chan *service.UpdateService {
	return s.notifyUpdateChan
}

func (s *Scheduler) CleanupRemainingServices() map[string]service.ServiceStatus {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.servicesStatusMap) == 0 {
		return nil
	}

	services := make(map[string]service.ServiceStatus)

	for uuid, status := range s.servicesStatusMap {
		if status == service.ServiceStatusSuccess || status == service.ServiceStatusFailed {
			delete(s.servicesStatusMap, uuid)
			continue
		}
		services[uuid] = status
	}

	return services
}
