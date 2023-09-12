package config

import (
	"flag"
	"git.miem.hse.ru/1206/app"
	"git.miem.hse.ru/1206/app/logger"
	perms "git.miem.hse.ru/1206/app/permissions-v2"
	"git.miem.hse.ru/1206/app/storage/stpg"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	GRPC        app.GRPCConfig            `yaml:"grpc"`
	Postgres    stpg.Config               `yaml:"psql"`
	Logger      logger.Config             `yaml:"logger"`
	Permissions perms.PermissionsDbConfig `yaml:"permissions"`
	Jaeger      logger.JaegerConfig       `yaml:"jaeger"`
	S3          S3                        `yaml:"s3"`
	Education   app.GRPCConfig            `yaml:"education"`
}

type S3 struct {
	Endpoint       string        `yaml:"endpoint"`
	PublicEndpoint string        `yaml:"publicEndpoint"` // Full endpoint for remote calling (with https://)
	Region         string        `yaml:"region"`
	BucketName     string        `yaml:"bucketName"`
	AccessKeyId    string        `yaml:"accessKeyId"`
	SecretKey      string        `yaml:"secretKey"`
	DisableSSL     bool          `yaml:"disableSSL"`
	UrlLifespan    time.Duration `yaml:"urlLifespan"`
	FileNameSalt   string        `yaml:"fileNameSalt"` // Salt for creating file storage key
}

var config *Config

func Get() *Config {
	if config == nil {
		config = &Config{}
	}
	return config
}

func Init() (*Config, error) {
	filePath := flag.String("c", "etc/config.yml", "Path to configuration file")
	flag.Parse()
	config = &Config{}
	data, err := os.ReadFile(*filePath)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}
