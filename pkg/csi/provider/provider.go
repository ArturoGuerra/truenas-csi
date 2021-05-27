package provider

import (
	"github.com/arturoguerra/truenas-csi/pkg/csi/service"
	"github.com/rexray/gocsi"
)

func New(svc service.Service) gocsi.StoragePluginProvider {
	return &gocsi.StoragePlugin{
		Controller: svc,
		Identity:   svc,
		Node:       svc,
		EnvVars: []string{
			gocsi.EnvVarSerialVolAccess + "=true",
			gocsi.EnvVarSpecValidation + "=true",
			gocsi.EnvVarRequirePubContext + "=true",
		},
	}
}
