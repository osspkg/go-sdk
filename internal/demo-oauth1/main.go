package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/deweppro/go-sdk/app"
	"github.com/deweppro/go-sdk/auth"
	"github.com/deweppro/go-sdk/log"
	"github.com/deweppro/go-sdk/webutil"
)

var (
	provConf = &auth.ConfigOAuth{
		Provider: []auth.ConfigOAuthItem{
			{
				Code:         "google",
				ClientID:     "****************.apps.googleusercontent.com",
				ClientSecret: "****************",
				RedirectURL:  "https://example.com/oauth/callback/google",
			},
		},
	}

	servConf = webutil.ConfigHttp{Addr: ":8080"}
)

func main() {
	ctx := app.NewContext()
	authServ := auth.NewOAuth(provConf)

	route := webutil.NewRouter()
	route.Route("/oauth/request/google", authServ.Request(auth.CodeGoogle), http.MethodGet)
	route.Route("/oauth/callback/google", authServ.CallBack(auth.CodeGoogle, oauthCallBackHandler), http.MethodGet)

	serv := webutil.NewServerHttp(servConf, route, log.Default())
	serv.Up(ctx) //nolint: errcheck
	<-time.After(60 * time.Minute)
	ctx.Close()
	serv.Down() //nolint: errcheck
}

const out = `
email: %s
name:  %s
ico:   %s
`

func oauthCallBackHandler(w http.ResponseWriter, _ *http.Request, u auth.UserOAuth) {
	w.WriteHeader(200)
	fmt.Fprintf(w, out, u.GetEmail(), u.GetName(), u.GetIcon())
}
