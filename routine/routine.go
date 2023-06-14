package routine

import (
	"context"
	"sync"
	"time"

	"github.com/osspkg/go-sdk/errors"
)

func Interval(ctx context.Context, interval time.Duration, call func(context.Context)) {
	call(ctx)

	go func() {
		tick := time.NewTicker(interval)
		defer tick.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				call(ctx)
			}
		}
	}()
}

func Retry(count int, ttl time.Duration, call func() error) error {
	var err error
	for i := 0; i < count; i++ {
		if e := call(); e != nil {
			err = errors.Wrap(err, errors.Wrapf(e, "[#%d]", i))
			time.Sleep(ttl)
			continue
		}
		return nil
	}
	return errors.Wrapf(err, "retry error")
}

func Parallel(calls ...func()) {
	var wg sync.WaitGroup
	for _, call := range calls {
		call := call
		wg.Add(1)
		go func() {
			call()
			wg.Done()
		}()
	}
	wg.Wait()
}
