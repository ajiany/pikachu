package feishu

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io/ioutil"

	"github.com/ajiany/pikachu/tools/httpx"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Client struct {
	cfg *Config
	rds redis.UniversalClient

	Api    *httpx.Httpx
	Token  *Token
	Signer *Signer
}

func NewClient(cfg Config, rds redis.UniversalClient) *Client {
	cli := new(Client)
	cli.cfg = &cfg
	cli.rds = rds
	host := fmt.Sprintf("%s://%s", cfg.Scheme, cfg.Host)

	cli.Api = httpx.NewHttpx(httpx.WithHost(host), httpx.WithUrlPrefix(cfg.UrlPrefix))
	cli.Token = NewToken(cli)
	cli.Signer = NewSigner(cli)
	cli.Api.SetBeforeRequestHook(func(req *httpx.Request, opts *httpx.Option) {
		req.Req.Header.Set("Content-Type", "application/test_helper; charset=utf-8")
		req.Req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.Token.GetToken()))
	})
	cli.Api.SetAfterRequestHook(func(resp *httpx.Response) error {
		var r ResponseData
		if err := resp.ParsedStruct(&r); err != nil {
			return err
		}

		if r.Code != 0 {
			return fmt.Errorf("Request fail: code -> %d, msg -> %s", r.Code, r.Msg)
		}

		return nil
	})

	return cli
}

func (cli *Client) UnmarshalAndVerifyEvent(c *gin.Context) (*EventData, error) {
	event, err := cli.UnmarshalEvent(c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 回调验证
	if event.Type == EventUrlVerification {
		return event, nil
	}

	if !cli.VerifySign(c, string(event.RawData)) {
		return nil, ErrInvalidSign
	}

	if !cli.VerifyToken(event) {
		return nil, ErrInvalidToken
	}

	return event, nil
}

func (cli *Client) VerifySign(c *gin.Context, body string) bool {
	return cli.Signer.Verify(c, body)
}

func (cli *Client) VerifyToken(event *EventData) bool {
	return event.VerificationToken() == cli.cfg.VerificationToken
}

func (cli *Client) UnmarshalEvent(c *gin.Context) (*EventData, error) {
	var event EventData

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	event.RawData = data

	// 解密
	if cli.cfg.Encrypted() {
		var encryptData EncryptEventData

		if err := json.Unmarshal(data, &encryptData); err != nil {
			return nil, errors.WithStack(err)
		}

		if encryptData.Encrypt != "" {
			eventStr, err := cli.Signer.Decrypt(encryptData.Encrypt)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			data = []byte(eventStr)
		}
	}

	logrus.Info(string(data))

	if err := json.Unmarshal(data, &event); err != nil {
		return nil, errors.WithStack(err)
	}

	return &event, nil
}
