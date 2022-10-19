package feishu

import (
	"context"
	"fmt"

	"github.com/ajiany/pikachu/tools/httpx"
	"github.com/pkg/errors"
)

type SendMessageParam struct {
	ReceiveId string `json:"receive_id"`
	Content   string `json:"content"`
	MsgType   string `json:"msg_type"`
}

func (cli *Client) SendMessage(ctx context.Context, recvType string, param SendMessageParam) error {
	url := fmt.Sprintf("/im/v1/messages?receive_id_type=%s", recvType)
	_, err := cli.Api.POST(ctx, url, &httpx.Option{Params: param})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type BulkSendMessageParam struct {
	Content       interface{} `json:"content,omitempty"`
	Card          interface{} `json:"card,omitempty"`
	MsgType       string      `json:"msg_type"`
	DepartmentIds []string    `json:"department_ids,omitempty"`
	OpenIds       []string    `json:"open_ids,omitempty"`
	UserIds       []string    `json:"user_ids,omitempty"`
	UnionIds      []string    `json:"union_ids,omitempty"`
}

func (cli *Client) BulkSendMessage(ctx context.Context, param BulkSendMessageParam) error {
	_, err := cli.Api.POST(ctx, "/message/v4/batch_send", &httpx.Option{Params: param})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
