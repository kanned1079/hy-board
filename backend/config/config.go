package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Auth     AuthConfig     `yaml:"auth"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Type     string         `yaml:"type"`
	Sqlite   SqliteConfig   `yaml:"sqlite"`
	Mysql    MysqlConfig    `yaml:"mysql"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type SqliteConfig struct {
	File string `yaml:"file"`
}

type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	Charset  string `yaml:"charset"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"ssl_mode"`
}

type AuthConfig struct {
	JWTSecret      string `yaml:"jwt_secret"`
	JWTExpireHours int    `yaml:"jwt_expire_hours"`
	UniProxyToken  string `yaml:"uniproxy_token"`
}

var GlobalConfig *Config

func LoadConfig(path string) *Config {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	GlobalConfig = &cfg
	return GlobalConfig
}
