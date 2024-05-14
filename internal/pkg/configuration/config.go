package configuration

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	localConfig = "../local_configuration"

	Environment = "ENVIRONMENT"
	LogLevel    = "LOGLEVEL"

	DBName     = "DB.NAME"
	DBPort     = "DB.PORT"
	DBHost     = "DB.HOST"
	DBUsername = "DB.USERNAME"
	DBPassword = "DB.PASSWORD"
	DBParams   = "DB.PARAMS"

	JWTSecret = "JWT.SECRET"
	Salt      = "JWT.BCRYPTSALT"

	S3ID     = "S3.ID"
	S3Secret = "S3.SECRETKEY"
	S3Bucket = "S3.BUCKETNAME"
	S3Region = "S3.REGION"
)

var requiredConfig = []string{
	Environment,
	LogLevel,
	DBName,
	DBPort,
	DBHost,
	DBUsername,
	DBPassword,
	DBParams,
	JWTSecret,
	Salt,
	S3ID,
	S3Secret,
	S3Bucket,
	S3Region,
}

type S3Config struct {
	ID     string
	Secret string
	Bucket string
	Region string
}

type Configuration struct {
	Environment string
	LogLevel    string

	DBName     string
	DBPort     string
	DBHost     string
	DBUsername string
	DBPassword string
	DBParams   string

	JWTSecret string
	Salt      string

	S3Config S3Config
}

func NewConfiguration() (*Configuration, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	homepath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	localConfigPath := fmt.Sprintf("%v/%v", homepath, localConfig)
	viper.AddConfigPath(localConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err)
	}
	// DB.NAME -> DB_NAME

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	s3Config := S3Config{
		ID:     viper.GetString(S3ID),
		Secret: viper.GetString(S3Secret),
		Bucket: viper.GetString(S3Bucket),
		Region: viper.GetString(S3Region),
	}

	config := &Configuration{
		Environment: viper.GetString(Environment),
		LogLevel:    viper.GetString(LogLevel),

		DBName:     viper.GetString(DBName),
		DBPort:     viper.GetString(DBPort),
		DBHost:     viper.GetString(DBHost),
		DBUsername: viper.GetString(DBUsername),
		DBPassword: viper.GetString(DBPassword),
		DBParams:   viper.GetString(DBParams),

		JWTSecret: viper.GetString(JWTSecret),
		Salt:      viper.GetString(Salt),

		S3Config: s3Config,
	}

	if err := validateConfig(requiredConfig); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(requiredConfig []string) error {
	for _, config := range requiredConfig {
		if viper.GetString(config) == "" {
			return fmt.Errorf("required config %s is not set", config)
		}
	}
	return nil
}
