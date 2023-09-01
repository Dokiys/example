package conf

type Config struct {
	One string
	Two string
}

func InitConf() *Config {
	return &Config{One: "1", Two: "2"}
}
