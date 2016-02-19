package instr

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"google.golang.org/grpc"

	"golang.org/x/net/context"
)

// Returns a grpc service descriptor, injecting instrumentation.
// Panics on errors.
// Sample usage:
// s := grpc.NewServer()
// s.RegisterService(Must("adder.Adder",
//     (*adder.AdderServer)(nil),
//     impl,
//     func(sn, method string, took time.Duration, err error) {
//         log.Println(sn, method, took, err)
//     }),
//     impl)
//
func Must(
	// Fully qualified service name. Something like "adder.Adder"
	serviceName string,
	// The pointer to the service interface. Used to check whether the user
	// provided implementation satisfies the interface requirements.
	// Something like (*adder.AdderServer)(nil)
	inter interface{},
	// Service implementation.
	impl interface{},
	// Instrumentation function. Will be called on every method invocation, successful or not.
	instrument func(service, method string, took time.Duration, err error)) *grpc.ServiceDesc {
	d, err := ServiceDesc(serviceName, inter, impl, instrument)
	if err != nil {
		log.Panic(err)
	}
	return d
}

// Returns a grpc service descriptor, injecting instrumentation.
func ServiceDesc(
	// Fully qualified service name. Something like "adder.Adder"
	serviceName string,
	// The pointer to the service interface. Used to check whether the user
	// provided implementation satisfies the interface requirements.
	// Something like (*adder.AdderServer)(nil)
	inter interface{},
	// Service implementation.
	impl interface{},
	// Instrumentation function. Will be called on every method invocation, successful or not.
	instrument func(service, method string, took time.Duration, err error)) (*grpc.ServiceDesc, error) {
	implType := reflect.TypeOf(impl)
	interType := reflect.TypeOf(inter).Elem()
	if !implType.Implements(interType) {
		return nil, fmt.Errorf("expected impl %s to implement %s", implType.String(), interType.String())
	}
	d := grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: inter,
		Streams:     []grpc.StreamDesc{},
	}
	for i := 0; i < implType.NumMethod(); i++ {
		method := implType.Method(i)
		if method.Type.NumIn() != 3 || method.Type.NumOut() != 2 {
			return nil, fmt.Errorf("unexpected function signature %s", method.Type)
		}
		recType, ctxType, reqType := method.Type.In(0), method.Type.In(1), method.Type.In(2)
		respType, errType := method.Type.Out(0), method.Type.Out(1)
		_, _ = recType, respType

		// Sanity checks.
		switch {
		case !ctxType.Implements(reflect.TypeOf((*context.Context)(nil)).Elem()):
			return nil, fmt.Errorf("expected context.Context, got %v", ctxType)
		case !errType.Implements(reflect.TypeOf((*error)(nil)).Elem()):
			return nil, fmt.Errorf("expected error got %v", errType)
		}

		handler := func(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
			in := reflect.New(reqType.Elem())
			if err := dec(in.Interface()); err != nil {
				return nil, err
			}
			now := time.Now()
			res := method.Func.Call([]reflect.Value{reflect.ValueOf(srv), reflect.ValueOf(ctx), in})
			took := time.Since(now)
			if len(res) != 2 {
				panic("wtf")
			}
			var err error
			out, e := res[0].Interface(), res[1].Interface()
			if e != nil {
				err = e.(error)
			}
			instrument(serviceName, method.Name, took, err)
			if err != nil {
				return nil, err
			}
			return out, nil
		}
		d.Methods = append(d.Methods,
			grpc.MethodDesc{MethodName: method.Name, Handler: handler})
	}
	return &d, nil
}
