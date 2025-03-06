package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	SSLMode         string        `mapstructure:"ssl_mode"`
}

type SMSConfig struct {
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	Endpoint        string `mapstructure:"endpoint"`
	SignName        string `mapstructure:"sign_name"`
	TemplateCode    string `mapstructure:"template_code"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type AppConfig struct {
	Env      string         `mapstructure:"env"`
	Database DatabaseConfig `mapstructure:"database"`
	SMS      SMSConfig      `mapstructure:"sms"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

var Cfg *AppConfig

func LoadConfig(path string) error {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	viper.BindEnv("env")
	viper.BindEnv("database.host", "DB_HOST")         // 数据库主机名环境变量
	viper.BindEnv("database.port", "DB_PORT")         // 数据库端口环境变量
	viper.BindEnv("database.user", "DB_USER")         // 数据库用户名环境变量
	viper.BindEnv("database.password", "DB_PASSWORD") // 数据库密码环境变量
	viper.BindEnv("database.name", "DB_NAME")         // 数据库名称环境变量

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("config file not found at %s", path)
		}
		return fmt.Errorf("failed to read config:%w", err)
	}

	//将配置文件的数据解码到AppConfig
	if err := viper.Unmarshal(&Cfg); err != nil {
		return fmt.Errorf("failed to unmarsha config:%w", err)
	}

	//通过env字段区分环境 例如开发环境（development）、测试环境（test）或生产环境（production）等
	if env := os.Getenv("APP_ENV"); env != "" {
		Cfg.Env = env
	}
	// 环境变量覆盖配置（可选）
	if ak := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"); ak != "" {
		Cfg.SMS.AccessKeyID = ak
	}
	if as := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"); as != "" {
		Cfg.SMS.AccessKeySecret = as
	}

	// 对数据库配置进行校验，确保关键字段不为空或无效
	if Cfg.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if Cfg.Database.Port == 0 {
		return fmt.Errorf("database port is required")
	}

	if Cfg.Database.MaxOpenConns == 0 {
		Cfg.Database.MaxOpenConns = 25
	}
	if Cfg.Database.MaxIdleConns == 0 {
		Cfg.Database.MaxIdleConns = 5 // 数据库空闲连接数默认为 5
	}
	if Cfg.Database.ConnMaxLifetime == 0 {
		Cfg.Database.ConnMaxLifetime = time.Hour // 数据库连接最大生命周期默认为 1 小时
	}
	if Cfg.Database.SSLMode == "" {
		Cfg.Database.SSLMode = "disable" // 数据库 SSL 模式默认为禁用
	}
	return nil
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.Name)
}
