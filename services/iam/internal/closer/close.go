package closer

import (
	"context"
	"errors"
	"slices"
	"sync"
)

type CloserFunc func() error

type Cleanup struct {
	mutex   sync.Mutex
	closers []CloserFunc
}

func (c *Cleanup) Add(closerFn CloserFunc) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.closers = append(c.closers, closerFn)
}

func (c *Cleanup) Close(ctx context.Context) error {
	var errs []error
	closeCh := make(chan struct{}, 1)

	c.mutex.Lock()
	closers := c.closers
	c.closers = nil
	c.mutex.Unlock()

	go func() {
		var wg sync.WaitGroup
		for _, closer := range slices.Backward(closers) {
			wg.Go(func() {
				if err := closer(); err != nil {
					errs = append(errs, err)
				}
			})
		}
		wg.Wait()
		closeCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-closeCh:
		return errors.Join(errs...)
	}
}
