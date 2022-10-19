package feishu

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/ajiany/pikachu/tools/httpx"
	"github.com/pkg/errors"
)

const (
	CustomRobotHost   = "https://open.feishu.cn"
	CustomRobotPrefix = "/open-apis/bot/v2"
)

var _customRobotClient *httpx.Httpx
var _customRobotClientOnce sync.Once

func customRobotCli() *httpx.Httpx {
	if _customRobotClient != nil {
		return _customRobotClient
	}

	_customRobotClientOnce.Do(func() {
		cli := httpx.NewHttpx(httpx.WithHost(CustomRobotHost), httpx.WithUrlPrefix(CustomRobotPrefix))

		cli.SetBeforeRequestHook(func(req *httpx.Request, opts *httpx.Option) {
			req.Req.Header.Set("Content-Type", "application/json; charset=utf-8")
		})

		cli.SetAfterRequestHook(func(resp *httpx.Response) error {
			var r ResponseData
			if err := resp.ParsedStruct(&r); err != nil {
				return err
			}

			if r.Code != 0 {
				return NewErrResponse(r)
			}

			return nil
		})

		_customRobotClient = cli
	})

	return _customRobotClient
}

type CustomRobot struct {
	Token  string
	Secret string
}

type BotMessage struct {
	Timestamp string                 `json:"timestamp,omitempty"`
	Sign      string                 `json:"sign,omitempty"`
	MsgType   string                 `json:"msg_type"`
	Content   map[string]interface{} `json:"content"`
}

func NewCustomRobot(token, secret string) *CustomRobot {
	r := new(CustomRobot)
	r.Token = token
	r.Secret = secret
	return r
}

func (cr *CustomRobot) SendMessage(ctx context.Context, msg BotMessage) error {
	// 签名
	if cr.EnableSign() {
		ts := time.Now().Unix()
		sign, err := cr.Sign(ts)
		if err != nil {
			return errors.WithStack(err)
		}

		msg.Timestamp = fmt.Sprint(ts)
		msg.Sign = sign
	}

	_, err := customRobotCli().POST(ctx, fmt.Sprintf("/hook/%s", cr.Token), &httpx.Option{Params: msg})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (cr *CustomRobot) EnableSign() bool {
	return cr.Secret != ""
}

func (cr *CustomRobot) Sign(timestamp int64) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + cr.Secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
