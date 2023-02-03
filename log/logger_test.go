package log_test

import (
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/deweppro/go-sdk/log"
	"github.com/stretchr/testify/require"
)

func TestUnit_New(t *testing.T) {
	require.NotNil(t, log.Default())

	filename, err := os.CreateTemp(os.TempDir(), "test_new_default-*.log")
	require.NoError(t, err)

	log.SetOutput(filename)
	log.SetLevel(log.LevelDebug)
	require.Equal(t, log.LevelDebug, log.GetLevel())

	go log.Infof("async %d", 1)
	go log.Warnf("async %d", 2)
	go log.Errorf("async %d", 3)
	go log.Debugf("async %d", 4)
	log.Infof("sync %d", 1)
	log.Warnf("sync %d", 2)
	log.Errorf("sync %d", 3)
	log.Debugf("sync %d", 4)
	log.WithFields(log.Fields{"ip": "0.0.0.0"}).Infof("context1")
	log.WithFields(log.Fields{"nil": nil}).Infof("context2")
	log.WithFields(log.Fields{"func": func() {}}).Infof("context3")

	<-time.After(time.Second * 1)
	log.Close()

	require.NoError(t, filename.Close())
	data, err := os.ReadFile(filename.Name())
	require.NoError(t, err)
	require.NoError(t, os.Remove(filename.Name()))

	sdata := string(data)
	require.Contains(t, sdata, `"lvl":"INF","msg":"async 1"`)
	require.Contains(t, sdata, `"lvl":"WRN","msg":"async 2"`)
	require.Contains(t, sdata, `"lvl":"ERR","msg":"async 3"`)
	require.Contains(t, sdata, `"lvl":"DBG","msg":"async 4"`)
	require.Contains(t, sdata, `"lvl":"INF","msg":"sync 1"`)
	require.Contains(t, sdata, `"lvl":"WRN","msg":"sync 2"`)
	require.Contains(t, sdata, `"lvl":"ERR","msg":"sync 3"`)
	require.Contains(t, sdata, `"msg":"context1","ctx":{"ip":"0.0.0.0"}`)
	require.Contains(t, sdata, `"msg":"context2","ctx":{"nil":null}`)
	require.Contains(t, sdata, `"msg":"context3","ctx":{"func":"unsupported field value: (func())`)
}

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()

	ll := log.New()
	ll.SetOutput(io.Discard)
	ll.SetLevel(log.LevelDebug)
	wg := sync.WaitGroup{}

	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		wg.Add(1)
		for p.Next() {
			ll.WithFields(log.Fields{"a": "b"}).Infof("hello")
		}
		wg.Done()
	})
	wg.Wait()
	ll.Close()
}
