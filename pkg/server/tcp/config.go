package tcp

type Config struct {
	Codec  string       `yaml:"Codec"`
	Secret string       `yaml:"Secret"`
	Host   string       `yaml:"Host"`
	Port   int          `yaml:"Port"`
	Redis  *RedisConfig `yaml:"Redis"`
}

type RedisConfig struct {
	Addr       string `yaml:"Addr"`
	Username   string `yaml:"Username"`
	Password   string `yaml:"Password"`
	DB         int    `yaml:"DB"`
	MaxRetries int    `yaml:"MaxRetries"`
}
