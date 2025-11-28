package configs

import (
	"os"

	"github.com/jinzhu/configor"
)

// Config 全局配置结构体
var Config = struct {
	App struct {
		Name string `yaml:"name"`
		Env  string `yaml:"env"`
		Port int    `yaml:"port"`
		JWT  struct {
			Secret    string `yaml:"secret"`
			ExpiresAt int64  `yaml:"expires_at"` // 有效期（秒）
			Issuer    string `yaml:"issuer"`
		} `yaml:"jwt"`
	} `yaml:"app"`
	Database struct {
		Driver          string `yaml:"driver"`
		DSN             string `yaml:"dsn"`
		MaxOpenConns    int    `yaml:"max_open_conns"`
		MaxIdleConns    int    `yaml:"max_idle_conns"`
		ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	} `yaml:"database"`
	Logger struct {
		Level    string `yaml:"level"`
		Path     string `yaml:"path"`
		Filename string `yaml:"filename"`
	} `yaml:"logger"`
}{}

// Init 初始化配置（加载 yaml 文件和环境变量）
func Init() error {
	// 优先加载环境变量（覆盖 yaml 配置）
	_ = os.Setenv("CONFIGOR_ENV_PREFIX", "APP")

	// 加载配置文件
	return configor.Load(&Config, "configs/app.yaml")
}
