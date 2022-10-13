package utils

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type RedisConfig struct {
	Network   string `json:"network,omitempty" yaml:"network"`
	Addr      string `json:"addr,omitempty" yaml:"addr"`
	Username  string `json:"username,omitempty" yaml:"username"`
	Password  string `json:"password,omitempty" yaml:"password"`
	Database  int    `json:"database,omitempty" yaml:"database"`
	KeyPrefix string `json:"key_prefix,omitempty" yaml:"key_prefix"`
}

func (r *RedisConfig) check() bool {
	if len(r.Addr) == 0 {
		log.Errorf("invalid redis parameter")
		return false
	}
	if len(r.Network) == 0 {
		r.Network = "tcp"
	}
	if len(r.KeyPrefix) == 0 {
		log.Errorf("invalid redis key prefix")
		return false
	}
	return true
}

type Config struct {
	Redis *RedisConfig `json:"redis,omitempty" yaml:"redis"`
}

func NewConfig(ctx context.Context, path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Errorf("os open err %+v", err)
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	cfg := &Config{}
	err = yaml.NewDecoder(f).Decode(cfg)
	if err != nil {
		log.Errorf("yaml decode err %+v", err)
		return nil, err
	}

	if !cfg.Check() {
		log.Errorf("check config failed")
		return nil, ErrInvalidParameter
	}

	return cfg, nil
}

func (c *Config) Check() bool {
	if c.Redis == nil {
		log.Errorf("redis config missing")
		return false
	}
	return c.Redis.check()
}

func NewTestConfig() *Config {
	return &Config{
		Redis: &RedisConfig{
			Network:   "tcp",
			Addr:      "127.0.0.1:6379",
			KeyPrefix: defaultPrefix,
		},
	}
}
