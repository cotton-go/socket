package codec

import (
	"encoding/base64"
	"encoding/json"

	"github.com/forgoer/openssl"
	"github.com/pkg/errors"
)

type DESECB struct {
	key []byte
}

// NewDESECB 创建一个新的 DESECB 实例
func NewDESECB(key string) ICodec {
	encode := &DESECB{key: []byte(key)}
	if key == "" {
		encode.key, _ = generateRandomKey(32)
	}

	// 你可能需要处理 key 的长度和格式，这里简单地使用 key 的字节数组
	return encode
}

// Encrypt 使用 AES 加密给定的数据
func (sc DESECB) encrypt(src []byte) (string, error) {
	dst, _ := openssl.DesECBEncrypt(src, []byte(sc.key), openssl.PKCS7_PADDING)
	return base64.StdEncoding.EncodeToString(dst), nil
}

// Decrypt 使用 AES 解密给定的数据
func (sc DESECB) decrypt(src string) ([]byte, error) {
	dst, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}

	return openssl.DesECBDecrypt(dst, []byte(sc.key), openssl.PKCS7_PADDING)
}

func (sc DESECB) Encode(value any) (any, error) {
	b, _ := json.Marshal(Event{Value: value})
	return sc.encrypt(b)
}

func (sc DESECB) Decode(value any) (any, error) {
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
