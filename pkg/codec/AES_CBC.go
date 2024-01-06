package codec

import (
	"encoding/base64"
	"encoding/json"

	"github.com/forgoer/openssl"
	"github.com/pkg/errors"
)

type AESCBC struct {
	key []byte
	iv  []byte
}

func NewAESCBC(key string) ICodec {
	if key == "" {
		newKey, _ := generateRandomKey(32)
		key = string(newKey)
	}

	// 你可能需要处理 key 的长度和格式，这里简单地使用 key 的字节数组
	return &AESCBC{key: []byte(key), iv: []byte(key)}
}

// Encrypt 使用 AES 加密给定的数据
func (sc AESCBC) encrypt(src []byte) (string, error) {
	dst, _ := openssl.AesCBCEncrypt(src, sc.key, sc.iv, openssl.PKCS7_PADDING)
	return base64.StdEncoding.EncodeToString(dst), nil
}

// Decrypt 使用 AES 解密给定的数据
func (sc AESCBC) decrypt(src string) ([]byte, error) {
	dst, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}

	return openssl.AesCBCDecrypt(dst, sc.key, sc.iv, openssl.PKCS7_PADDING)
}

func (sc AESCBC) Encode(value any) (any, error) {
	b, _ := json.Marshal(Event{Value: value})
	return sc.encrypt(b)
}

func (sc AESCBC) Decode(value any) (any, error) {
	data, err := sc.decrypt(value.(string))
	if err != nil {
		return nil, errors.Wrap(err, "解密失败[1001]")
	}

	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, errors.Wrap(err, "解析失败[1002]")
	}

	return event.Value, nil
}
