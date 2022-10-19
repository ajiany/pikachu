package feishu

type Config struct {
	Scheme    string
	Host      string
	UrlPrefix string

	AppId             string
	AppSecret         string
	EncryptKey        string
	VerificationToken string
}

func DefaultCfg() Config {
	return Config{
		Scheme:    "https",
		Host:      "open.feishu.cn",
		UrlPrefix: "/open-apis",
	}
}

func (cfg *Config) Encrypted() bool {
	return cfg.EncryptKey != ""
}
