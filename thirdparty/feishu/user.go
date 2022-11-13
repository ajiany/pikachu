package feishu

import (
	"context"
	"fmt"

	"github.com/ajiany/pikachu/tools/httpx"
	"github.com/pkg/errors"
)

func (cli *Client) GetUserDetail(ctx context.Context, userId, userIdType string) (*User, error) {
	query := map[string]interface{}{
		"user_id_type": "user_id",
	}
	r, err := cli.Api.GET(ctx, fmt.Sprintf("/contact/v3/users/%s", userId), &httpx.Option{Query: query})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resp struct {
		Data struct {
			User User `test_helper:"user"`
		} `test_helper:"data"`
	}
	if err := r.ParsedStruct(&resp); err != nil {
		return nil, errors.WithStack(err)
	}

	return &resp.Data.User, nil
}

func (cli *Client) GetUserByEmailOrMobile(ctx context.Context, emails, mobiles []string) ([]UserIdItem, error) {
	query := map[string]interface{}{
		"user_id_type": "user_id",
	}
	param := map[string]interface{}{
		"emails":  emails,
		"mobiles": mobiles,
	}
	r, err := cli.Api.POST(ctx, "/contact/v3/users/batch_get_id", &httpx.Option{Query: query, Params: param})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resp struct {
		Data struct {
			UserList []UserIdItem `test_helper:"user_list"`
		} `test_helper:"data"`
	}
	if err := r.ParsedStruct(&resp); err != nil {
		return nil, errors.WithStack(err)
	}

	return resp.Data.UserList, nil
}
