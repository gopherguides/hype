# Context Rules

According the <godoc>context</godoc> documentation there are rules that must be followed when using the <godoc>context</godoc> package, <ref>doc</ref>.

<figure id="doc" type="listing">

- Programs that use Contexts should follow these rules to keep interfaces consistent across packages and enable static analysis tools to check context propagation.
- Do not store Contexts inside a struct type; instead, pass a Context explicitly to each function that needs it. The Context should be the first parameter, typically named ctx.
- Do not pass a `nil` Context, even if a function permits it. Pass <godoc>context#TODO</godoc> if you are unsure about which Context to use.
- Use context Values only for request-scoped data that transits processes and APIs, not for passing optional parameters to functions.
- The same Context may be passed to functions running in different goroutines; Contexts are safe for simultaneous use by multiple goroutines.

<figcaption>Rules for using contexts.</figcaption>

</figure>
