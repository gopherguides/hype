package demo

import (
	"context"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/markbates/contextual"
)

// snippet: func
func Listener(ctx context.Context, t testing.TB) {
	t.Log("waiting for context to finish")

	// wait for the context to finish
	<-ctx.Done()

}

// snippet: func

// snippet: example

// use syscall.SIGUSR2 to test
const TEST_SIGNAL = syscall.SIGUSR2

func Test_Signals(t *testing.T) {

	// create a background context
	ctx := context.Background()

	// wrap the context with a context
	// that will self cancel after 5 seconds
	// if the context is not finished
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// wrap the context with a context
	// that will self cancel if the system
	// receives a TEST_SIGNAL
	sigCtx, cancel := signal.NotifyContext(ctx, TEST_SIGNAL)
	defer cancel()

	print(t, sigCtx)

	// launch a goroutine to wait for the context
	// to finish
	go Listener(sigCtx, t)

	// snippet: kill
	// launch a goroutine to send a TEST_SIGNAL
	// to the system after 1 second
	go func() {
		time.Sleep(time.Second)

		t.Log("sending test signal")

		// send the TEST_SIGNAL to the system
		syscall.Kill(syscall.Getpid(), TEST_SIGNAL)
	}()
	// snippet: kill

	// wait for the context to finish
	select {
	case <-ctx.Done():
		t.Log("context finished")
	case <-sigCtx.Done():
		t.Log("signal received")
		t.Log("successfully completed")
		return
	}

	err := ctx.Err()
	if err == nil {
		return
	}

	// if we receive a DeadlineExceeded error then
	// the context timed out and the signal was never
	// received.
	if err == context.DeadlineExceeded {
		t.Fatal("unexpected error", err)
	}

}

// snippet: example

func print(t testing.TB, ctx context.Context) {
	t.Helper()

	s, err := contextual.String(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf(s)
}
