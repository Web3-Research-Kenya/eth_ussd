package handlers

import (
	"ussd_eth_v2/internal/database"
	"ussd_eth_v2/internal/eth"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Handler struct {
	DB     database.Service
	Tree   *MenuTree
	Dat    Data
	Client *ethclient.Client
}

func NewHandler(db database.Service) *Handler {
	return &Handler{
		Tree:   NewMenuTree(),
		DB:     db,
		Dat:    Data{},
		Client: eth.Connect(),
	}
}
