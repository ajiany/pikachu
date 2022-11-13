package feishu

import (
	"context"
	"fmt"

	"github.com/ajiany/pikachu/tools/httpx"
	"github.com/pkg/errors"
)

type SendMessageParam struct {
	ReceiveId string `test_helper:"receive_id"`
	Content   string `test_helper:"content"`
	MsgType   string `test_helper:"msg_type"`
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
	Content       interface{} `test_helper:"content,omitempty"`
	Card          interface{} `test_helper:"card,omitempty"`
	MsgType       string      `test_helper:"msg_type"`
	DepartmentIds []string    `test_helper:"department_ids,omitempty"`
	OpenIds       []string    `test_helper:"open_ids,omitempty"`
	UserIds       []string    `test_helper:"user_ids,omitempty"`
	UnionIds      []string    `test_helper:"union_ids,omitempty"`
}

func (cli *Client) BulkSendMessage(ctx context.Context, param BulkSendMessageParam) error {
	_, err := cli.Api.POST(ctx, "/message/v4/batch_send", &httpx.Option{Params: param})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
