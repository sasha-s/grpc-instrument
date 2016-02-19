package adder

import (
	"time"

	"golang.org/x/net/context"
)

type Impl struct {
	Delay time.Duration
}

func (a Impl) Add(ctx context.Context, req *AddRequest) (*AddReply, error) {
	time.Sleep(a.Delay)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return &AddReply{req.A + req.B}, nil
}

func (a Impl) Add2(ctx context.Context, req *Add2Request) (*Add2Reply, error) {
	time.Sleep(a.Delay)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return &Add2Reply{req.A ^ req.B}, nil
}
