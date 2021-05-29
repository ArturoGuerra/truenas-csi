package service

import (
	"os"
	"strings"
	"sync"

	"github.com/arturoguerra/truenas-csi/pkg/mounter"
	"github.com/arturoguerra/truenas-csi/pkg/truenasapi"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/go-playground/validator/v10"
	"github.com/rexray/gocsi"
	"github.com/sirupsen/logrus"
)

const (
	// Name : CSI Driver name
	Name = "csi.truenas.ar2ro.io"
	// VendorVersion : CSI version
	VendorVersion = "1.0.0"
	// UnixSocketPrefix : CSI Driver unix prefix
	UnixSocketPrefix = "unix://"
)

// Manifest : CSI Driver manifest
var Manifest = map[string]string{
	"url": "https://github.com/arturoguerra/truenas-csi",
}

type (
	// Service : combines all the CSI Driver spec
	Service interface {
		csi.ControllerServer
		csi.IdentityServer
		csi.NodeServer
	}

	service struct {
		NodeID    string
		ClusterID string
		TClient   truenasapi.Client
		Logger    *logrus.Logger
		Validator *validator.Validate
		Mounter   mounter.Mounter

		// CreateVolume
		CVMux sync.Mutex
		// DeleteVolume
		DVMux sync.Mutex
		// ControllerPublushVolume
		CPVMux sync.Mutex
		// ControllerUnpublishVolume
		CUVMux sync.Mutex
		// NodeStageVolume
		NSVMux sync.Mutex
		// NodeUnstageVolume
		NUSVMux sync.Mutex
		// NodePublishVolume
		NPVMux sync.Mutex
		// NodeUnpublishVolume
		NUPVMux sync.Mutex
	}
)

// Work around if node dies and old csi.sock is left behind.
// NOTE: This can cause issues if two instances of a node are scheduled in the same node but that would be an extreme edge case.
func init() {
	sockPath := os.Getenv(gocsi.EnvVarEndpoint)
	sockPath = strings.TrimPrefix(sockPath, UnixSocketPrefix)
	if len(sockPath) > 1 {
		os.Remove(sockPath)
	}
}

// New : returns a new instance of the csi driver
func New(logger *logrus.Logger, tnapi truenasapi.Client) Service {
	return &service{
		NodeID:    os.Getenv("NODE_ID"),
		ClusterID: os.Getenv("CLUSTER_ID"),
		TClient:   tnapi,
		Logger:    logger,
		Validator: validator.New(),
		Mounter:   mounter.New(logger),
	}
}
