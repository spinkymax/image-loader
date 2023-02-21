package main

import (
	"github.com/sirupsen/logrus"
	"image-loader/internal/repository"
	"image-loader/internal/server"
	"image-loader/internal/service"
)

func main() {
	userRepo := repository.NewUserRepo("myNewFile.json")

	controller := service.NewController(userRepo)

	loger := logrus.New()

	srv := server.NewServer(":8000", loger, controller)

	srv.RegisterRoutes()
	srv.StartRouter()

}
