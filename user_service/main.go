package main

import (
	"context"
	"log"
	"user_service/internal/config"
	"user_service/internal/domain"
	"user_service/internal/handler"
	"user_service/internal/httpapi"
	"user_service/internal/repository"
	"user_service/internal/service"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()

	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}
	db.AutoMigrate(&domain.User{})

	mc, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal("Can't connect to MinIO:", err)
	}

	ctx := context.Background()
	exists, _ := mc.BucketExists(ctx, "avatars")
	if !exists {
		mc.MakeBucket(ctx, "avatars", minio.MakeBucketOptions{})
	}

	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo, cfg, mc)
	h := handler.NewUserHandler(svc)
	e := httpapi.NewRouter(h)

	log.Printf("User Service is running at %s", cfg.HTTPAddress)
	e.Start(cfg.HTTPAddress)
}
