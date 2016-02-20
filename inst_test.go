package instr

import (
	"log"
	"net"
	"testing"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/sasha-s/grpc-instrument/adder"
)

func BenchmarkInstrumented(b *testing.B) {
	client, done := client(true)
	defer done()
	bench(b, client)
}

func BenchmarkDirect(b *testing.B) {
	client, done := client(false)
	defer done()
	bench(b, client)
}

func TestInstrument(t *testing.T) {
	tcs := []struct {
		method         string
		a, b           int32
		sleepFor       time.Duration
		expectedError  error
		expectedResult int32
	}{
		{"Add", 1, 2, 0, nil, 3},
		{"Add", -1, 1, 0, nil, 0},
		{"Add", -1, -2, time.Millisecond * 1, nil, -3},
		{"Add", -1, 1, time.Millisecond * 10, context.DeadlineExceeded, 0},
		{"Add2", 1 + 2 + 4, 2, time.Millisecond * 1, nil, 5},
	}

	for _, tc := range tcs {
		grpcServer := grpc.NewServer()
		impl := adder.Impl{tc.sleepFor}
		type r struct {
			service string
			method  string
			took    time.Duration
			err     error
		}
		rC := make(chan r, 1)
		d, err := ServiceDesc("adder.Adder", (*adder.AdderServer)(nil), impl, func(sn, n string, took time.Duration, err error) {
			rC <- r{sn, n, took, err}
		})
		if err != nil {
			panic(err)
		}
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			panic(err)
		}
		// Make sure the server exits. This happens only when we return from the TestInstrument function (after all test cases are processed) but it'd fine.
		defer l.Close()
		grpcServer.RegisterService(d, impl)
		go func() {
			log.Println("grpc server done:", grpcServer.Serve(l))
		}()

		conn, err := grpc.Dial(l.Addr().String(), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second))
		if err != nil {
			panic(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*5)
		defer cancel()
		var result int32
		var gotErr error
		switch tc.method {
		case "Add":
			var resp *adder.AddReply
			resp, gotErr = adder.NewAdderClient(conn).Add(ctx, &adder.AddRequest{tc.a, tc.b})
			if gotErr == nil {
				result = resp.R
			}
		case "Add2":
			var resp *adder.Add2Reply
			resp, gotErr = adder.NewAdderClient(conn).Add2(ctx, &adder.Add2Request{tc.a, tc.b})
			if gotErr == nil {
				result = resp.R
			}
		}
		// Get instrumentation results.
		ir := <-rC
		switch {
		case tc.expectedResult != result:
			log.Fatalf("expected result %d got %d", tc.expectedResult, result)
		case desc(tc.expectedError) != desc(gotErr):
			log.Fatalf("expected error %v got %v", tc.expectedError, gotErr)
		case tc.sleepFor >= ir.took:
			log.Fatalf("expected the call to take at least %v, took %v", tc.sleepFor, ir.took)
		case tc.sleepFor+time.Millisecond*4 < ir.took:
			log.Fatalf("expected the call to take at no longer than %v+4ms, took %v", tc.sleepFor, ir.took)
		case tc.expectedError != ir.err:
			log.Fatalf("expected error %v got %v", tc.expectedError, ir.err)
		case tc.method != ir.method:
			log.Fatalf("expected method %s got %s", tc.method, ir.method)
		}
	}
}

func bench(b *testing.B, client adder.AdderClient) {
	b.ResetTimer()
	req := &adder.AddRequest{1, 2}
	ctx := context.TODO()
	for i := 0; i < b.N; i++ {
		c, cancel := context.WithTimeout(ctx, time.Millisecond*20)
		client.Add(c, req)
		cancel()
	}
}

func client(shouldInstr bool) (client adder.AdderClient, done func() error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	s := grpc.NewServer()
	impl := adder.Impl{}
	if shouldInstr {
		// Add NOP instrumentation.
		s.RegisterService(Must("adder.Adder", (*adder.AdderServer)(nil), impl, func(sn, n string, took time.Duration, err error) {}), impl)
	} else {
		adder.RegisterAdderServer(s, impl)
	}
	go func() {
		s.Serve(l)
	}()
	conn, err := grpc.Dial(l.Addr().String(), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second))
	if err != nil {
		panic(err)
	}
	return adder.NewAdderClient(conn), l.Close
}

func desc(err error) interface{} {
	if err == nil {
		return nil
	}
	return grpc.ErrorDesc(err)
}
