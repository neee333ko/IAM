package posixsignal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/neee333ko/IAM/pkg/shutdown"
	"github.com/neee333ko/log"
)

type PosixSignalManager struct {
	signals []os.Signal
}

func NewPSManager(signals ...os.Signal) shutdown.ShutdownManager {
	m := &PosixSignalManager{
		signals: make([]os.Signal, 0, 2),
	}

	if len(signals) == 0 {
		m.signals = append(m.signals, syscall.SIGINT, syscall.SIGTERM)
	}

	m.signals = append(m.signals, signals...)

	return m
}

func (m *PosixSignalManager) ShutdownStart() error {
	log.Info("Graceful Shutdown Started.")
	return nil
}

func (m *PosixSignalManager) ShutdownFinish() error {
	log.Info("Graceful Shutdown Finished.")
	os.Exit(0)
	return nil
}

func (m *PosixSignalManager) GetName() string {
	return "posixsignal"
}

func (m *PosixSignalManager) Start(gs *shutdown.GracefulShutdown) {
	go func() {
		var c chan os.Signal = make(chan os.Signal, len(m.signals))
		signal.Notify(c, m.signals...)

		<-c

		gs.StartShutdown(m)
	}()
}
