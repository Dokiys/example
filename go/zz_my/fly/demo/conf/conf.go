package conf

import (
	"encoding/json"

	"github.com/Dokiys/go_test/go/zz_my/fly"
)

type Config struct {
	One string `json:"one"`
	Two string `json:"two"`
}

type BootConfigLoader struct{}

func (b *BootConfigLoader) Load(v any) {
	err := json.Unmarshal([]byte(`{"one":"1", "two": "2"}`), v)
	if err != nil {
		panic(err)
	}
}

func NewConfig(loader fly.ConfigLoader) (result *Config) {
	loader.Load(&result)
	return
}
