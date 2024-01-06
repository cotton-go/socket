package codec

import (
	"encoding/base64"
	"encoding/json"

	"github.com/forgoer/openssl"
	"github.com/pkg/errors"
)

type DESCBC struct {
	key []byte
	iv  []byte
}

func NewDESCBC(key string) ICodec {
	if key == "" {
		newKey, _ := generateRandomKey(32)
		key = string(newKey)
	}

	// 你可能需要处理 key 的长度和格式，这里简单地使用 key 的字节数组
	return &DESCBC{key: []byte(key), iv: []byte(key)}
}

// Encrypt 使用 AES 加密给定的数据
func (sc DESCBC) encrypt(src []byte) (string, error) {
	dst, _ := openssl.DesCBCEncrypt(src, sc.key, sc.iv, openssl.PKCS7_PADDING)
	return base64.StdEncoding.EncodeToString(dst), nil
}

// Decrypt 使用 AES 解密给定的数据
func (sc DESCBC) decrypt(src string) ([]byte, error) {
	dst, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}

	return openssl.DesCBCDecrypt(dst, sc.key, sc.iv, openssl.PKCS7_PADDING)
}

func (sc DESCBC) Encode(value any) (any, error) {
	b, _ := json.Marshal(Event{Value: value})
	return sc.encrypt(b)
}

func (sc DESCBC) Decode(value any) (any, error) {
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
