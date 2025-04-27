package client

import (
	"context"
	"github.com/drinkthere/okx"
	"github.com/drinkthere/okx/api"
	"log"
	"syncPrice/config"
)

type OkxClient struct {
	Client *api.Client
}

func (okxClient *OkxClient) Init(cfg *config.Config) bool {
	var dest okx.Destination
	if cfg.OkxColo == "zoneB" {
		dest = okx.ColoServer
	} else if cfg.OkxColo == "zoneD" {
		dest = okx.ColoDServer
	} else {
		dest = okx.NormalServer
	}

	ctx := context.Background()

	var client *api.Client
	var err error
	if cfg.OkxLocalIP == "" {
		client, err = api.NewClient(ctx, cfg.OkxAPIKey, cfg.OkxSecretKey, cfg.OkxPassword, dest)
	} else {
		client, err = api.NewClientWithIP(ctx, cfg.OkxAPIKey, cfg.OkxSecretKey, cfg.OkxPassword, dest, cfg.OkxLocalIP)
	}
	if err != nil {
		log.Fatal(err)
		return false
	}

	okxClient.Client = client
	return true
}
