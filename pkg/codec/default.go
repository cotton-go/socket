package codec

// Default 结构体，用于存储默认值
type Default struct{}

func NewDefault() *Default {
	return &Default{}
}

// Encode函数，用于对任意类型的值进行编码。
//
// 参数：
//   - value any 需要编码的值
//
// 返回值：
//   - any 编码后的值
//   - error 返回错误信息
func (Default) Encode(value any) (any, error) {
	// 实现逻辑
	return value, nil
}

// Decode函数，用于对任意类型的值进行解码。
//
// 参数：
//   - value any 需要解码的值
//
// 返回值：
//   - any 解码后的值
//   - error 返回错误信息
func (Default) Decode(value any) (any, error) {
	// 实现逻辑
	return value, nil
}
