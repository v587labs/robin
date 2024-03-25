package gorm

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
)

type Config struct {
	Group           string `mapstructure:"group" toml:"group" json:"group,omitempty" yaml:"group" `
	Type            string `mapstructure:"type" toml:"type" json:"type,omitempty" yaml:"type"`
	Host            string `mapstructure:"host" toml:"host" json:"host,omitempty" yaml:"host"`
	Port            int    `mapstructure:"port" toml:"port" json:"port,omitempty" yaml:"port"`
	User            string `mapstructure:"user" toml:"user" json:"user,omitempty" yaml:"user"`
	Pass            string `mapstructure:"pass" toml:"pass" json:"pass,omitempty" yaml:"pass"`
	Name            string `mapstructure:"name" toml:"name" json:"name,omitempty" yaml:"name"`
	ApplicationName string `mapstructure:"applicationName" toml:"applicationName" json:"application_name,omitempty" yaml:"applicationName"`
	MaxIdleCount    int    `mapstructure:"maxIdleCount" toml:"maxIdleCount" json:"maxIdleCount,omitempty" yaml:"maxIdleCount"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns" toml:"maxOpenConns" json:"maxOpenConns,omitempty" yaml:"maxOpenConns"`
	MaxLifetime     int64  `mapstructure:"maxLifetime" toml:"maxLifetime" json:"maxLifetime,omitempty" yaml:"maxLifetime"`
	Debug           bool   `mapstructure:"debug" toml:"debug" json:"debug,omitempty" yaml:"debug"`
	SslMode         string `mapstructure:"sslMode" toml:"sslMode" json:"sslMode,omitempty" yaml:"sslMode"`
}

func (cfg Config) BuildDsn() (string, error) {
	//	postgresql://[user[:password]@][netloc][:port][/dbname][?params]
	switch cfg.Type {
	case Postgres:
		return cfg.buildPostgresqlDsn()
	case SqlLite:
		return cfg.buildSqliteDsn()
	default:
		return "", fmt.Errorf("unsupported driver: %q", cfg.Type)
	}

}

func (cfg Config) BuildGormConfig() *gorm.Config {
	//	postgresql://[user[:password]@][netloc][:port][/dbname][?params]
	switch cfg.Type {
	case Postgres:
		return &gorm.Config{}
	case SqlLite:
		return &gorm.Config{}
	default:
		return &gorm.Config{}
	}

}

func (cfg Config) GetMaxLifetime() time.Duration {
	d := time.Duration(cfg.MaxLifetime)
	if d < time.Second {
		return d * time.Second
	}
	return d
}

func (cfg Config) buildPostgresqlDsn() (string, error) {
	//	postgresql://[user[:password]@][netloc][:port][/dbname][?params]
	dsn := strings.Builder{}
	dsn.WriteString("postgresql://")
	if len(cfg.User) > 0 {
		dsn.WriteString(cfg.User)
		dsn.WriteByte(':')
		dsn.WriteString(cfg.Pass)
		dsn.WriteByte('@')
	}
	if len(cfg.Host) > 0 {
		dsn.WriteString(cfg.Host)
		dsn.WriteByte(':')
		dsn.WriteString(strconv.Itoa(cfg.Port))
	}
	if len(cfg.Name) > 0 {
		dsn.WriteByte('/')
		dsn.WriteString(cfg.Name)
	}
	params := strings.Builder{}
	if len(cfg.ApplicationName) > 0 {
		params.WriteString("application_name=")
		params.WriteString(cfg.ApplicationName)
		params.WriteByte('&')
	}
	if len(cfg.SslMode) > 0 {
		params.WriteString("sslmode=")
		params.WriteString(cfg.SslMode)
		params.WriteByte('&')
	}

	if params.Len() > 0 {

		dsn.WriteByte('?')
		dsn.WriteString(params.String())
	}

	return dsn.String(), nil
}

func (cfg Config) buildSqliteDsn() (string, error) {
	//	postgresql://[user[:password]@][netloc][:port][/dbname][?params]
	// user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	dsn := strings.Builder{}
	if len(cfg.User) > 0 {
		dsn.WriteByte(' ')
		dsn.WriteString("user=")
		dsn.WriteString(cfg.User)
	}
	if len(cfg.Pass) > 0 {
		dsn.WriteByte(' ')
		dsn.WriteString("password=")
		dsn.WriteString(cfg.Pass)
	}

	if len(cfg.Host) > 0 {
		dsn.WriteByte(' ')
		dsn.WriteString("host=")
		dsn.WriteString(cfg.Host)
	}
	if cfg.Port > 0 {
		dsn.WriteByte(' ')
		dsn.WriteString("port=")
		dsn.WriteString(strconv.Itoa(cfg.Port))
	}
	if len(cfg.Name) > 0 {
		dsn.WriteByte(' ')
		dsn.WriteString("dbname=")
		dsn.WriteString(cfg.Name)
	}
	params := strings.Builder{}
	if len(cfg.ApplicationName) > 0 {
		params.WriteString("application_name=")
		params.WriteString(cfg.ApplicationName)
		params.WriteByte('&')
	}
	if len(cfg.SslMode) > 0 {
		dsn.WriteByte(' ')
		dsn.WriteString("sslmode=")
		dsn.WriteString(cfg.SslMode)
	}

	return strings.TrimSpace(dsn.String()), nil
}

func (cfg Config) GetDialector() gorm.Dialector {
	dsn, err := cfg.BuildDsn()
	if err != nil {
		panic(err)
	}
	switch cfg.Type {
	case Postgres:
		return postgres.Open(dsn)
	case SqlLite:
		return sqlite.Open(dsn)
	default:
		panic(fmt.Errorf("unsupported driver: %q", cfg.Type))
	}
}
