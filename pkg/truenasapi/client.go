package truenasapi

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type (
	/* Client (TrueNAS API Endpoints wrapper) */
	Client interface {
		/* Volume Management */
		GetVolume(vol string) (*Volume, error)
		CreateVolume(vol VolumeOpts) (*Volume, error)
		DeleteVolume(vol string) error

		/* ISCSI Management */
		GetISCSIDevice(id string) (*ISCSIDevice, error)
		CreateISCSIDevice(dev ISCSIDeviceOpts) (*ISCSIDevice, error)
		DeleteISCSIDevice(id string) error

		GetISCSIID(volpath string) string

		/* Node Management */
		GetNode(node string) (*Node, error)
		SetNode(node Node) error
	}

	/* Implements the following functions:
	Public:
	GetVolume
	CreateVolume
	DeleteVolume

	GetISCSIDevice
	CreateISCSIDevice
	DeleteISCSIDevice

	GetNode
	SetNode

	Private:

	*/
	client struct {
		Logger     *logrus.Logger
		HttpClient *http.Client
		Config     *Config
		Validate   *validator.Validate
	}
)

func New(logger *logrus.Logger, cfg *Config) (Client, error) {
	validate := validator.New()
	httpClient := newHttpClient(cfg.APIKey)

	c := &client{
		Logger:     logger,
		HttpClient: httpClient,
		Config:     cfg,
		Validate:   validate,
	}

	return c, nil
}

func NewDefault(logger *logrus.Logger) (Client, error) {
	config, err := NewDefaultConfig()
	if err != nil {
		return nil, err
	}

	return New(logger, config)
}
