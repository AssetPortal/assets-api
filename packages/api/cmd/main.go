package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/AssetPortal/assets-api/pkg/adapters/assets"
	"github.com/AssetPortal/assets-api/pkg/adapters/auth"
	"github.com/AssetPortal/assets-api/pkg/adapters/storage"
	"github.com/AssetPortal/assets-api/pkg/adapters/tokens"
	"github.com/AssetPortal/assets-api/pkg/app"
	"github.com/AssetPortal/assets-api/pkg/config"
	"github.com/AssetPortal/assets-api/pkg/middleware"
	"github.com/AssetPortal/assets-api/pkg/service"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func main() {

	cfg, err := config.MustGetConfig()
	if err != nil {
		log.Fatalf("error getting the configuration: %v", err)
	}
	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalf("error parsing the log level: %v", err)
	}
	logger := logrus.New()
	logger.SetLevel(lvl)

	// Open a database connection
	sqlDB, err := sql.Open("postgres", cfg.RWDBURL)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	defer sqlDB.Close()

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(10)                  // Maximum open connections
	sqlDB.SetMaxIdleConns(5)                   // Maximum idle connections
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)  // Idle timeout
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Max lifetime of a connection

	// Initialize Bun with database/sql and PostgreSQL dialect
	db := bun.NewDB(sqlDB, pgdialect.New())
	tokensRepository := tokens.NewTokensRepository(db)
	assetsRepository := assets.NewAssetsRepository(db)
	httpClient := &http.Client{
		Timeout: cfg.AuthConfiguration.HTTPTimeout,
	}
	authClient := auth.NewPolkadotClient(cfg.AuthConfiguration.APIURL, httpClient)
	authMiddleware := middleware.NewPolkadotAuth(tokensRepository, authClient)

	staticCreds := credentials.NewStaticCredentialsProvider(
		cfg.BucketConfiguration.AccessKey,
		cfg.BucketConfiguration.SecretKey,
		cfg.BucketConfiguration.Session,
	)

	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background(),
		awsConfig.WithRegion(cfg.BucketConfiguration.Region),
		awsConfig.WithCredentialsProvider(staticCreds),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}
	storageClient := storage.NewS3Uploader(awsCfg, cfg.BucketConfiguration.Name)
	assetsApp := app.NewAssetsApp(cfg, db, tokensRepository, assetsRepository, storageClient, logger)
	service := service.NewService(assetsApp, authMiddleware)

	service.Setup()
	service.Start()
}
