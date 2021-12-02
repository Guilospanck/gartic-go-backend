package environments

import "os"

func initializeDevEnvironmentVariables() {
	// Database
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USERNAME", "postgres")
	os.Setenv("DB_PASSWORD", "123456")
	os.Setenv("DB_DATABASE_NAME", "gartic")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("DB_TIMEZONE", "America/Sao_Paulo")
}

func initializeStagingEnvironmentVariables() {
}

func initializeProductionEnvironmentVariables() {
	// Database
	os.Setenv("DB_HOST", "postgres")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USERNAME", "postgres")
	os.Setenv("DB_PASSWORD", "123456")
	os.Setenv("DB_DATABASE_NAME", "gartic")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("DB_TIMEZONE", "America/Sao_Paulo")
}

func init() {
	switch os.Getenv("GO_ENV") {
	case "development":
		initializeDevEnvironmentVariables()
	case "staging":
		initializeStagingEnvironmentVariables()
	case "production":
		initializeProductionEnvironmentVariables()
	default:
		initializeDevEnvironmentVariables()
	}
}
