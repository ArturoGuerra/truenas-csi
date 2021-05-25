package main

import (
	"fmt"

	logging "github.com/arturoguerra/go-logging"
	"github.com/arturoguerra/truenas-csi/pkg/truenasapi"
)

func main() {
	logger := logging.New()

	config, err := truenasapi.NewDefaultConfig()
	if err != nil {
		panic(err)
	}
	client, err := truenasapi.New(logger, config)
	if err != nil {
		panic(err)
	}

	if err = client.DeleteVolume("rabbitTank/test/pain2"); err != nil {
		fmt.Println(err)
	}

	/*
		volopts := truenasapi.VolumeOpts{
			Name:         "rabbitTank/test/pain2",
			Type:         "VOLUME",
			VolSize:      1024,
			VolBlockSize: "512",
		}

		volume, err := client.CreateVolume(volopts)
		if err != nil {
			fmt.Print(err)
		}

		fmt.Printf("%v\n", volume)
	*/
}
