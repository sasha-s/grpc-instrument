# Server-side Instrumentation hooks for golang GRPC.

## Why?
Latency monitoring.
Also see [grpc/grpc-go: Instrumentation hooks issue](https://github.com/grpc/grpc-go/issues/240)

## How?
Reflection.

## Usage
```go
s := grpc.NewServer()
impl := adder.Impl{}
s.RegisterService(Must("adder.Adder",
    (*adder.AdderServer)(nil),
    impl,
    func(sn, method string, took time.Duration, err error) {
	    log.Println(sn, method, took, err)
    }),
    impl)```

instead of

```
...
		adder.RegisterAdderServer(s, impl)
```


