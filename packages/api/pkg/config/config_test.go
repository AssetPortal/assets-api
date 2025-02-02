package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/AssetPortal/assets-api/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestMustGetConfig_Success(t *testing.T) {
	// Set the environment variables
	os.Setenv("SERVICE_ADDRESS", ":8080")
	os.Setenv("HTTP_TIMEOUT", "30s")
	os.Setenv("DATABASE_URL", "localhost:5432")
	os.Setenv("MAX_REQUESTS_PER_SECOND", "20")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("TOKEN_EXPIRATION", "10m")
	os.Setenv("AUTH_API_URL", "http://auth.example.com")
	os.Setenv("AUTH_HTTP_TIMEOUT", "15s")
	os.Setenv("BUCKET_ACCESS_KEY", "my-access-key")
	os.Setenv("BUCKET_SECRET_KEY", "my-secret-key")
	os.Setenv("BUCKET_SESSION", "my-session")
	os.Setenv("BUCKET_REGION", "us-east-1")
	os.Setenv("BUCKET_NAME", "my-bucket")

	// Call MustGetConfig
	cfg, err := config.MustGetConfig()

	// Assert no error
	assert.NoError(t, err)

	// Assert the configuration values are correctly set
	assert.Equal(t, ":8080", cfg.ServiceAddress)
	assert.Equal(t, 30*time.Second, cfg.HTTPTimeout)
	assert.Equal(t, "localhost:5432", cfg.DatabaseConfiguration.URL)
	assert.Equal(t, 20, cfg.MaxRequestsPerSecond)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, 10*time.Minute, cfg.TokenExpiration)
	assert.Equal(t, "http://auth.example.com", cfg.AuthConfiguration.APIURL)
	assert.Equal(t, 15*time.Second, cfg.AuthConfiguration.HTTPTimeout)
	assert.Equal(t, "my-access-key", cfg.BucketConfiguration.AccessKey)
	assert.Equal(t, "my-secret-key", cfg.BucketConfiguration.SecretKey)
	assert.Equal(t, "my-session", cfg.BucketConfiguration.Session)
	assert.Equal(t, "us-east-1", cfg.BucketConfiguration.Region)
	assert.Equal(t, "my-bucket", cfg.BucketConfiguration.Name)
}

func TestMustGetConfig_DefaultValues(t *testing.T) {
	// Unset the environment variables
	os.Unsetenv("SERVICE_ADDRESS")
	os.Unsetenv("HTTP_TIMEOUT")
	os.Unsetenv("RW_DB_URL")
	os.Unsetenv("MAX_REQUESTS_PER_SECOND")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("TOKEN_EXPIRATION")
	os.Unsetenv("AUTH_API_URL")
	os.Unsetenv("AUTH_HTTP_TIMEOUT")
	os.Unsetenv("BUCKET_ACCESS_KEY")
	os.Unsetenv("BUCKET_SECRET_KEY")
	os.Unsetenv("BUCKET_SESSION")
	os.Unsetenv("BUCKET_REGION")
	os.Unsetenv("BUCKET_NAME")

	// Call MustGetConfig
	cfg, err := config.MustGetConfig()

	// Assert no error
	assert.NoError(t, err)

	// Assert default values are used
	assert.Equal(t, ":8000", cfg.ServiceAddress)
	assert.Equal(t, 20*time.Second, cfg.HTTPTimeout)
	assert.Equal(t, "", cfg.DatabaseConfiguration.URL)
	assert.Equal(t, 3, cfg.MaxRequestsPerSecond)
	assert.Equal(t, "warn", cfg.LogLevel)
	assert.Equal(t, 5*time.Minute, cfg.TokenExpiration)
	assert.Equal(t, "", cfg.AuthConfiguration.APIURL)
	assert.Equal(t, 20*time.Second, cfg.AuthConfiguration.HTTPTimeout)
	assert.Equal(t, "test", cfg.BucketConfiguration.AccessKey)
	assert.Equal(t, "test", cfg.BucketConfiguration.SecretKey)
	assert.Equal(t, "", cfg.BucketConfiguration.Session)
	assert.Equal(t, "", cfg.BucketConfiguration.Region)
	assert.Equal(t, "", cfg.BucketConfiguration.Name)
}
