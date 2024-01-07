package grpc

// Option 类型是一个函数，接收一个 *Server 类型的指针参数。
type Option func(*Server)

// WithServerHost 为 Server 设置主机名。
//
// 参数：
// - host string 要设置的主机名。
//
// 返回值：
// - Option 一个 Option 类型的函数，用于接收一个 *Server 并对其进行配置。
func WithServerHost(host string) Option {
	return func(s *Server) {
		// 为 Server 设置主机名。
		s.host = host
	}
}

// WithServerPort 为 Server 设置端口号。
//
// 参数：
// - port int 要设置的端口号。
//
// 返回值：
// - Option 一个 Option 类型的函数，用于接收一个 *Server 并对其进行配置。
func WithServerPort(port int) Option {
	return func(s *Server) {
		// 为 Server 设置端口号。
		s.port = port
	}
}
