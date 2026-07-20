package conf

import (
	"log"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ========== 配置结构体定义 (补全所有字段) ==========

// Config 总配置
type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Logger LoggerConfig `mapstructure:"logger"`
	ORM    ORMConfig    `mapstructure:"orm"`
	JWT    JWTConfig    `mapstructure:"jwt"`
}

// AppConfig 应用配置 (包含你截图里的所有字段)
type AppConfig struct {
	Name          string `mapstructure:"name"`
	Version       string `mapstructure:"version"`
	Mode          string `mapstructure:"mode"`
	Addr          string `mapstructure:"addr"`
	Host          string `mapstructure:"host"`
	Resource      string `mapstructure:"resource"`
	FfprobePath   string `mapstructure:"ffprobe_path"`
	Env           string `mapstructure:"env"`
	MpHost        string `mapstructure:"mp_host"`
	AdOctopusHost string `mapstructure:"ad_octopus_host"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

// ORMConfig 数据库配置
type ORMConfig struct {
	Dsn          string `mapstructure:"dsn"` // 直接使用完整 DSN（优先）
	Driver       string `mapstructure:"driver"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	ParseTime    bool   `mapstructure:"parse_time"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	Expire time.Duration `mapstructure:"expire"`
	Issuer string        `mapstructure:"issuer"`
}

// Conf 全局配置实例 (公开变量)
var Conf *Config

// setDefaults 设置默认值，确保 YAML 缺失字段时程序不会崩溃
func setDefaults() {
	viper.SetDefault("app.mode", "debug")
	viper.SetDefault("app.addr", ":8080")
	viper.SetDefault("app.env", "local")
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("orm.driver", "mysql")
	viper.SetDefault("orm.host", "127.0.0.1")
	viper.SetDefault("orm.port", 3306)
	viper.SetDefault("orm.user", "root")
	viper.SetDefault("orm.dbname", "ginvueblog")
	viper.SetDefault("orm.charset", "utf8mb4")
	viper.SetDefault("orm.parse_time", true)
	viper.SetDefault("orm.max_idle_conns", 10)
	viper.SetDefault("orm.max_open_conns", 100)
	viper.SetDefault("jwt.expire", 7200)
	viper.SetDefault("jwt.issuer", "ginVueBlog")
}

// Init 初始化配置：读取 YAML -> 环境变量覆盖 -> 热加载
func Init(configPath string) *Config {
	// 1. 设置默认值
	setDefaults()

	// 2. 读取 YAML 配置文件 (支持绝对或相对路径)
	viper.SetConfigFile(configPath) // 自动识别后缀 (yml/yaml)
	viper.SetConfigType("yaml")     // 显式声明类型

	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件找不到，只警告，不 panic (允许完全依赖环境变量)
		log.Printf("WARNING: 读取配置文件失败 (%s)，将完全依赖环境变量与默认值\n", err)
	} else {
		log.Printf("INFO: 成功加载配置文件: %s", configPath)
	}

	// 3. 【核心修改】开启环境变量支持，自动覆盖 YAML 字段
	viper.AutomaticEnv()
	// 将配置键中的 "." 替换为 "_" 以实现 APP_MODE -> app.mode
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 4. 反序列化到 Conf 指针 (修复了原来空指针的问题)
	Conf = &Config{}
	if err := viper.Unmarshal(Conf); err != nil {
		log.Fatalf("FATAL: 解析配置失败: %v", err)
	}

	// 5. 【企业级特性】配置文件热加载
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("INFO: 配置文件发生变更 (%s)，重新加载...", e.Name)
		newConf := &Config{}
		if err := viper.Unmarshal(newConf); err != nil {
			log.Printf("ERROR: 热加载失败: %v", err)
			return
		}
		// 原子性更新全局配置
		Conf = newConf
		log.Println("INFO: 配置热加载成功")
	})

	return Conf
}
