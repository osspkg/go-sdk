package plugins

import (
	"github.com/deweppro/go-sdk/log"
)

var (
	//StdOutLog simple stdout debug log
	StdOutLog = func() log.Logger {
		l := log.Default()
		l.SetLevel(log.LevelDebug)
		l.SetOutput(StdOutWriter)
		return l
	}()
)
