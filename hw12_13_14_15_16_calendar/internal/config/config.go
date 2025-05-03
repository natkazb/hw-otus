package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf `yaml:"logger"`
	HTTP    HTTPConf   `yaml:"http"`
	GRPC    GRPCConf   `yaml:"grpc"`
	Storage Storage    `yaml:"storage"`
	Rabbit  RabbitConf `yaml:"rabbit"`
}

type LoggerConf struct {
	Level    string `yaml:"level"`
	FileName string `yaml:"fileName"`
}

type HTTPConf struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	Timeout string `yaml:"timeout"`
}

type GRPCConf struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	Timeout string `yaml:"timeout"`
}

type Storage struct {
	Type string  `yaml:"type"`
	SQL  SQLConf `yaml:"sql"`
}

type SQLConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbName"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
	Driver   string `yaml:"driver"`
}

type RabbitConf struct {
	Host      string        `yaml:"host"`
	Port      int           `yaml:"port"`
	User      string        `yaml:"user"`
	Password  string        `yaml:"password"`
	QueueName string        `yaml:"queueName"`
	Timeout   time.Duration `yaml:"timeout"`
}

func NewConfig(filePath string) (Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("error in opening file %s: %w", filePath, err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("error in decoding %s: %w", filePath, err)
	}
	return config, nil
}
