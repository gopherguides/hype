# Context

<metadata>
class: center, middle, inverse
duration: 45m
exercises: 2 Shared Lab, 1 Student Labs
level: Intermediate
name: gRPC Middleware
topic: Distributed Computing
</metadata>

Introduced in Go 1.7, the <godoc>context</godoc> package was introduced to provide a cleaner way, than the use of channels, of managing cancellation and timeouts across goroutines.

While the scope, and API footprint of the package is pretty small, it was a welcome addition to the language when introduced.

The <godoc>context</godoc> package, <ref>context</ref>, defines the <godoc>context#Context</godoc> type, which carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.

Context is, mostly, used for controlling concurrent subsystems in your application. This week we will cover the different kinds of behavior with contexts including canceling, timeouts, and values. We'll also see how we can clean up a lot of code involving channels by using contexts.

<figure id="context" type="listing">
<go doc="-short context"></go>
<figcaption>The <godoc>context</godoc> package.</figcaption>
</figure>

<include src="basics/basics.md"></include>

<include src="rules/rules.md"></include>

<include src="nodes/nodes.md"></include>

<include src="values/values.md"></include>

<include src="cancellation/cancellation.md"></include>

<include src="timeouts/timeouts.md"></include>

<include src="errors/errors.md"></include>

<include src="signals/signals.md"></include>

# Summary

In this chapter we explore the concept of contexts in Go. We learn that contexts are a way to manage cancellation, timeouts, and other request-scoped values across API boundaries and between processes. We also learn how to use contexts to clean up a lot of code involving channels, such as listening for system signals. We discussed the nodal hierarchy of how the <godoc>context</godoc> package wraps a new <godoc>context#Context</godoc> around a parent <godoc>context#Context</godoc>. We learned the different was to cancel a <godoc>context#Context</godoc> and how to use multiple <godoc>context#Context</godoc>s to confirm shutdown behavior. The <godoc>context</godoc> package, while small, is a very powerful tool for managing concurrency in your application.
