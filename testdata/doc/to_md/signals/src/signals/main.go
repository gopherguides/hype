package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// snippet: example
func main() {

	// create a background context
	ctx := context.Background()

	// wrap the context with a timeout
	// of 50 milliseconds to ensure the application
	// will eventually exit
	ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()

	// wrap the context with a context
	// that will be cancelled when an
	// interrupt signal is received (ctrl-c)
	ctx, cancel = signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// lauch a goroutine that will
	// trigger an interrupt signal
	// after 10 milliseconds (ctrl-c)
	go func() {
		time.Sleep(10 * time.Millisecond)

		fmt.Println("sending ctrl-c")

		// send the interrupt signal
		// to the current process
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	fmt.Println("waiting for context to finish")

	// wait for the context to finish
	<-ctx.Done()

	fmt.Printf("context finished: %v\n", ctx.Err())

}

// snippet: example
