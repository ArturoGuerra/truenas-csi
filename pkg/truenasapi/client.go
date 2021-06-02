package truenasapi

import (
	"errors"
	"fmt"
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
		// The stuff bellow is to reduce the amount of api calls required
		portal int
		// The stuff bellow is to reduce the amount of api calls required
		baseIQN string
		// portal ip
		portalAddr string
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

	global, err := c.getGlobal()
	if err != nil {
		return nil, err
	}

	portal, err := c.getPortal()
	if err != nil {
		return nil, err
	}

	if len(portal.Listen) != 1 {
		return nil, errors.New("invalid number portal addresses")
	} else if portal.Listen[0].IP == "0.0.0.0" {
		return nil, errors.New("invalid portal address")
	}

	c.portal = portal.ID
	c.portalAddr = fmt.Sprintf("%s:%d", portal.Listen[0].IP, portal.Listen[0].Port)
	c.baseIQN = global.Basename

	return c, nil
}

func NewDefault(logger *logrus.Logger) (Client, error) {
	config, err := NewDefaultConfig()
	if err != nil {
		return nil, err
	}

	return New(logger, config)
}
