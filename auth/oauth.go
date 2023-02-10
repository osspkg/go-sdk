package auth

import (
	"net/http"
	"sync"
)

/**********************************************************************************************************************/

type (
	ConfigOAuthItem struct {
		Code         string `yaml:"code"`
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
		RedirectURL  string `yaml:"redirect_url"`
	}

	ConfigOAuth struct {
		Provider []ConfigOAuthItem `yaml:"oauth"`
	}

	configIsp struct {
		State       string
		AuthCodeKey string
		RequestURL  string
	}
)

/**********************************************************************************************************************/

type (
	OAuth struct {
		config *ConfigOAuth
		list   map[string]OAuthProvider
		mux    sync.RWMutex
	}

	CallBackOAuth func(http.ResponseWriter, *http.Request, UserOAuth)
)

func NewOAuth(c *ConfigOAuth) *OAuth {
	return &OAuth{
		config: c,
		list:   make(map[string]OAuthProvider),
	}
}

func (v *OAuth) Up() error {
	v.AddProviders(
		&IspYandex{},
	)
	return nil
}

func (v *OAuth) Down() error {
	return nil
}

func (v *OAuth) Request(name string) func(http.ResponseWriter, *http.Request) {
	p, err := v.GetProvider(name)
	if err != nil {
		return func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error())) //nolint: errcheck
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, p.AuthCodeURL(), http.StatusMovedPermanently)
	}
}

func (v *OAuth) CallBack(name string, call CallBackOAuth) func(w http.ResponseWriter, r *http.Request) {
	p, err := v.GetProvider(name)
	if err != nil {
		return func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error())) //nolint: errcheck
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get(p.AuthCodeKey())
		u, err := p.Exchange(r.Context(), code)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error())) //nolint: errcheck
			return
		}
		call(w, r, u)
	}
}
