package feishu

import (
	"context"

	"github.com/ajiany/pikachu/tools/httpx"
	"github.com/pkg/errors"
)

func (cli *Client) GetChats(ctx context.Context, pageToken string) ([]Chat, *Pagination, error) {
	query := map[string]interface{}{
		// "user_id_type": "user_id",
		"page_token": pageToken,
		"page_size":  100,
	}
	r, err := cli.Api.GET(ctx, "/im/v1/chats", &httpx.Option{Query: query})
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	var resp struct {
		Data struct {
			Pagination
			Items []Chat `json:"items"`
		} `json:"data"`
	}
	if err := r.ParsedStruct(&resp); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return resp.Data.Items, &resp.Data.Pagination, nil
}

func (cli *Client) AllChats(ctx context.Context) ([]Chat, error) {
	var pageToken string
	chats := make([]Chat, 0)

	for {
		items, pag, err := cli.GetChats(ctx, pageToken)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		chats = append(chats, items...)

		if len(items) == 0 || !pag.HasMore {
			break
		}
	}

	return chats, nil
}
