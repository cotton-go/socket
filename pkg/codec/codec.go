package codec

// ICodec 接口定义了编码和解码方法，用于对任意类型的数据进行编解码操作。
type ICodec interface {
	// Encode 方法接受一个任意类型的参数 value,返回一个 any 类型的结果和一个 error 类型的错误信息。
	Encode(value any) (any, error)

	// Decode 方法接受一个任意类型的参数 value,返回一个 any 类型的结果和一个 error 类型的错误信息。
	Decode(value any) (any, error)
}
