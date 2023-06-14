package routine_test

import (
	"fmt"
	"testing"

	"github.com/osspkg/go-sdk/routine"
)

func TestUnit_Parallel(t *testing.T) {
	routine.Parallel(
		func() {
			fmt.Println("a")
		}, func() {
			fmt.Println("b")
		}, func() {
			fmt.Println("c")
		},
	)
}
