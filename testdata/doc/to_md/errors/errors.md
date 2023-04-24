# Context Errors

In a complex system, or even in a small one, when a <godoc>context#Context</godoc> is cancelled, we need a way to know what caused the cancellation. It is possible that <godoc>context#Context</godoc> was cancelled by a <godoc>context#CancelFunc</godoc> successfully, if it was cancelled because it timed out, or some other reason.

The <godoc>context#Context.Err</godoc> method, <ref>context.err</ref> returns the error that caused the context to be cancelled.

<figure id="context.err" type="listing">

<go doc="context.Context.Err"></go>

<figcaption>The <godoc>context#Context.Err</godoc> method.</figcaption>
</figure>

## Context Cancelled Error

The <godoc>context</godoc> package defines two different <godoc>builtin#error</godoc> variables that can be used to check an <godoc>builtin#error</godoc> that was returned from <godoc>context#Context.Err</godoc> method.

The first is <godoc>context#Canceled</godoc>, <ref>canceled.doc</ref>, which is returned when the context is cancelled through the use of a <godoc>context#CancelFunc</godoc> function. This <godoc>builtin#error</godoc> is considered to indicate a "successful" cancellation.

<figure id="canceled.doc" type="listing">

<go doc="context.Canceled"></go>

<figcaption>The <godoc>context#Canceled</godoc> <godoc>builtin#error</godoc>.</figcaption>
</figure>

Consider <ref>canceled.example</ref>. When we first check the <godoc>context#Context.Err</godoc> method, it returns `nil`. After we call the <godoc>context#CancelFunc</godoc> function provided by <godoc>context#WithCancel</godoc>, the <godoc>context#Context.Err</godoc> method returns a <godoc>context#Canceled</godoc> error.

<figure id="canceled.example" type="listing">

<go src="src/canceled" run="." code="main.go#example"></go>

<figcaption>Checking for cancellation errors.</figcaption>
</figure>

As we can see from the output in <ref>canceled.example</ref>, repeated calls to the <godoc>context#Context.Err</godoc> method return the same <godoc>context#Canceled</godoc> error.

## Context Deadline Exceeded Error

When a <godoc>context#Context</godoc> is cancelled due to a deadline, or timeout, being exceeded, the <godoc>context#Context.Err</godoc> method returns a <godoc>context#DeadlineExceeded</godoc> error, <ref>exceeded.doc</ref>.

<figure id="exceeded.doc" type="listing">

<go doc="context.DeadlineExceeded"></go>

<figcaption>The <godoc>context#DeadlineExceeded</godoc> <godoc>builtin#error</godoc>.</figcaption>
</figure>

Consider <ref>exceeded.example</ref>. We create a <godoc>context#Context</godoc> that will self cancel after 1 second. When we check <godoc>context#Context.Err</godoc> method, before the <godoc>context#Context</godoc> times out, it returns `nil`.

<figure id="exceeded.example" type="listing">

<go src="src/deadline" run="." code="main.go#example"></go>

<figcaption>Checking for deadline exceeded errors.</figcaption>
</figure>

As we can see from the output, the <godoc>context#Context</godoc> times out after the specified time, and the <godoc>context#Context.Err</godoc> method returns a <godoc>context#DeadlineExceeded</godoc> error.
