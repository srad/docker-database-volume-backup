package conf

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"log"
	"os"
	"strings"
)

type AppConfig struct {
	Host              string `validate:"required"`
	User              string `validate:"required"`
	Password          string `validate:"required"`
	Database          string `validate:"required"`
	Cron              string `validate:"required"`
	BackupOnStart     bool
	BasicAuthPassword string
}

func LoadConfig() AppConfig {
	secretMap, err := readEnv([]string{"MYSQL_HOST", "MYSQL_USER", "MYSQL_PASSWORD", "MYSQL_DATABASE", "BASIC_AUTH_PASSWORD"})
	if err != nil {
		log.Fatal(err)
	}

	appConfig := AppConfig{
		Host:              secretMap["MYSQL_HOST"],
		User:              secretMap["MYSQL_USER"],
		Password:          secretMap["MYSQL_PASSWORD"],
		Database:          secretMap["MYSQL_DATABASE"],
		Cron:              os.Getenv("BACKUP_CRON"),
		BackupOnStart:     os.Getenv("BACKUP_ON_START") == "true",
		BasicAuthPassword: secretMap["BASIC_AUTH_PASSWORD"],
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	errValidate := validate.Struct(appConfig)
	if errValidate != nil {
		log.Fatal(errValidate)
	}

	return appConfig
}

func readEnv(secrets []string) (map[string]string, error) {
	result := make(map[string]string)

	for _, secret := range secrets {
		// First try to read env as secret...
		env := os.Getenv(secret)
		if env != "" {
			result[secret] = env
			continue
		}
		// ... then try to find docker secret
		envFile := os.Getenv(secret + "_FILE")
		if envFile == "" {
			return nil, fmt.Errorf("missing environment variable: %s", secret)
		}
		text, err := readFileAsString(envFile)
		if err != nil {
			return nil, fmt.Errorf("error reading secret %s: %v", secret, err)
		}
		result[secret] = text
	}

	return result, nil
}

func readFileAsString(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close() // Ensure the file is closed when the function exits

	// Read the entire file content
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Convert content to string and return
	return strings.TrimSpace(string(content)), nil
}
