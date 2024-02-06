package application

import (
	"os"
	"strconv"
)

type Config struct {
	RedisAddress string
	Port         uint16
}

func LoadConfig() Config {
	cfg := Config{
		RedisAddress: "localhost:6379",
		Port:         3000,
	}

	if redisAddress, ok := os.LookupEnv("REDIS_ADDR"); ok {
		cfg.RedisAddress = redisAddress
	}

	if serverPort, ok := os.LookupEnv("ORDERS_SERVICE_SERVER_PORT"); ok {
		serverPortUInt64, err := strconv.ParseUint(serverPort, 10, 16)
		if err == nil {
			cfg.Port = uint16(serverPortUInt64)
		}
	}

	return cfg
}
