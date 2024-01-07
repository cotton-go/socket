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

// NewAESECB 创建一个新的 AES ECB 编码器实例
//
// 参数：
// - key string: 密钥，长度可以是任意值。如果为空字符串，将自动生成一个随机密钥。
//
// 返回值：
// - ICodec: 一个实现了 ICodec 接口的 AES ECB 编码器实例。
func NewAESECB(key string) ICodec {
	if key == "" {
		newKey, _ := generateRandomKey(32)
		key = string(newKey)
	}

	// 你可能需要处理 key 的长度和格式，这里简单地使用 key 的字节数组
	return &AESECB{key: []byte(key)}
}

// generateRandomKey 生成一个指定长度的随机密钥
//
// 参数：
// - length int: 需要生成的密钥长度。
//
// 返回值：
// - []byte: 一个长度为 length 的随机密钥字节数组。
// - error: 如果在生成密钥过程中出现错误，则返回相应的错误信息。
func generateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// encrypt 对输入的数据进行 ECB 模式加密，并返回加密后的 Base64 字符串。
//
// 参数：
// - src []byte: 需要加密的数据。
//
// 返回值：
// - string: 加密后的 Base64 字符串。如果加密过程中出现错误，则返回相应的错误信息。
func (sc *AESECB) encrypt(src []byte) (string, error) {
	dst, _ := openssl.AesECBEncrypt(src, sc.key, openssl.PKCS7_PADDING)
	return base64.StdEncoding.EncodeToString(dst), nil
}

// decrypt 对输入的 Base64 字符串进行解码，然后进行 ECB 模式解密，最后返回解密后的数据。
//
// 参数：
// - src string: 需要解密的 Base64 字符串。
//
// 返回值：
// - []byte: 解密后的数据。如果解密过程中出现错误，则返回相应的错误信息。
func (sc *AESECB) decrypt(src string) ([]byte, error) {
	dst, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}

	return openssl.AesECBDecrypt(dst, sc.key, openssl.PKCS7_PADDING)
}

// Encode 是 AESECB 类型的一个方法，用于对输入的值进行编码。
//
// 参数：
// - value any: 需要编码的值，类型为任意类型。
//
// 返回值：
// - any: 编码后的值，类型为任意类型。
// - error: 返回错误信息，如果编码过程中出现错误，则返回相应的错误信息。
func (sc AESECB) Encode(value any) (any, error) {
	// 使用 sonic.Marshal 对输入的值进行序列化
	b, _ := sonic.Marshal(Event{Value: value})
	// 调用 encrypt 方法对序列化后的值进行加密，并返回结果
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
func (sc AESECB) Decode(value any) (any, error) {
	// 使用 sc.decrypt 方法对输入的值进行解密
	data, err := sc.decrypt(value.(string))
	if err != nil {
		// 如果解密过程中出现错误，则返回错误信息
		return nil, errors.Wrap(err, "解密失败[1001]")
	}

	var event Event
	// 使用 sonic.Unmarshal 方法对解密后的数据进行解析
	if err := sonic.Unmarshal(data, &event); err != nil {
		// 如果解析过程中出现错误，则返回错误信息
		return nil, errors.Wrap(err, "解析失败[1002]")
	}

	// 返回解析后的值和 nil 表示没有错误
	return event.Value, nil
}
