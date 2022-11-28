package main

import (
	"context"
	"io"
)

type readerCtx struct {
	ctx        context.Context
	r          io.Reader
	total      uint64
	updateFunc func(total uint64)
}

func (r *readerCtx) Read(p []byte) (n int, err error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	n, err = r.r.Read(p)
	if err != nil {
		return n, err
	}
	r.total += uint64(n)
	r.updateFunc(r.total)
	return n, nil
}

// NewReader gets a context-aware io.Reader.
func NewReaderCtx(ctx context.Context, r io.Reader, updateFunc func(total uint64)) io.Reader {
	return &readerCtx{
		ctx:        ctx,
		r:          r,
		updateFunc: updateFunc,
	}
}
