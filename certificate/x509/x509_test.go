package x509_test

import (
	"testing"
	"time"

	"github.com/osspkg/go-sdk/certificate/x509"
)

func TestUnit_X509(t *testing.T) {
	conf := &x509.Config{
		Organization: "Demo Inc.",
	}

	crt, err := x509.NewCertCA(conf, time.Hour*24*365*10, "Demo Root R1")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(string(crt.Private), string(crt.Public))

	crt, err = x509.NewCert(conf, time.Hour*24*90, 2, crt, "example.com", "*.example.com")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(string(crt.Private), string(crt.Public))
}
