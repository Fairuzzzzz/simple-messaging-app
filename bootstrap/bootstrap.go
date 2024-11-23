package bootstrap

import (
	"io"
	"log"
	"os"

	"github.com/Fairuzzzzz/fiber-boostrap/app/ws"
	"github.com/Fairuzzzzz/fiber-boostrap/pkg/database"
	"github.com/Fairuzzzzz/fiber-boostrap/pkg/env"
	"github.com/Fairuzzzzz/fiber-boostrap/pkg/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func NewApplication() *fiber.App {
	env.SetupEnvFile()
	SetupLogfile()
	database.SetupDatabase()
	database.SetupMongoDB()
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Get("/dashboard", monitor.New())
	router.InstallRouter(app)

	go ws.ServeWSMessaging(app)

	return app
}

func SetupLogfile() {
	logfile, err := os.OpenFile(
		"./logs/simple_messaging_app.log",
		os.O_CREATE|os.O_APPEND|os.O_RDWR,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(mw)
}
