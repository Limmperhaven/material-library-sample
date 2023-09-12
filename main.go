package main

import (
	"context"
	"git.miem.hse.ru/1206/app"
	"git.miem.hse.ru/1206/app/logger"
	"git.miem.hse.ru/1206/app/storage/stpg"
	"git.miem.hse.ru/1206/material-library/internal/client"
	"git.miem.hse.ru/1206/material-library/internal/config"
	"git.miem.hse.ru/1206/material-library/internal/domain"
	"git.miem.hse.ru/1206/material-library/internal/service"
	"log"
	"os"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("error initializating etc %s", err.Error())
	}

	lg := logger.Init(&cfg.Logger)

	jaegerTP, err := logger.NewJaeger(&cfg.Jaeger)
	if err != nil {
		lg.Warn(err)
	}

	if err := stpg.InitConnect(&cfg.Postgres); err != nil {
		lg.Fatal(err)
	}

	permsClient, err := client.NewPermissionsClient(&cfg.Permissions)
	if err != nil {
		lg.Fatal(err)
	}

	s3Client, err := client.NewS3Client(&cfg.S3)
	if err != nil {
		lg.Fatal(err)
	}

	educationClient, err := client.NewEducationClient(&cfg.Education)
	if err != nil {
		lg.Fatal(err)
	}

	usecase := domain.NewUsecase(permsClient, s3Client, educationClient)
	server, err := service.NewLibraryServer(&cfg.GRPC, usecase)
	if err != nil {
		lg.Fatal(err)
	}

	server.Run()
	lg.Info("started grpc server ", cfg.GRPC.Host, ":", cfg.GRPC.Port)

	app.Lock(make(chan os.Signal, 1))

	jaegerTP.Shutdown(context.Background())
}
