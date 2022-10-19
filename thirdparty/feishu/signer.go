package feishu

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type Signer struct {
	cli *Client
}

func NewSigner(cli *Client) *Signer {
	s := new(Signer)
	s.cli = cli
	return s
}

// 验证签名
func (s *Signer) Verify(c *gin.Context, body string) bool {
	sign := c.GetHeader("X-Lark-Signature")
	timestamp := c.GetHeader("X-Lark-Request-Timestamp")
	nonce := c.GetHeader("X-Lark-Request-Nonce")

	return s.verify(sign, timestamp, nonce, body)
}

func (s *Signer) verify(sign, timestamp, nonce, body string) bool {
	vsign := s.calculateSign(timestamp, nonce, s.cli.cfg.EncryptKey, body)
	return vsign == sign
}

// 解密
func (s *Signer) Decrypt(encrypt string) (string, error) {
	return s.decrypt(encrypt, s.cli.cfg.EncryptKey)
}

func (s *Signer) calculateSign(timestamp, nonce, encryptKey, bodystring string) string {
	var b strings.Builder
	b.WriteString(timestamp)
	b.WriteString(nonce)
	b.WriteString(encryptKey)
	b.WriteString(bodystring)
	bs := []byte(b.String())

	h := sha256.New()
	h.Write(bs)
	bs = h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func (s *Signer) decrypt(encrypt string, key string) (string, error) {
	buf, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", fmt.Errorf("base64StdEncode Error[%v]", err)
	}
	if len(buf) < aes.BlockSize {
		return "", fmt.Errorf("cipher too short")
	}
	keyBs := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(keyBs[:sha256.Size])
	if err != nil {
		return "", fmt.Errorf("AESNewCipher Error[%v]", err)
	}
	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]
	// CBC mode always works in whole blocks.
	if len(buf)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buf, buf)
	n := strings.Index(string(buf), "{")
	if n == -1 {
		n = 0
	}
	m := strings.LastIndex(string(buf), "}")
	if m == -1 {
		m = len(buf) - 1
	}
	return string(buf[n : m+1]), nil
}
