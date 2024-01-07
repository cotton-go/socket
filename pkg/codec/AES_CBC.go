package codec

import (
	"encoding/base64"

	"github.com/bytedance/sonic"
	"github.com/forgoer/openssl"
	"github.com/pkg/errors"
)

// AESCBC 结构体，用于存储 AES-CBC 加密算法的密钥和初始向量
type AESCBC struct {
	key []byte // 密钥
	iv  []byte // 初始向量
}

// NewAESCBC 函数，用于创建一个新的 AES-CBC 编码器实例
//
// 参数：
//   - key string 密钥，长度为 32 字节
//
// 返回值：
//   - ICodec 返回一个实现了 ICodec 接口的 AESCBC 结构体实例
func NewAESCBC(key string) ICodec {
	if key == "" {
		newKey, _ := generateRandomKey(32) // 如果密钥为空，则生成一个随机密钥
		key = string(newKey)               // 将密钥转换为字符串类型
	}

	// 你可能需要处理 key 的长度和格式，这里简单地使用 key 的字节数组
	return &AESCBC{key: []byte(key), iv: []byte(key)} // 返回一个包含密钥和初始向量的 AESCBC 结构体实例
}

// encrypt 函数，用于对给定的字节数组进行 AES-CBC 加密
//
// 参数：
//   - src []byte 需要加密的字节数组
//
// 返回值：
//   - string 返回经过 Base64 编码后的加密结果字符串
//   - error 返回错误信息，如果加密过程中出现异常则返回该异常
func (sc AESCBC) encrypt(src []byte) (string, error) {
	// 使用 OpenSSL 库中的 AesCBCEncrypt 函数进行加密，并将结果转换为 Base64 编码的字符串
	dst, _ := openssl.AesCBCEncrypt(src, sc.key, sc.iv, openssl.PKCS7_PADDING)
	return base64.StdEncoding.EncodeToString(dst), nil
}

// decrypt 是 AESCBC 结构体中的一个方法，用于对输入的密文进行解密。
//
// 参数：
// - src string: 需要解密的密文，类型为字符串。
//
// 返回值：
// - []byte: 解密后的明文，类型为字节切片。
// - error: 返回错误信息，如果解密过程中出现错误，则返回相应的错误信息。
func (sc AESCBC) decrypt(src string) ([]byte, error) {
	// 对输入的密文进行 base64 解码
	dst, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}

	// 使用 OpenSSL 库的 AesCBCDecrypt 方法对解码后的密文进行 AES-CBC 解密
	return openssl.AesCBCDecrypt(dst, sc.key, sc.iv, openssl.PKCS7_PADDING)
}

// Encode 是 AESCBC 结构体中的一个方法，用于对输入的值进行编码。
//
// 参数：
// - value any: 需要编码的值，类型为任意类型。
//
// 返回值：
// - any: 编码后的值，类型为任意类型。
// - error: 返回错误信息，如果编码过程中出现错误，则返回相应的错误信息。
func (sc AESCBC) Encode(value any) (any, error) {
	// 使用 sonic.Marshal 方法将输入的值编码为 Event 结构体
	b, _ := sonic.Marshal(Event{Value: value})
	// 调用 encrypt 方法对编码后的数据进行加密，并返回加密后的结果
	return sc.encrypt(b)
}

// Decode 是一个解密函数，用于解密输入的值并返回解密后的值和错误信息。
//
// 参数：
// - value any: 需要解密的值，类型为任意类型。
//
// 返回值：
// - any: 解密后的值，类型为任意类型。
// - error: 返回错误信息，如果解密或解析过程中出现错误，则返回相应的错误信息。
func (sc AESCBC) Decode(value any) (any, error) {
	// 使用 sc.decrypt 方法对输入的值进行解密
	data, err := sc.decrypt(value.(string))
	if err != nil {
		// 如果解密过程中出现错误，则返回错误信息
		return nil, errors.Wrap(err, "解密失败[1001]")
	}

	var event Event
	// 使用 sonic.Unmarshal 方法将解密后的数据解析为 Event 结构体
	if err := sonic.Unmarshal(data, &event); err != nil {
		// 如果解析过程中出现错误，则返回错误信息
		return nil, errors.Wrap(err, "解析失败[1002]")
	}

	// 返回解密后的值和 nil 表示没有错误
	return event.Value, nil
}
