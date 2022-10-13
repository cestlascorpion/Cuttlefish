package utils

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	// TODO:
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
	return false
}
