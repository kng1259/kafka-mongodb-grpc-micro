package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig
	Kafka   KafkaConfig
	MongoDB MongoDBConfig
	GRPC    GRPCConfig
}

type ServerConfig struct {
	Port int
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
	Enabled bool // Add this field
}

type MongoDBConfig struct {
	URI        string
	Database   string
	Collection string
}

type GRPCConfig struct {
	Host    string
	Port    int
	Timeout time.Duration
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Set default values
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("kafka.topic", "products")
	viper.SetDefault("kafka.groupid", "product-consumer-group")
	viper.SetDefault("kafka.enabled", true) // Default to enabled
	viper.SetDefault("mongodb.database", "productdb")
	viper.SetDefault("mongodb.collection", "products")
	viper.SetDefault("grpc.host", "localhost")
	viper.SetDefault("grpc.port", 9080)
	viper.SetDefault("grpc.timeout", 10)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("unable to decode config: %v", err))
	}

	return &config
}
