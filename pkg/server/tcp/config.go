package tcp

type Config struct {
	Codec  string
	Secret string
	Host   string
	Port   int
	Cache  string
	// Redis  cache.Redis
}
