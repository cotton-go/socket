package codec

import "strings"

// ICodec 接口定义了编码和解码方法，用于对任意类型的数据进行编解码操作。
type ICodec interface {
	// Encode 方法接受一个任意类型的参数 value,返回一个 any 类型的结果和一个 error 类型的错误信息。
	Encode(value any) (any, error)

	// Decode 方法接受一个任意类型的参数 value,返回一个 any 类型的结果和一个 error 类型的错误信息。
	Decode(value any) (any, error)
}

// New 返回一个基于提供的类型和密钥的新 ICodec 实现。
//
// 参数：
// - typec: 一个表示编解码器类型的字符串。
// - secret: 一个表示加密密钥的字符串。
//
// 返回值：
// - resp: 一个根据提供的类型确定的 ICodec 实例。
func New(typec, secret string) ICodec {
	var resp ICodec
	switch strings.ToUpper(typec) {
	case "AESCBC":
		resp = NewAESCBC(secret)
	case "AESECB":
		resp = NewAESECB(secret)
	case "DESCBC":
		resp = NewDESCBC(secret)
	case "DESECB":
		resp = NewDESECB(secret)
	default:
		resp = NewDefault()
	}

	return resp
}
