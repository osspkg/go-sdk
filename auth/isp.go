package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/deweppro/go-sdk/errors"
	"github.com/deweppro/go-sdk/ioutil"
	"golang.org/x/oauth2"
)

var (
	errProviderFail = errors.New("provider not found")
)

type (
	UserOAuth interface {
		GetName() string
		GetEmail() string
		GetIcon() string
	}

	OAuthProvider interface {
		Code() string
		Config(conf ConfigOAuthItem)
		AuthCodeURL() string
		AuthCodeKey() string
		Exchange(ctx context.Context, code string) (UserOAuth, error)
	}
)

func (v *OAuth) AddProviders(p ...OAuthProvider) {
	v.mux.Lock()
	defer v.mux.Unlock()

	for _, item := range p {
		for _, cp := range v.config.Provider {
			if cp.Code == item.Code() {
				item.Config(cp)
				v.list[item.Code()] = item
			}
		}
	}
}

func (v *OAuth) GetProvider(name string) (OAuthProvider, error) {
	v.mux.RLock()
	defer v.mux.RUnlock()

	p, ok := v.list[name]
	if !ok {
		return nil, errProviderFail
	}
	return p, nil
}

/**********************************************************************************************************************/

type oauth2Config interface {
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	Client(ctx context.Context, t *oauth2.Token) *http.Client
}

func oauth2ExchangeContext(
	ctx context.Context, code string, uri string, srv oauth2Config, model json.Unmarshaler,
) error {
	tok, err := srv.Exchange(ctx, code)
	if err != nil {
		return errors.Wrapf(err, "exchange to oauth service")
	}
	client := srv.Client(ctx, tok)
	resp, err := client.Get(uri)
	if err != nil {
		return errors.Wrapf(err, "client request to oauth service")
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "read response from oauth service")
	}
	if err = json.Unmarshal(b, model); err != nil {
		return errors.Wrapf(err, "decode oauth model")
	}
	return nil
}
