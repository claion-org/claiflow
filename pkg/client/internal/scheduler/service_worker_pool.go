package scheduler

import (
	"context"

	"github.com/claion-org/claiflow/pkg/client/internal/executor"
	"github.com/claion-org/claiflow/pkg/client/internal/service"
	"github.com/claion-org/claiflow/pkg/client/log"
)

func runServiceWorker(ctx context.Context, id int, services chan<- chan *service.Service, updates chan<- *service.UpdateService) {
	reqCh := make(chan *service.Service)
	defer close(reqCh)
	for {
		select {
		case services <- reqCh:
			serv := <-reqCh
			if serv == nil {
				continue
			}

			log.Debugf("service-worker received service. worker_id=%d, service=%s\n", id, serv.GetId())

			se := executor.NewServiceExecutor(*serv, updates)
			se.Execute()
		case <-ctx.Done():
			log.Debugf("cancelled worker. worker_id=%d, error=%v\n", id, ctx.Err())
			return
		}
	}
}

type serviceWorkerPool struct {
	limit    int
	services chan chan *service.Service
}

func NewServiceWorkerPool(limit int) *serviceWorkerPool {
	return &serviceWorkerPool{
		limit:    limit,
		services: make(chan chan *service.Service),
	}
}

func (wp *serviceWorkerPool) Run(ctx context.Context, updates chan<- *service.UpdateService) {
	log.Debugf("running service-workers. count=%d\n", wp.limit)
	for i := 0; i < wp.limit; i++ {
		go runServiceWorker(ctx, i, wp.services, updates)
	}
}
