package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/gommon/color"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

type DefaultConfig struct {
	Apps Apps `json:"apps"`
}

type Apps struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Address  string `json:"address"`
	HttpPort string `json:"httpPort"`
}

type MySQLDB struct {
	DSN                string
	MinIdleConnections int
	MaxOpenConnections int
	MaxLifetime        int
	LogMode            bool
}

type PostgreSQLDB struct {
	Username           string
	Password           string
	Name               string
	Schema             string
	Host               string
	Port               int
	SSLMode            string
	MinIdleConnections int
	MaxOpenConnections int
	MaxLifetime        int
	LogMode            bool
}

var (
	// v viper global variable
	v *viper.Viper
	// env used environment
	env string
)

func init() {
	// initializing viper
	v = viper.New()
}

// Load set viper configuration based on yaml files
// default configuration files are located in base path config.yaml.dist
func Load(env, localConfigName string) (err error) {
	v.SetDefault("apps.name", "personalities")
	v.SetDefault("appName", "personalities")
	v.SetDefault("apps.version", "dev")
	v.SetDefault("apps.address", "0.0.0.0")
	v.SetDefault("apps.httpPort", "8080")
	v.SetDefault("logger.sys.stdout", true)
	v.SetDefault("postgresql.skeleton.sslMode", "require")
	v.SetEnvPrefix("PERSONALITIES")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	return localConfig(env, localConfigName)
}

func localConfig(envShortCode, configName string) (err error) {
	var confEnv string
	v.SetConfigType("")

	if envShortCode != "" {
		confEnv = strings.ToLower(fmt.Sprintf("%s.%s", configName, envShortCode))
	} else {
		confEnv = strings.ToLower(configName)
	}

	if _, err := os.Stat(configName); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	baseName := filepath.Base(configName)
	if baseName == ".env" || strings.HasPrefix(baseName, ".env.") {
		if err := loadDotEnv(configName); err != nil {
			return err
		}
	} else {
		v.SetConfigFile(configName)
		if err := v.ReadInConfig(); err != nil {
			return err
		}
	}
	color.Println(color.Green(fmt.Sprintf("⇨ config: '%v'", confEnv)))

	if envShortCode != "" {
		env = envShortCode
		color.Println(color.Green(fmt.Sprintf("⇨ environment: %s", env)))
	}

	return
}

// LoadFromFile load configuration from file
func LoadFromFile(path string) {
	if err := loadDotEnv(path); err != nil {
		panic(err)
	}
}

func loadDotEnv(path string) error {
	values, err := gotenv.Read(path)
	if err != nil {
		return fmt.Errorf("read dotenv file: %w", err)
	}

	for key, value := range values {
		envName := dotenvEnvName(key)
		if _, exists := os.LookupEnv(envName); exists {
			continue
		}
		if err := os.Setenv(envName, value); err != nil {
			return fmt.Errorf("set environment variable %s: %w", envName, err)
		}
	}
	return nil
}

func dotenvEnvName(key string) string {
	normalized := strings.ToUpper(key)
	if normalized == "APP_ENV" || normalized == "CONFIG_FILE" || strings.HasPrefix(normalized, "PERSONALITIES_") {
		return normalized
	}
	key = strings.NewReplacer(".", "_").Replace(normalized)
	return "PERSONALITIES_" + key
}

// GetEnv getting env
func GetEnv() string {
	return env
}

// GetConfig gets the global Viper instance.
func GetConfig() *viper.Viper {
	return v
}

// Set sets the value for the key in the override register.
func Set(key string, value interface{}) { v.Set(key, value) }

// Get returns an interface. For a specific value use one of the Get____ methods.
func Get(key string) interface{} { return v.Get(key) }

// GetString returns the value associated with the key as an string.
func GetString(key string) string { return v.GetString(key) }

// GetBool returns the value associated with the key as a bool.
func GetBool(key string) bool { return v.GetBool(key) }

// GetInt returns the value associated with the key as an integer.
func GetInt(key string) int { return v.GetInt(key) }

// GetInt64 returns the value associated with the key as an integer.
func GetInt64(key string) int64 { return v.GetInt64(key) }

// GetFloat64 returns the value associated with the key as a float64.
func GetFloat64(key string) float64 { return v.GetFloat64(key) }

// GetDuration returns the value associated with the key as a duration.
func GetDuration(key string) time.Duration { return v.GetDuration(key) }

// GetStringSlice returns the value associated with the key as a slice of strings.
func GetStringSlice(key string) []string { return v.GetStringSlice(key) }

// GetStringMap returns the value associated with the key as a map of interfaces.
func GetStringMap(key string) map[string]interface{} { return v.GetStringMap(key) }

// GetStringMapString returns the value associated with the key as a map of strings.
func GetStringMapString(key string) map[string]string { return v.GetStringMapString(key) }

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func GetStringMapStringSlice(key string) map[string][]string { return v.GetStringMapStringSlice(key) }
