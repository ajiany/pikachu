package feishu

import (
	"context"
	"time"

	"github.com/ajiany/pikachu/tools/httpx"
	"github.com/pkg/errors"
)

// GetApprovalInstanceIds 批量获取审批实例ID
func (cli *Client) GetApprovalInstanceIds(ctx context.Context, approvalCode string, startTime, endTime time.Time, page, limit int64) ([]string, error) {
	param := map[string]interface{}{
		"approval_code": approvalCode,
		"start_time":    startTime.Unix() * 1000, // 毫秒
		"end_time":      endTime.Unix() * 1000,
		"offset":        (page - 1) * limit,
		"limit":         limit,
	}
	r, err := cli.Api.POST(ctx, "https://www.feishu.cn/approval/openapi/v2/instance/list", &httpx.Option{Params: param})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resp struct {
		Data struct {
			Ids []string `test_helper:"instance_code_list"`
		} `test_helper:"data"`
	}
	if err := r.ParsedStruct(&resp); err != nil {
		return nil, errors.WithStack(err)
	}

	return resp.Data.Ids, nil
}

// AllApprovalInstanceIds 批量获取审批实例ID
func (cli *Client) AllApprovalInstanceIds(ctx context.Context, approvalCode string, startTime, endTime time.Time) ([]string, error) {
	page, limit := int64(1), int64(100)
	ids := make([]string, 0)

	retry := 3
	for {
		items, err := cli.GetApprovalInstanceIds(ctx, approvalCode, startTime, endTime, page, limit)
		if err != nil {
			if retry == 0 {
				return nil, errors.WithStack(err)
			}
			retry -= 1
			continue
		}

		ids = append(ids, items...)

		if len(items) < int(limit) {
			break
		}

		page += 1
		retry = 3
	}

	return ids, nil
}

// GetApprovalInstance 获取单个审批实例详情
func (cli *Client) GetApprovalInstance(ctx context.Context, instanceCode string) (*ApprovalInstance, error) {
	query := map[string]interface{}{
		"instance_code": instanceCode,
	}
	r, err := cli.Api.GET(ctx, "https://www.feishu.cn/approval/openapi/v2/instance/get", &httpx.Option{Query: query})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resp struct {
		Data ApprovalInstance `test_helper:"data"`
	}
	if err := r.ParsedStruct(&resp); err != nil {
		return nil, errors.WithStack(err)
	}

	return &resp.Data, nil
}
