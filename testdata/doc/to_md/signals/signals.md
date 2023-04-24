# Listening for System Signals with Context

Previously, when discussing channels, we saw how to capture system signals, such as `ctrl-c`, using <godoc>os/signal#Notify</godoc>. The <godoc>os/signal#NotifyContext</godoc> function, <ref>notify.doc</ref>, is a variant of <godoc>os/signal#Notify</godoc> that takes a <godoc>context#Context</godoc> as an argument. In return, we are given a <godoc>context#Context</godoc> that will be canceled when the signal is received.

<figure id="notify.doc" type="listing">

<go doc="os/signal.NotifyContext"></go>

<figcaption>The <godoc>os/signal#NotifyContext</godoc> function.</figcaption>
</figure>

Consider <ref>signals.example</ref>. We use <godoc>os/signal#NotifyContext</godoc> to listen for `ctrl-c`. This function returns a wrapped <godoc>context#Context</godoc> that will cancel when the signal is received. It also returns a <godoc>context#CancelFunc</godoc> that can be used to cancel <godoc>context#Context</godoc> when needed.

<figure id="signals.example" type="listing">
<code src="src/signals/main.go#example"></code>
<figcaption>Listening for system signals.</figcaption>
</figure>

## Testing Signals

Testing system signals is tricky and care must be taken not to accidentally exit your running tests. Unfortunately, the <godoc>syscall</godoc> package does not provide a "test" signal, or a way to implement a test signal.

We can use <godoc>syscall#SIGUSR1</godoc> or <godoc>syscall#SIGUSR2</godoc> in our tests as these are allocated to the developer to use for their own purposes.

When we are testing signals, we are testing a **global** signal, that will caught by anyone else who is listening to that signal. Because of this we want to make that when testing signals we aren't running the tests in parallel and that we don't have other tests also listening to the same signal.

Consider <ref>testing.func</ref>. How do we test that the `Listener` function will respond properly to a signal? We don't want to make that the responsibility of the `Listener` function, it already has a <godoc>context#Context</godoc> that it can listen to for cancellation. The `Listener` function doesn't care why it was told to stop listening, it just needs to stop listening. This could be because we receive an interrupt signal, because a deadline has passed, or because the application no longer needs the `Lister` function to keep running.

<figure id="testing.func" type="listing">
<code src="src/testing/signals_test.go#func"></code>
<figcaption>The `Listener` function.</figcaption>
</figure>

In, <ref>testing.example</ref>, before we call the `Listener` function, we first create a <godoc>context#Context</godoc> that will self-cancel after 5 seconds if nothing else happens. We then wrap that <godoc>context#Context</godoc> with one received from the <godoc>os/signal#NotifyContext</godoc> function, that will self-cancel when the system receives a `TEST_SIGNAL` signal.

Our test blocks with a `select` waiting for either <godoc>context#Context</godoc> to be cancelled, and then respond accordingly.

<figure id="testing.example" type="listing">
<code src="src/testing/signals_test.go#example"></code>
<figcaption>Testing the `Listener` function.</figcaption>
</figure>

Inside the test, in a goroutine, we can trigger the `TEST_SIGNAL` signal by sending it to the current process, <godoc>syscall#Getpid</godoc>, with the <godoc>syscall#Kill</godoc> function.

<figure id="testing.kill" type="listing">
<go test="-v" src="src/testing" code="signals_test.go#kill"></go>
<figcaption>Sending a `TEST_SIGNAL` signal.</figcaption>
</figure>
