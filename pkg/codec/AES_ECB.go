package codec

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/bytedance/sonic"
	"github.com/forgoer/openssl"
	"github.com/pkg/errors"
)

// AESECB 结构体
type AESECB struct {
	key []byte
}

// NewAESECB 创建一个新的 AESECB 实例
func NewAESECB(key string) ICodec {
	if key == "" {
		newKey, _ := generateRandomKey(32)
		key = string(newKey)
	}

	// 你可能需要处理 key 的长度和格式，这里简单地使用 key 的字节数组
	return &AESECB{key: []byte(key)}
}

// GenerateRandomKey 生成指定长度的随机字节
func generateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Encrypt 使用 AES 加密给定的数据
func (sc AESECB) encrypt(src []byte) (string, error) {
	dst, _ := openssl.AesECBEncrypt(src, []byte(sc.key), openssl.PKCS7_PADDING)
	return base64.StdEncoding.EncodeToString(dst), nil
}

// Decrypt 使用 AES 解密给定的数据
func (sc AESECB) decrypt(src string) ([]byte, error) {
	dst, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}

	return openssl.AesECBDecrypt(dst, []byte(sc.key), openssl.PKCS7_PADDING)
}

func (sc AESECB) Encode(value any) (any, error) {
	b, _ := sonic.Marshal(Event{Value: value})
	return sc.encrypt(b)
}

func (sc AESECB) Decode(value any) (any, error) {
	data, err := sc.decrypt(value.(string))
	if err != nil {
		return nil, errors.Wrap(err, "解密失败[1001]")
	}

	var event Event
	if err := sonic.Unmarshal(data, &event); err != nil {
		return nil, errors.Wrap(err, "解析失败[1002]")
	}

	return event.Value, nil
}
