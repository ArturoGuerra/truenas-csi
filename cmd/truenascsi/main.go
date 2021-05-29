package main

import (
	"context"

	logging "github.com/arturoguerra/go-logging"
	"github.com/arturoguerra/truenas-csi/pkg/csi/service"
	"github.com/arturoguerra/truenas-csi/pkg/csi/provider"
	"github.com/arturoguerra/truenas-csi/pkg/truenasapi"
	"github.com/rexray/gocsi"
)

func main() {
	logger := logging.New()

	tclient, err := truenasapi.NewDefault(logger)
	if err != nil {
		logger.Panic(err)
	}

	logger.Info("Starting GoCSI for TrueNAS...")
	svc := service.New(logger, tclient)
	gocsi.Run(
		context.Background(),
		service.Name,
		"CSI Driver for TrueNAS",
		"",
		provider.New(svc),
	)

}
