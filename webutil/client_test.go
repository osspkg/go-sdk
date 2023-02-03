package webutil_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/deweppro/go-sdk/webutil"
	"github.com/stretchr/testify/require"
)

type (
	TestModel struct {
		Val struct {
			Page struct {
				Name string `json:"name"`
			} `json:"page"`
		}
	}
)

func (v *TestModel) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &v.Val)
}

func TestUnit_NewClientHttp_JSON(t *testing.T) {
	model := TestModel{}
	cli := webutil.NewClientHttp()
	err := cli.Call(context.TODO(), http.MethodGet, "https://www.githubstatus.com/api/v2/status.json", nil, &model)
	require.NoError(t, err)
	require.Equal(t, "GitHub", model.Val.Page.Name)
}

func TestUnit_NewClientHttp_Bytes(t *testing.T) {
	var model []byte
	cli := webutil.NewClientHttp()
	err := cli.Call(context.TODO(), http.MethodGet, "https://www.githubstatus.com/api/v2/status.json", nil, &model)
	require.NoError(t, err)
	require.Contains(t, string(model), ",\"name\":\"GitHub\",")
}
