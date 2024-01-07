package codec

import (
	"encoding/base64"

	"github.com/bytedance/sonic"
	"github.com/forgoer/openssl"
	"github.com/pkg/errors"
)

// DESECB 是一个使用 DES ECB 模式的解密器
type DESECB struct {
	key []byte // 密钥
}

// NewDESECB 创建一个新的 DESECB 实例
//
// 参数：
// - key string 密钥，长度为0时会自动生成一个随机密钥
//
// 返回值：
// - ICodec 返回一个 DESECB 实例
func NewDESECB(key string) ICodec {
	encode := &DESECB{key: []byte(key)}
	if key == "" {
		encode.key, _ = generateRandomKey(32)
	}

	// 你可能需要处理 key 的长度和格式，这里简单地使用 key 的字节数组
	return encode
}

// encrypt 使用 DES ECB 模式加密数据
//
// 参数：
// - src []byte 需要加密的数据
//
// 返回值：
// - string 加密后的数据，以 base64 编码的字符串形式返回
// - error 返回错误信息，如果加密失败则返回 nil
func (sc DESECB) encrypt(src []byte) (string, error) {
	// 实现 DES ECB 加密逻辑
	dst, _ := openssl.DesECBEncrypt(src, []byte(sc.key), openssl.PKCS7_PADDING)
	return base64.StdEncoding.EncodeToString(dst), nil
}

// decrypt 使用 DES ECB 模式解密数据
//
// 参数：
// - src string 需要解密的数据，以 base64 编码的字符串形式传入
//
// 返回值：
// - []byte 解密后的数据，以字节数组形式返回
// - error 返回错误信息，如果解密失败则返回 err
func (sc DESECB) decrypt(src string) ([]byte, error) {
	// 实现 DES ECB 解密逻辑
	dst, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}

	return openssl.DesECBDecrypt(dst, []byte(sc.key), openssl.PKCS7_PADDING)
}

// Encode 将任意类型的数据进行加密并返回加密后的字符串和错误信息
//
// 参数：
// - value any 需要加密的数据
//
// 返回值：
// - any 加密后的数据的任意类型表示，以及错误信息，如果加密失败则返回 nil 以及错误信息 "解密失败[1001]"
func (sc DESECB) Encode(value any) (any, error) {
	// 实现数据的加密逻辑，并将结果转换为字符串格式进行返回
	b, _ := sonic.Marshal(Event{Value: value})
	return sc.encrypt(b)
}

// Decode 将任意类型的数据进行解密并返回解密后的数据的任意类型表示和错误信息
//
// 参数：
// - value any 需要解密的数据，应为已加密的字符串形式的数据
//
// 返回值：
// - any 解密后的数据的任意类型表示，以及错误信息，如果解密失败则返回 nil 以及错误信息 "解析失败[1002]"
func (sc DESECB) Decode(value any) (any, error) {
	// 实现数据的解密逻辑，并将结果转换为原始的任意类型表示进行返回
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
