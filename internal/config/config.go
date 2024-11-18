package config

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    PostgresHost     string
    PostgresPort     string
    PostgresUser     string
    PostgresPassword string
    PostgresDB       string

    RestHost string
    RestPort string

    GrpcHost string
    GrpcPort string
}

func New() (*Config, error) {
    if err := godotenv.Load(); err != nil {
        return nil, fmt.Errorf("error loading .env file: %w", err)
    }

    return &Config{
        PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
        PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
        PostgresUser:     getEnv("POSTGRES_USER", "clicks_user"),
        PostgresPassword: getEnv("POSTGRES_PASSWORD", "clicks_password"),
        PostgresDB:       getEnv("POSTGRES_DB", "clicks_db"),

        RestHost: getEnv("REST_HOST", "0.0.0.0"),
        RestPort: getEnv("REST_PORT", "8080"),

        GrpcHost: getEnv("GRPC_HOST", "0.0.0.0"),
        GrpcPort: getEnv("GRPC_PORT", "50051"),
    }, nil
}

func (c *Config) GetPostgresDSN() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        c.PostgresHost,
        c.PostgresPort,
        c.PostgresUser,
        c.PostgresPassword,
        c.PostgresDB,
    )
}

func (c *Config) GetRestAddress() string {
    return fmt.Sprintf("%s:%s", c.RestHost, c.RestPort)
}

func (c *Config) GetGrpcAddress() string {
    return fmt.Sprintf("%s:%s", c.GrpcHost, c.GrpcPort)
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
