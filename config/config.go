package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `yaml:"port"`
}

// DatabaseConfig 数据库配置 (支持 MySQL / PostgreSQL)
type DatabaseConfig struct {
	Driver   string `yaml:"driver"`   // "mysql" 或 "postgresql"
	Host     string `yaml:"host"`     // 数据库地址
	Port     int    `yaml:"port"`     // 数据库端口
	User     string `yaml:"user"`     // 用户名
	Password string `yaml:"password"` // 密码
	DBName   string `yaml:"dbname"`   // 数据库名
}

// Load 从 YAML 配置文件加载
//
// 查找顺序:
//  1. 环境变量 CONFIG_PATH 指定的路径 (仅用于定位文件, 不覆盖值)
//  2. 当前目录下的 config.yaml
//  3. 可执行文件同目录下的 config.yaml
func Load() *Config {
	cfg := &Config{
		Server: ServerConfig{Port: "8000"},
		Database: DatabaseConfig{
			Driver: "mysql", Host: "localhost", Port: 3306,
			User: "root", Password: "", DBName: "unigo",
		},
	}

	configPath := findConfigFile()
	if configPath == "" {
		log.Println("[Config] 未找到配置文件，使用默认值")
		return cfg
	}

	if err := loadFromFile(cfg, configPath); err != nil {
		log.Fatalf("[FATAL] 加载配置文件失败 (%s): %v", configPath, err)
	}
	log.Printf("[Config] 配置加载成功: %s", configPath)
	return cfg
}

// DSN 生成数据库连接字符串
func (d DatabaseConfig) DSN() string {
	switch d.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			d.User, d.Password, d.Host, d.Port, d.DBName)
	case "postgresql":
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			d.Host, d.User, d.Password, d.DBName, d.Port)
	default:
		return ""
	}
}

// ============================================================
//  内部函数
// ============================================================

func findConfigFile() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" && fileExists(path) {
		return path
	}

	for _, name := range []string{"config.yaml", "config.yml"} {
		if fileExists(name) {
			return name
		}
	}

	exePath, err := os.Executable()
	if err == nil {
		for _, name := range []string{"config.yaml", "config.yml"} {
			if p := filepath.Join(filepath.Dir(exePath), name); fileExists(p) {
				return p
			}
		}
	}
	return ""
}

func loadFromFile(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, cfg)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
