package posixsignal

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/neee333ko/IAM/pkg/shutdown"
)

func waitSig(c chan int, t *testing.T) {
	select {
	case <-c:

	case <-time.After(time.Second):
		t.Error("Error in PosixSignal")
	}
}

func TestDefaultPosixSignal(t *testing.T) {
	c := make(chan int, 100)

	gs := shutdown.NewGS()
	m := NewPSManager()
	gs.AddShutdownManager(m)
	gs.AddShutdownCallback(shutdown.CallBackFn(func(s string) error {
		c <- 1
		return nil
	}))
	gs.Start()

	time.Sleep(time.Millisecond)

	syscall.Kill(os.Getpid(), syscall.SIGINT)

	waitSig(c, t)

	gs.Start()

	time.Sleep(time.Millisecond)

	syscall.Kill(os.Getpid(), syscall.SIGTERM)

	waitSig(c, t)
}

func TestCustomPosixSignal(t *testing.T) {
	c := make(chan int, 100)

	gs := shutdown.NewGS()
	m := NewPSManager(syscall.SIGHUP)
	gs.AddShutdownManager(m)
	gs.AddShutdownCallback(shutdown.CallBackFn(func(s string) error {
		c <- 1
		return nil
	}))
	gs.Start()

	time.Sleep(time.Millisecond)

	syscall.Kill(os.Getpid(), syscall.SIGHUP)

	waitSig(c, t)
}
