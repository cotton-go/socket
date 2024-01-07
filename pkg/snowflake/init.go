package snowflake

// 声明一个全局的 Snowflake 指针变量 workerID
var workerID *Snowflake

// init 初始化函数
func init() {
	// 使用 NewSnowflake 函数创建一个新的 Snowflake 实例，并将其赋值给 workerID
	workerID, _ = NewSnowflake(0)
}

// Next 生成一个新的 ID
//
// 参数：无
//
// 返回值：
//   - int64 返回新生成的 ID
func Next() int64 {
	// 调用 workerID 的 Generate 方法生成一个新的 ID,并将其赋值给 id
	id, _ := workerID.Generate()
	// 返回新生成的 ID
	return id
}
