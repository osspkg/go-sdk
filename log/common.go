package log

import "io"

// nolint: golint
const (
	levelFatal uint32 = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

var levels = map[uint32]string{
	levelFatal: "FAT",
	LevelError: "ERR",
	LevelWarn:  "WRN",
	LevelInfo:  "INF",
	LevelDebug: "DBG",
}

type Fields map[string]interface{}

type Sender interface {
	PutEntity(v *entity)
	SendMessage(level uint32, call func(v *message))
	Close()
}

// Writer writer interface
type Writer interface {
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

// Logger base interface
type Logger interface {
	SetOutput(out io.Writer)
	SetLevel(v uint32)
	GetLevel() uint32
	Close()

	WithFields(Fields) Writer
	Writer
}
