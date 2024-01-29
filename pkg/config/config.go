package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// 默认值
const (
	DefaultTokenExpire     time.Duration = 3600
	DefaultTokenToleration time.Duration = 1200
	DefaultSecret          string        = "2YejrzYBZzr1An5QSkbB3vKiGQYmRGZyUSGugAub0a39QFdFg1DyFdtMbbIEAY94"
)

type (
	// Config 配置文件
	Config struct {
		Logging  Logging
		Server   Server
		Database Database
		Cache    Cache
		Security Security
		LDAP     LDAP
	}

	// Logging 日志配置
	Logging struct {
		Format     string `yaml:"format" mapstructure:"format"`
		Path       string `yaml:"path" mapstructure:"path"`
		Level      string `yaml:"level" mapstructure:"level"`
		MaxAge     int    `yaml:"max_age" mapstructure:"max_age"`
		MaxSize    int    `yaml:"max_size" mapstructure:"max_size"`
		MaxBackups int    `yaml:"max_backups" mapstructure:"max_backups"`
	}

	// Server HTTP服务端相关配置
	// Key Cert Host 是证书相关的配置
	Server struct {
		BindAddress        string        `yaml:"bind_address" mapstructure:"bind_address"`
		InsecurePort       string        `yaml:"insecure_port"  mapstructure:"insecure_port"`
		SecurePort         string        `yaml:"secure_port"  mapstructure:"secure_port"`
		ReadTimeout        time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
		WriteTimeout       time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
		TerminationTimeout time.Duration `yaml:"termination_timeout" mapstructure:"termination_timeout"`
		SecureKey          string        `yaml:"secure_key" mapstructure:"secure_key"`
		SecureCert         string        `yaml:"secure_cert" mapstructure:"secure_cert"`
		SecureHost         string        `yaml:"secure_host" mapstructure:"secure_host"`
	}

	// Database 数据库相关配置
	Database struct {
		Host            string        `yaml:"host" mapstructure:"host"`
		Port            string        `yaml:"port" mapstructure:"port"`
		User            string        `yaml:"user" mapstructure:"user"`
		Password        string        `yaml:"password" mapstructure:"password"`
		DatabaseName    string        `yaml:"database_name" mapstructure:"database_name"`
		MaxIdleConns    int           `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
		MaxOpenConns    int           `yaml:"max_open_conns" mapstructure:"max_open_conns"`
		ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	}
	DB struct {
		Type        string   `mapstructure:"type" json:"type" yaml:"type"`
		TablePrefix string   `mapstructure:"table-prefix" json:"tablePrefix" yaml:"table-prefix"`
		Mysql       Database `yaml:"mysql,omitempty"`
		// Sqlite3     Sqlite3 `yaml:"sqlite3,omitempty"`
	}
	// Cache 缓存相关配置
	Cache struct {
	}

	// Security 安全相关的配置
	Security struct {
		AdminUser           string        `yaml:"admin_user" mapstructure:"admin_user"`
		AdminPassword       string        `yaml:"admin_password" mapstructure:"admin_password"`
		AdminEmail          string        `yaml:"admin_email" mapstructure:"admin_email"`
		AdminPhone          string        `yaml:"admin_phone" mapstructure:"admin_phone"`
		SecretKey           string        `yaml:"secret_key" mapstructure:"secret_key"`
		TokenExpireTime     time.Duration `yaml:"token_expire_time" mapstructure:"token_expire_time"`
		TokenTolerationTime time.Duration `yaml:"token_toleration_time" mapstructure:"token_toleration_time"`
	}

	// LDAP LDAP认证相关的配置
	LDAP struct {
		Enable                         bool   `yaml:"enable" mapstructure:"enable"`
		Host                           string `yaml:"host" mapstructure:"host"`
		Port                           string `yaml:"port" mapstructure:"port"`
		BindDB                         string `yaml:"bind_dn" mapstructure:"bind_dn"`
		BindPwd                        string `yaml:"bind_pwd" mapstructure:"bind_pwd"`
		SearchBaseDNS                  string `yaml:"search_base_dns" mapstructure:"search_base_dns"`
		SearchFilter                   string `yaml:"search_filter" mapstructure:"search_filter"`
		GroupSearchFilter              string `yaml:"group_search_filter" mapstructure:"group_search_filter"`
		GroupSearchBaseDNS             string `yaml:"group_search_base_dns" mapstructure:"group_search_base_dns"`
		GroupSearchFilterUserAttribute string `yaml:"group_search_filter_user_attribute" mapstructure:"group_search_filter_user_attribute"`
	}
)

// String 将配置文件输出为字符串
func (config *Config) String() string {
	str, _ := yaml.Marshal(config)
	return "\n" + string(str)
}

// Load 加载配置文件
func (config *Config) Load(configFile string) error {
	// TODO: 判断文件大小
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("config read error: %v", err)
	}
	//将配置文件中的数据，写入 config 结构体中
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return fmt.Errorf("config yaml unmarshal error: %s", err)
	}
	config.SetDefault()
	return nil
}

// SetDefault 设置默认值
func (config *Config) SetDefault() {
	defaultTokenExpireTime(config)
	defaultTokenTolerationTime(config)
	defaultSecretKey(config)
}

func defaultTokenExpireTime(cfg *Config) {
	if cfg.Security.TokenExpireTime == 0 {
		cfg.Security.TokenExpireTime = DefaultTokenExpire
	}
}
func defaultTokenTolerationTime(cfg *Config) {
	if cfg.Security.TokenTolerationTime == 0 {
		cfg.Security.TokenTolerationTime = DefaultTokenToleration
	}
}

func defaultSecretKey(cfg *Config) {
	if cfg.Security.SecretKey == "" {
		cfg.Security.SecretKey = DefaultSecret
	}
}
