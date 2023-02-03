package context

import (
	cc "context"
	"reflect"
)

func Combine(multi ...cc.Context) cc.Context {
	ctx, cancel := cc.WithCancel(cc.Background())

	go func() {
		cases := make([]reflect.SelectCase, 0, len(multi))
		for _, vv := range multi {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(vv.Done()),
			})
		}
		chosen, _, _ := reflect.Select(cases)
		switch chosen {
		default:
			cancel()
		}
	}()

	return ctx
}
