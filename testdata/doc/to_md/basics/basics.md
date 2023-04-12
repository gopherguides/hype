# The Context Interface

The <godoc>context#Context</godoc> interface, <ref>context.doc</ref>, consists of four methods. These methods provide us the ability to listen for cancellation and timeout events, retrieve values from the context hierarchy, and finally, a way to check what `error`, if any, caused the context to be canceled.

<figure id="context.doc" type="listing">

```godoc
type Context interface {
  Deadline() (deadline time.Time, ok bool)
  Done() <-chan struct{}
  Err() error
  Value(key interface{}) interface{}
}
```

<figcaption>The <godoc>context#Context</godoc> interface.</figcaption>
</figure>

We can see, in <ref>context.doc</ref>, that the <godoc>context#Context</godoc> interface implements several of the channel patterns we have already seen, such as have a `Done` channel that can be listened to for cancellation.

We will cover each of these methods in more detail later. For now, let's briefly look at each one of them.

## Context#Deadline

The <godoc>context#Context.Deadline</godoc> method, <ref>deadline.doc</ref>, can be used to check if a context has a cancellation deadline set, and if so, what that deadline is.

<figure id="deadline.doc" type="listing">
<go doc="context.Context.Deadline"></go>
<figcaption>The <godoc>context#Context.Deadline</godoc> method.</figcaption>
</figure>

## Context#Done

The <godoc>context#Context.Done</godoc> method, <ref>done.doc</ref>, can be used to listen for cancellation events. This is similar to how we can listen for a channel being closed, but it is more flexible.

<figure id="done.doc" type="listing">
<go doc="context.Context.Done"></go>
<figcaption>The <godoc>context#Context.Done</godoc> method.</figcaption>
</figure>

## Context#Err

The <godoc>context#Context.Err</godoc> method, <ref>err.doc</ref>, can be used to check if a context has been canceled.

<figure id="err.doc" type="listing">
<go doc="context.Context.Err"></go>
<figcaption>The <godoc>context#Context.Err</godoc> method.</figcaption>
</figure>

## Context#Value

The <godoc>context#Context.Value</godoc> method, <ref>value.doc</ref>, can be used to retrieve values from the context hierarchy.

<figure id="value.doc" type="listing">
<go doc="context.Context.Value"></go>
<figcaption>The <godoc>context#Context.Value</godoc> method.</figcaption>
</figure>

## Helper Functions

As we will see the <godoc>context</godoc> package provides a number of useful helper functions for wrapping a <godoc>context#Context</godoc> making the need for custom implementations of the <godoc>context#Context</godoc> interface less common.

<figure id="helper.doc" type="listing">
<go doc="-short context"></go>
<figcaption>The <godoc>context</godoc> package.</figcaption>
</figure>

## The Background Context

While often we might be given a <godoc>context#Context</godoc>, we might also be the one start a <godoc>context#Context</godoc>. The most common way to provide a quick and easy way to start a <godoc>context#Context</godoc> is to use the <godoc>context#Background</godoc> function, <ref>background.doc</ref>.

<figure id="background.doc" type="listing">
<go doc="context.Background"></go>
<figcaption>The <godoc>context#Background</godoc> function.</figcaption>
</figure>

In, <ref>empty</ref>, we print the <godoc>context#Context</godoc> returned by <godoc>context#Background</godoc>. As we can see from the output, the context is empty.

<figure id="empty" type="listing">
<go run="main.go" src="src/background/empty" code="main.go#example"></go>
<figcaption>The <godoc>context#Background</godoc> function.</figcaption>
</figure>

## Default Implementations

The <godoc>context#Background</godoc> interface, while empty, does provide default implementations of the <godoc>context#Context</godoc> interface. Because of this the <godoc>context#Background</godoc> context is almost always used as the base of a new <godoc>context#Context</godoc> hierarchy.

<figure id="implementation" type="listing">
<go run="main.go" src="src/background/implementation" code="main.go#example"></go>
<figcaption>The <godoc>context#Background</godoc> function provides default implementation of the <godoc>context#Context</godoc> interface.</figcaption>
</figure>
