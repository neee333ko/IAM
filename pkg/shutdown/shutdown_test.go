package shutdown_test

import (
	"testing"
	"time"

	"github.com/neee333ko/IAM/pkg/shutdown"
	mock "github.com/neee333ko/IAM/pkg/shutdown/mock"
	"go.uber.org/mock/gomock"
)

func waitSig(c chan int, t *testing.T) {
	select {
	case <-c:

	case <-time.After(time.Second):
		t.Error("Error in PosixSignal")
	}
}

func TestShutdown(t *testing.T) {
	c := make(chan int, 100)

	ctrl := gomock.NewController(t)
	m := mock.NewMockShutdownManager(ctrl)
	gs := shutdown.NewGS()

	m.EXPECT().Start(gs).Do(func(gs *shutdown.GracefulShutdown) {
		gs.StartShutdown(m)
	})
	m.EXPECT().ShutdownStart().DoAndReturn(func() error {
		return nil
	})
	m.EXPECT().ShutdownFinish().DoAndReturn(func() error {
		return nil
	})
	m.EXPECT().GetName().DoAndReturn(func() string {
		return "nil"
	})

	gs.AddShutdownManager(m)
	gs.AddShutdownCallback(shutdown.CallBackFn(
		func(s string) error {
			c <- 1
			return nil
		},
	))
	gs.Start()

	waitSig(c, t)
}
