# DEPRECATED

Use [grpc.UnaryServerInterceptor](https://godoc.org/google.golang.org/grpc#UnaryServerInterceptor) and friends instead.

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
s.RegisterService(instr.Must("adder.Adder",
    (*adder.AdderServer)(nil),
    impl,
    func(sn, method string, took time.Duration, err error) {
	    log.Println(sn, method, took, err)
    }),
    impl)
```

instead of

```
...
    adder.RegisterAdderServer(s, impl)
```

## Limitations
Does not support streaming RPC (yet).

## Benchmarks

```
PASS
BenchmarkInstrumented-8	   20000	     95334 ns/op
BenchmarkDirect-8      	   20000	     97192 ns/op
ok  	github.com/sasha-s/grpc-instrument	5.807s
```

Adds about 2% overhead when both server and client are on the same machine the the method in question is almost NOP (adding two numbers). In reality the overhead should be negligible.
