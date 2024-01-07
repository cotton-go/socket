package codec

import (
	"encoding/base64"

	"github.com/bytedance/sonic"
	"github.com/forgoer/openssl"
	"github.com/pkg/errors"
)

// DESCBC 是一个DES加密算法的CBC模式实现
type DESCBC struct {
	key []byte
	iv  []byte
}

// NewDESCBC 创建一个新的 DESCBC 实例
// 参数：
//   - key string 密钥，如果为空，则生成一个随机密钥
//
// 返回值：
//   - ICodec 返回一个 DESCBC 实例
func NewDESCBC(key string) ICodec {
	if key == "" {
		newKey, _ := generateRandomKey(32)
		key = string(newKey)
	}

	// 你可能需要处理 key 的长度和格式，这里简单地使用 key 的字节数组
	return &DESCBC{key: []byte(key), iv: []byte(key)}
}

// encrypt 使用给定的密钥和初始化向量对输入的字节切片进行加密，并返回加密后的base64编码字符串和错误信息。
//
// 参数：
// - src []byte 需要加密的数据，以字节切片形式给出。
//
// 返回值：
// - string 加密后的数据，以base64编码字符串形式返回。
// - error 如果在加密过程中发生错误，则返回错误信息。
func (sc DESCBC) encrypt(src []byte) (string, error) {
	// 实现加密逻辑
	dst, _ := openssl.DesCBCEncrypt(src, sc.key, sc.iv, openssl.PKCS7_PADDING)
	return base64.StdEncoding.EncodeToString(dst), nil
}

// decrypt 使用给定的密钥和初始化向量对输入的base64编码字符串进行解密，并返回解密后的字节切片和错误信息。
//
// 参数：
// - src string 需要解密的数据，以base64编码字符串形式给出。
//
// 返回值：
// - []byte 解密后的数据，以字节切片形式返回。
// - error 如果在解密过程中发生错误，则返回错误信息。
func (sc DESCBC) decrypt(src string) ([]byte, error) {
	// 实现解密逻辑
	dst, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}

	return openssl.DesCBCDecrypt(dst, sc.key, sc.iv, openssl.PKCS7_PADDING)
}

// Encode 将输入的任意类型数据使用给定的密钥和初始化向量进行加密后，再进行base64编码，并返回编码后的数据和错误信息。
//
// 参数：
// - value any 需要加密的数据，可以是任意类型。
//
// 返回值：
// - any 编码后的数据，以任意类型形式返回。
// - error 如果在加密或编码过程中发生错误，则返回错误信息。
func (sc DESCBC) Encode(value any) (any, error) {
	// 实现编码逻辑
	b, _ := sonic.Marshal(Event{Value: value})
	return sc.encrypt(b)
}

// Decode 将输入的base64编码的任意类型数据使用给定的密钥和初始化向量进行解密，然后解析为原始数据并返回原始数据和错误信息。
//
// 参数：
// - value any 需要解密并解析的数据，可以是任意类型。需要先将其转换为base64编码的字符串形式。
//
// 返回值：
// - any 解密并解析后的数据，以任意类型形式返回。
// - error 如果在解密或解析过程中发生错误，则返回错误信息。如果解密失败，会返回特定的错误信息"解密失败[1001]"。如果解析失败，会返回特定的错误信息"解析失败[1002]"。
func (sc DESCBC) Decode(value any) (any, error) {
	// 实现解密和解析逻辑
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
