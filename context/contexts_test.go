package context_test

import (
	ccc "context"
	"fmt"
	"testing"
	"time"

	"github.com/deweppro/go-sdk/context"
)

func TestUnit_Combine(t *testing.T) {
	cc, cancel := ccc.WithCancel(ccc.Background())
	c := context.Combine(ccc.Background(), ccc.Background(), cc)
	if c == nil {
		t.Fatalf("contexts.Combine returned nil")
	}

	select {
	case <-c.Done():
		t.Fatalf("<-c.Done() == it should block")
	default:
	}

	cancel()
	<-time.After(time.Second)

	select {
	case <-c.Done():
	default:
		t.Fatalf("<-c.Done() it shouldn't block")
	}

	if got, want := fmt.Sprint(c), "context.Background.WithCancel"; got != want {
		t.Fatalf("contexts.Combine() = %q want %q", got, want)
	}
}
