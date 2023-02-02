package syscall

import (
	"os"
	"os/signal"
	"strconv"
	s "syscall"
)

// OnStop calling a function if you send a system event stop
func OnStop(callFunc func()) {
	quit := make(chan os.Signal, 4)
	signal.Notify(quit, os.Interrupt, s.SIGINT, s.SIGTERM, s.SIGKILL)
	<-quit

	callFunc()
}

// OnUp calling a function if you send a system event SIGHUP
func OnUp(callFunc func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, s.SIGHUP)
	<-quit

	callFunc()
}

// OnCustom calling a function if you send a system custom event
func OnCustom(callFunc func(), sig ...os.Signal) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, sig...)
	<-quit

	callFunc()
}

// Pid write pid file
func Pid(filename string) error {
	pid := strconv.Itoa(s.Getpid())
	return os.WriteFile(filename, []byte(pid), 0755)
}
