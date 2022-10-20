package main

import (
	"github.com/ALTA-Group-Project-Social-Media-Apps/Social-Media-Apps/config"
	"github.com/ALTA-Group-Project-Social-Media-Apps/Social-Media-Apps/features/user/delivery"
	"github.com/ALTA-Group-Project-Social-Media-Apps/Social-Media-Apps/features/user/repository"
	"github.com/ALTA-Group-Project-Social-Media-Apps/Social-Media-Apps/features/user/services"
	database "github.com/ALTA-Group-Project-Social-Media-Apps/Social-Media-Apps/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	cfg := config.NewConfig()
	db := database.InitDB(cfg)
	database.MigrateDB(db)
	uRepo := repository.New(db)
	uService := services.New(uRepo)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	delivery.New(e, uService)

	e.Logger.Fatal(e.Start(":8000"))
}
