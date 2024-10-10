package config

import (
    "log"
    "github.com/spf13/viper"
)

type Config struct {
    ServerAddress   string
    PostgresConn    string
    PostgresUsername string
    PostgresPassword string
}

func LoadConfig() (*Config, error) {
    viper.AutomaticEnv()

    config := &Config{
        ServerAddress:    viper.GetString("SERVER_ADDRESS"),
        PostgresConn:     viper.GetString("POSTGRES_CONN"),
        PostgresUsername: viper.GetString("POSTGRES_USERNAME"),
        PostgresPassword: viper.GetString("POSTGRES_PASSWORD"),
    }

    if config.ServerAddress == "" || config.PostgresConn == "" {
        log.Fatal("Missing required environment variables")
    }

    return config, nil
}
