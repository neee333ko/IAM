package shutdown

import "sync"

type ShutdownCallBack interface {
	OnShutdown(string) error
}

type CallBackFn func(string) error

func (cb CallBackFn) OnShutdown(ManagerName string) error {
	return cb(ManagerName)
}

type ErrorHandler interface {
	OnError(error)
}

type HandleError func(error)

func (h HandleError) OnError(err error) {
	h(err)
}

type ShutdownManager interface {
	GetName() string
	Start(*GracefulShutdown)
	ShutdownStart() error
	ShutdownFinish() error
}

type GracefulShutdown struct {
	managers     []ShutdownManager
	callbacks    []ShutdownCallBack
	errorHandler ErrorHandler
}

func NewGS() *GracefulShutdown {
	return &GracefulShutdown{
		managers:  make([]ShutdownManager, 0, 3),
		callbacks: make([]ShutdownCallBack, 0, 10),
	}
}

func (gs *GracefulShutdown) Start() {
	for _, m := range gs.managers {
		m.Start(gs)
	}
}

func (gs *GracefulShutdown) SetErrorHandler(h ErrorHandler) {
	gs.errorHandler = h
}

func (gs *GracefulShutdown) AddShutdownManager(m ShutdownManager) {
	gs.managers = append(gs.managers, m)
}

func (gs *GracefulShutdown) AddShutdownCallback(cb ShutdownCallBack) {
	gs.callbacks = append(gs.callbacks, cb)
}

func (gs *GracefulShutdown) HandleError(err error) {
	if err != nil && gs.errorHandler != nil {
		gs.errorHandler.OnError(err)
	}
}

func (gs *GracefulShutdown) StartShutdown(m ShutdownManager) {
	gs.HandleError(m.ShutdownStart())

	var wg sync.WaitGroup
	for _, cb := range gs.callbacks {
		wg.Add(1)
		go func(ShutdownCallBack) {
			defer wg.Done()
			gs.HandleError(cb.OnShutdown(m.GetName()))
		}(cb)
	}
	wg.Wait()

	gs.HandleError(m.ShutdownFinish())
}
