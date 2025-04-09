package utils

import "context"

func NewWithCancel(parent context.Context) (context.Context, context.CancelFunc) {
	closeCh := make(chan struct{})
	doneCh := make(chan struct{})

	ctx, _ := context.WithCancel(parent)

	go func() {
		defer close(doneCh)
		for {
			select {
			case <-ctx.Done():
				return
			case <-closeCh:
				return
			}
		}
	}()

	cancelFunc := func() {
		close(closeCh)
		<-doneCh
	}

	return ctx, cancelFunc
}
