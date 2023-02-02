package shell_test

import (
	"context"
	"testing"

	"github.com/deweppro/go-sdk/shell"
)

func TestUnit_ShellCall(t *testing.T) {
	sh := shell.New()
	sh.SetDir("/tmp")
	sh.SetEnv("LANG", "en_US.UTF-8")

	out, err := sh.Call(context.TODO(), "ls -la /tmp")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(string(out))
}
