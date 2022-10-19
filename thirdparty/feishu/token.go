package feishu

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type Token struct {
	cli *Client
}

func NewToken(cli *Client) *Token {
	return &Token{
		cli: cli,
	}
}

func (t *Token) GetToken() string {
	token, err := t.tokenFromCache()
	if err == nil && token != "" {
		return token
	}

	token, err = t.ReflushToken()
	if err != nil {
		logrus.WithError(err).Error("reflush token fail")
	}
	return token
}

func (t *Token) ReflushToken() (string, error) {
	token, err := t.cli.GetTenantAccessToken(context.TODO())
	if err != nil {
		return "", err
	}

	err = t.cachedToken(token)
	return token.Token, err
}

func (t *Token) cachedToken(token *TenantAccessToken) error {
	expire := token.Expire - 10*60 // 提前10分钟过期
	err := t.cli.rds.Do(context.TODO(), "SET", t.cacheKey(), token.Token, "EX", expire).Err()
	if err != nil {
		return err
	}
	return nil
}

func (t *Token) tokenFromCache() (string, error) {
	token, err := t.cli.rds.Get(context.TODO(), t.cacheKey()).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (t *Token) cacheKey() string {
	return fmt.Sprintf("atlassian_bot:feishu:%s:access_token", t.cli.cfg.AppId)
}
