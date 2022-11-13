package feishu

import (
	"context"
	"fmt"

	"github.com/ajiany/pikachu/tools/httpx"
	"github.com/pkg/errors"
)

// GetUserGroups 查询用户组列表
// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/contact-v3/group/simplelist
func (cli *Client) GetUserGroups(ctx context.Context, pageToken string) ([]UserGroup, *Pagination, error) {
	query := map[string]interface{}{
		"page_token": pageToken,
		"page_size":  100,
	}
	r, err := cli.Api.GET(ctx, "/contact/v3/group/simplelist", &httpx.Option{Query: query})
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	var resp struct {
		Data struct {
			Pagination
			Items []UserGroup `test_helper:"grouplist"`
		} `test_helper:"data"`
	}
	if err := r.ParsedStruct(&resp); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return resp.Data.Items, &resp.Data.Pagination, nil
}

// AllUserGroups 获取所有用户组
func (cli *Client) AllUserGroups(ctx context.Context) ([]UserGroup, error) {
	var pageToken string
	groups := make([]UserGroup, 0)

	for {
		items, pag, err := cli.GetUserGroups(ctx, pageToken)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		groups = append(groups, items...)

		if len(items) == 0 || !pag.HasMore {
			break
		}
	}

	return groups, nil
}

// UserGroupAddMember 添加用户组成员
func (cli *Client) AddUserGroupMember(ctx context.Context, groupId, userId string) error {
	url := fmt.Sprintf("/contact/v3/group/%s/member/add", groupId)
	param := map[string]interface{}{
		"member_type":    "user",
		"member_id_type": "user_id",
		"member_id":      userId,
	}
	r, err := cli.Api.POST(ctx, url, &httpx.Option{Params: param})
	if err != nil {
		if r != nil && r.Resp.StatusCode == 400 {
			var resp ResponseData
			if err := r.ParsedStruct(&resp); err != nil {
				return errors.WithStack(err)
			}

			// 用户已加入用户组
			if resp.Code == 42005 {
				return nil
			}
		}
		return errors.WithStack(err)
	}

	return nil
}
