package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

// snippet: start
type Monitor struct {
	cancel context.CancelFunc
}

func (m *Monitor) Start(ctx context.Context) context.Context {

	// start the monitor with the given context
	go m.listen(ctx)

	// create a new context that will be canceled
	// when the monitor is shut down
	ctx, cancel := context.WithCancel(context.Background())

	// hold on to the cancellation function
	// when context that started the manager is canceled
	// this cancellation function will be called.
	m.cancel = cancel

	// return the new, cancellable, context.
	// clients can listen to this context
	// for cancellation to ensure the
	// monitor is properly shut down.
	return ctx
}

// snippet: start

// snippet: listen
func (m *Monitor) listen(ctx context.Context) {
	defer m.cancel()

	// create a new ticker channel to listen to
	tick := time.NewTicker(time.Millisecond * 10)
	defer tick.Stop()

	// use an infinite loop to continue to listen
	// to new messages after the select statement
	// has been executed
	for {
		select {
		case <-ctx.Done(): // listen for context cancellation
			// shut down if the context is canceled
			fmt.Println("shutting down monitor")

			// if the monitor was told to shut down
			// then it should call its cancel function
			// so the client will know that the monitor
			// has properly shut down.
			m.cancel()

			// return from the function
			return
		case <-tick.C: // listen to the ticker channel
			// and print a message every time it ticks
			fmt.Println("monitor check")
		}
	}

}

// snippet: listen

// snippet: main
func main() {

	// create a new background context
	ctx := context.Background()

	// wrap the background context with a
	// cancellable context.
	// this context can be listened to any
	// children of this context for notification
	// of application shutdown/cancellation.
	ctx, cancel := context.WithCancel(ctx)

	// ensure the cancel function is called
	// to shut down the monitor when the program
	// is exits
	defer cancel()

	// launch a goroutine to cancel the application
	// context after a short while.
	go func() {
		time.Sleep(time.Millisecond * 50)

		// cancel the application context
		// this will shut the monitor down
		cancel()
	}()

	// create a new monitor
	mon := Monitor{}

	// start the monitor with the application context
	// this will return a context that can be listened to
	// for cancellation signaling the monitor has shut down.
	ctx = mon.Start(ctx)

	// block the application until either the context
	// is canceled or the application times out
	select {
	case <-ctx.Done(): // listen for context cancellation
		// success shutdown
		os.Exit(0)
	case <-time.After(time.Second * 2): // timeout after 2 second
		fmt.Println("timed out while trying to shut down the monitor")

		// check if there was an error from the
		// monitor's context
		if err := ctx.Err(); err != nil {
			fmt.Printf("error: %s\n", err)
		}

		// non-successful shutdown
		os.Exit(1)
	}
}

// snippet: main

/*
// snippet: out
monitor check
monitor check
monitor check
monitor check
monitor check
shutting down monitor
// snippet: out
*/
