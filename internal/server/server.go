package server

import (
	"github.com/gofiber/fiber/v2"

	"ussd_eth_v2/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "ussd_eth_v2",
			AppName:      "ussd_eth_v2",
		}),

		db: database.New(),
	}

	return server
}
