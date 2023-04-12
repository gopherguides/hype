# Timeouts and Deadlines

In addition to allowing us to manually cancel a <godoc>context#Context</godoc>, the <godoc>context</godoc> package also provides mechanisms for creating a <godoc>context#Context</godoc> that will self-cancel after, or at, a given time. Using these mechanics allows us to control how long to run some before we give up and assume that the operation has failed.

## Cancelling at a Specific Time

The <godoc>context</godoc> package provides two functions for creating time based, self-cancelling, a <godoc>context#Context</godoc>; <godoc>context#WithTimeout</godoc> and <godoc>context#WithDeadline</godoc>.

<figure id="withdeadline.doc" type="listing">
<go doc="context.WithDeadline"></go>
<figcaption>The <godoc>context#WithDeadline</godoc> function.</figcaption>
</figure>

When using <godoc>context#WithDeadline</godoc>, <ref>withdeadline.doc</ref>, we need to provide an **absolute** time at which the <godoc>context#Context</godoc> should be cancelled. That means we need an exact date/time we want this <godoc>context#Context</godoc> to be cancelled, for example `March 14, 2029 3:45pm`.

Consider <ref>deadline.example</ref>. In it, we create a new <godoc>time#Time</godoc> for `January 1, 2030 00:00:00` and use it to create a <godoc>context#Context</godoc> that will self-cancel at that date and time.

<figure id="deadline.example" type="listing">
<go src="src/with-deadline" run="." code="/main.go#example"></go>
<figcaption>Using <godoc>context#WithDeadline</godoc>.</figcaption>
</figure>

## Cancelling After a Duration

While being able to cancel a <godoc>context#Context</godoc> at a particular time is useful, more often than not we want to cancel a <godoc>context#Context</godoc> after a certain amount of time has passed.

<figure id="withtimeout.doc" type="listing">
<go doc="context.WithTimeout"></go>
<figcaption>The <godoc>context#WithTimeout</godoc> function.</figcaption>
</figure>

When using <godoc>context#WithTimeout</godoc>, <ref>withtimeout.doc</ref>a, we need to provide an **relative** <godoc>time#Duration</godoc> at which the <godoc>context#Context</godoc> should be cancelled.

Consider <ref>timeout.example</ref>. In it, we create a new self-cancelling <godoc>context#Context</godoc> that will self-cancel after 5 seconds using <godoc>context#WithTimeout</godoc>.

<figure id="timeout.example" type="listing">
<go src="src/with-timeout" run="." code="/main.go#example"></go>
<figcaption>Using <godoc>context#WithTimeout</godoc>.</figcaption>
</figure>

Functionally, we could have used <godoc>context#WithDeadline</godoc> instead, but <godoc>context#WithTimeout</godoc> is more convenient when we want to cancel a <godoc>context#Context</godoc> after a certain amount of time has passed.
