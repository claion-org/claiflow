package fetcher

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/claion-org/claiflow/pkg/client/internal/service"
	"github.com/claion-org/claiflow/pkg/client/log"
)

func (f *Fetcher) RebounceClientPod(serviceId string) {
	t := time.Now()
	log.Debugf("client-pod rebounce: start\n")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	up := service.CreateUpdateService(serviceId, 1, 0, service.StepStatusProcessing, service.Result{}, t, time.Time{})
	if err := f.serverAPI.UpdateServices(ctx, up); err != nil {
		log.Errorf("client-pod rebounce: failed to update service status('processing'): error=%v\n", err)
	}

	log.Debugf("client-pod rebounce: polling stop\n")

	done := make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
		}()

		log.Debugf("client-pod rebounce: waiting to process the remaining services: waiting_timeout=30s\n")

		for {
			<-time.After(time.Second * 3)
			services := f.RemainServices()
			if len(services) == 0 {
				break
			}

			buf := bytes.Buffer{}
			buf.WriteString("client-pod rebounce: waiting remain services:")
			for uuid, status := range services {
				buf.WriteString(fmt.Sprintf("\n\tuuid: %s, status: %s", uuid, status.String()))
			}
			log.Debugf(buf.String() + "\n")
		}
	}()

	select {
	case <-time.After(time.Second * 30):
		log.Debugf("client-pod rebounce: waiting timeout\n")
	case <-done:
		log.Debugf("client-pod rebounce: waiting done\n")
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel2()

	up = service.CreateUpdateService(serviceId, 1, 0, service.StepStatusSuccess, service.Result{Body: "client pod rebounce will be complete"}, t, time.Now())
	if err := f.serverAPI.UpdateServices(ctx2, up); err != nil {
		log.Errorf("client-pod rebounce: failed to update service status('success'): error=%v\n", err)
	}

	f.Cancel()
}
