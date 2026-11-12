package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env       string   `mapstructure:"env"`
	HTTP      HTTP     `mapstructure:"http"`
	OrderHTTP HTTP     `mapstructure:"order_http"`
	Log       Log      `mapstructure:"log"`
	PG        PG       `mapstructure:"pg"`
	Redis     Redis    `mapstructure:"redis"`
	Auth      Auth     `mapstructure:"auth"`
	Wechat    Wechat   `mapstructure:"wechat"`
	Identity  Identity `mapstructure:"identity"`
}

type HTTP struct {
	Port            string        `mapstructure:"port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type PG struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Database        string        `mapstructure:"database"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	TimeZone        string        `mapstructure:"time_zone"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type Redis struct {
	Addr         string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type Auth struct {
	Issuer          string        `mapstructure:"issuer"`
	TokenSecret     string        `mapstructure:"token_secret"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
}

type Wechat struct {
	APIBaseURL string `mapstructure:"api_base_url"`
	MiniAppID  string `mapstructure:"mini_app_id"`
	MiniSecret string `mapstructure:"mini_secret"`
	OpenAppID  string `mapstructure:"open_app_id"`
	OpenSecret string `mapstructure:"open_secret"`
}

type Identity struct {
	Provider  string        `mapstructure:"provider"`
	VerifyURL string        `mapstructure:"verify_url"`
	AppID     string        `mapstructure:"app_id"`
	APIKey    string        `mapstructure:"api_key"`
	Timeout   time.Duration `mapstructure:"timeout"`
}

// Load 加载配置。默认优先读取config/config.yaml，兼容旧的configs目录。
func Load(configDir string) (*Config, error) {
	if configDir == "" {
		// 优先使用项目实际的config目录；保留configs作为兼容路径。
		for _, dir := range []string{"config", filepath.Join(exeDir(), "config"), "configs", filepath.Join(exeDir(), "configs")} {
			if st, err := os.Stat(filepath.Join(dir, "config.yaml")); err == nil && !st.IsDir() {
				configDir = dir
				break
			}
		}
		if configDir == "" {
			configDir = "config"
		}
	}

	env := strings.TrimSpace(os.Getenv("NX_ENV"))

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configDir)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config.yaml: %w", err)
	}
	// 优先NX_ENV，其次基础YAML中的env，留空时采用开发环境。
	if env == "" {
		env = strings.TrimSpace(v.GetString("env"))
	}
	if env == "" {
		env = "dev"
	}
	// 环境专属配置（可选）
	v.SetConfigName("config." + env)
	if err := v.MergeInConfig(); err != nil {
		// 不存在则忽略
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("merge config.%s.yaml: %w", env, err)
		}
	}

	// 环境变量覆盖：http.port -> NX_HTTP__PORT
	v.SetEnvPrefix("NX")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	v.AutomaticEnv()

	// 环境变量覆盖 duration 等字段时 viper 不做类型转换，这里对关键项手动兜底
	if s := os.Getenv("NX_HTTP__PORT"); s != "" {
		v.Set("http.port", s)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	cfg.Env = env
	return &cfg, nil
}

func isDir(p string) bool {
	st, err := os.Stat(p)
	return err == nil && st.IsDir()
}

func exeDir() string {
	// 拿到当前正在运行的这个程序自己的完整路径
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}
