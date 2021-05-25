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

	if err = client.DeleteVolume("rabbitTank/test/pain4"); err != nil {
		fmt.Println(err)
	}

	volopts := truenasapi.VolumeOpts{
		Name:         "rabbitTank/test/pain4",
		Type:         "VOLUME",
		VolSize:      1024,
		VolBlockSize: "512",
	}

	volume, err := client.CreateVolume(volopts)
	if err != nil {
		fmt.Print(err)
	}

	fmt.Printf("%v\n", volume)
	iscsiopts := truenasapi.ISCSIDeviceOpts{
		Name:      "pain4",
		Alias:     "",
		Mode:      "ISCSI",
		Portal:    1,
		Initiator: 1,
		Disk:      "zvol/rabbitTank/test/pain4",
		Path:      "zvol/rabbitTank/test/pain4",
		BlockSize: 512,
		Comment:   "Test",
		Enabled:   true,
		LunID:     0,
		RO:        false,
		Type:      "DISK",
	}

	iscsidevice, err := client.CreateISCSIDevice(iscsiopts)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(iscsidevice)

	iscsidevice2, err := client.GetISCSIDevice(iscsidevice.TargetID, iscsidevice.ExtentID, iscsidevice.TargetExtentID)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(iscsidevice2)

	if err = client.DeleteISCSIDevice(iscsidevice.TargetID, iscsidevice.ExtentID, iscsidevice.TargetExtentID); err != nil {
		fmt.Println(err)
	}

}
