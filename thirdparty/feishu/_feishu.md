# 飞书快速开发SDK

## 初始化feishu client
1. 定义飞书appid和飞书app_secret
2. 定义加密参数
```go
var _feishu *feishu.Client
var _feishuOnce sync.Once

func Feishu() *feishu.Client {
	_feishuOnce.Do(func() {
		cfg := feishu.DefaultCfg()
		cfg.AppId = Cfg.FeishuAppId
		cfg.AppSecret = Cfg.FeishuAppSecret
		cfg.EncryptKey = Cfg.EventEncryptKey
		// cfg.VerificationToken = Cfg.EventEncryptKey

		_feishu = feishu.NewClient(cfg, Redis())
	})

	return _feishu
}
```