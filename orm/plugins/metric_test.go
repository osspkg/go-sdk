package plugins

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMetric(t *testing.T) {
	w := &bytes.Buffer{}
	tl := &stdout{Writer: w}

	demo1 := NewMetric(nil)
	demo1.ExecutionTime("hello1", func() {})

	demo2 := NewMetric(tl)
	demo2.ExecutionTime("hello2", func() {})

	result := w.String()
	require.NotContains(t, result, "hello1")
	require.Contains(t, result, "hello2")
}
