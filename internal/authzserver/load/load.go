package load

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/neee333ko/IAM/pkg/storage"
	"github.com/neee333ko/log"
)

type Loader interface {
	Reload() error
}

type Reloader struct {
	ctx context.Context
	l   Loader
}

func NewReloader(ctx context.Context, l Loader) *Reloader {
	return &Reloader{
		ctx: ctx,
		l:   l,
	}
}

func shouldReload() ([]func(), bool) {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	if len(reloadQueue) == 0 {
		return nil, false
	}

	fns := reloadQueue
	reloadQueue = make([]func(), 0)

	return fns, true
}

var (
	reloadChan  chan func()   = make(chan func())
	reloadQueue []func()      = make([]func(), 0)
	queueMutex  *sync.RWMutex = new(sync.RWMutex)
)

func (r *Reloader) Run() {
	go r.Listen()
	go r.ReloadChanLoop()
	go r.ReloadQueueLoop()
	r.DoReload()
}

func (r *Reloader) Listen() {
	redisclient := &storage.RedisCluster{}

	for {
		err := redisclient.StartPubSubHandler(Channel, redisSignalHandler(nil, nil))

		if errors.Is(err, storage.ErrRedisIsDown) {
			log.Warn("connected to redis failed, retry in 10 mins.")
		}

		log.Warn("retry to connect to redis in 10 mins.")
		time.Sleep(10 * time.Minute)
	}
}

func (r *Reloader) ReloadQueueLoop(cbs ...func()) {
	timer := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-r.ctx.Done():
			return
		case <-timer.C:
			fns, ok := shouldReload()
			if !ok {
				continue
			}
			startTime := time.Now()

			r.DoReload()

			for _, fn := range fns {
				if fn != nil {
					fn()
				}
			}

			if len(cbs) != 0 {
				cbs[0]()
			}

			log.Info(fmt.Sprintf("reload completed, cost %v seconds.\n", time.Since(startTime).Seconds()))
		}
	}
}

func (r *Reloader) ReloadChanLoop(cbs ...func()) {
	for {
		select {
		case <-r.ctx.Done():
			return
		case fn := <-reloadChan:
			queueMutex.Lock()
			reloadQueue = append(reloadQueue, fn)
			queueMutex.Unlock()

			log.Info("reload queued.")

			if len(cbs) != 0 {
				cbs[0]()
			}
		}
	}
}

func (r *Reloader) DoReload() {
	log.Info("reloading...")

	if err := r.l.Reload(); err != nil {
		log.Error(fmt.Sprintf("reloading error: %s", err.Error()))
	}
}
