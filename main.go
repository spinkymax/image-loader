package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/spinkymax/image-loader/internal/config"
	"github.com/spinkymax/image-loader/internal/filestore"
	"github.com/spinkymax/image-loader/internal/repository"
	"github.com/spinkymax/image-loader/internal/server"
	"github.com/spinkymax/image-loader/internal/service"
	"github.com/spinkymax/image-loader/internal/telegram"
)

//	@title			Example Project API
//	@version		1.0
//	@description	Это API учебного проекта

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8000
//	@BasePath	/

func main() {

	ctx := context.Background()

	logger := logrus.New()

	cfg := &config.Config{}

	err := cfg.Process()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info(cfg.DB.Driver)

	db, err := sqlx.Connect(cfg.DB.Driver, fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s", cfg.DB.User,
		cfg.DB.Name, cfg.DB.SSLMode, cfg.DB.Password))
	if err != nil {
		logger.Fatal(err)
	}

	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.KeyID, cfg.Minio.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		logger.Fatal(err)
	}

	ok, err := minioClient.BucketExists(ctx, cfg.Minio.Bucket)
	if err != nil {
		logger.Fatal(err)
	}

	if !ok {
		err = minioClient.MakeBucket(context.Background(), cfg.Minio.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			logger.Fatal(err)
		}
	}

	fileStore := filestore.NewMinio(minioClient, cfg.Minio.Bucket)

	userRepo := repository.NewUserRepo(db, cfg.DB)

	imageRepo := repository.NewImageRepo(db, cfg.DB)

	err = userRepo.RunMigrations()
	if err != nil {
		logger.Warning(err)
	}

	controller := service.NewController(userRepo, imageRepo, cfg, fileStore)

	srv := server.NewServer(cfg.Port, logger, controller, cfg)
	srv.RegisterRoutes()

	bot, err := telegram.NewBot(cfg.TgBot.APIKey, logger, controller)
	if err != nil {
		logger.Fatal(err)
	}

	go bot.StartBot()

	srv.StartServer()
}
