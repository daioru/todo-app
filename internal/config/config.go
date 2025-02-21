package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

var cfg *Config

type DB struct {
	DSN             string        `yaml:"DSN"`
	MaxOpenConns    int           `yaml:"maxOpenConns"`
	MaxIdleConns    int           `yaml:"MaxIdleConns"`
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"`
}

func (db *DB) GetDSN() string {
	return db.DSN
}

func (db *DB) GetMaxOpenConns() int {
	return db.MaxOpenConns
}

func (db *DB) GetMaxIdleConns() int {
	return db.MaxIdleConns
}

func (db *DB) GetConnMaxIdleTime() time.Duration {
	return db.ConnMaxIdleTime
}

func (db *DB) GetConnMaxLifetime() time.Duration {
	return db.ConnMaxLifetime
}

func ReadConfigYML(configYML string) error {
	if cfg != nil {
		return nil
	}

	file, err := os.Open(configYML)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	return nil
}

type Config struct {
	DB DB `yaml:"db"`
}

func GetConfigInstance() Config {
	if cfg != nil {
		return *cfg
	}

	return Config{}
}
