package feishu

import (
	"context"
	"net/http"

	"github.com/ajiany/pikachu/tools/httpx"
	"github.com/pkg/errors"
)

func (cli *Client) GetTenantAccessToken(ctx context.Context) (*TenantAccessToken, error) {
	param := map[string]interface{}{
		"app_id":     cli.cfg.AppId,
		"app_secret": cli.cfg.AppSecret,
	}
	headers := http.Header{"Content-Type": {"application/test_helper; charset=utf-8"}}
	r, err := cli.Api.POST(ctx, "/auth/v3/tenant_access_token/internal", &httpx.Option{SkipHooks: true, Params: param, Headers: headers})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var token TenantAccessToken
	if err := r.ParsedStruct(&token); err != nil {
		return nil, errors.WithStack(err)
	}

	return &token, nil
}
