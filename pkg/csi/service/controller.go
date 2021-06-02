package service

import (
	"context"
	"fmt"

	"github.com/arturoguerra/truenas-csi/pkg/truenasapi"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
CreateVolume
Parameters:
	Name
	capacity_range
	volume_capabilities
	parameters
	secrets
	volume_content_source
	accessibility_requirements
Errors:
	InvalidArgument
	NotFound
	AlreadyExists
	ResourceExhuasted
	OutOfRange
	Internal
*/
func (s *service) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	s.CVMux.Lock()
	defer s.CVMux.Unlock()
	s.Logger.Info("CreateVolume")
	defer s.Logger.Info("CreateVolume")

	volname := req.GetName()

	params, err := s.parseParams(req.GetParameters())
	if err != nil {
		s.Logger.Error(err)
		return nil, status.Error(codes.InvalidArgument, "")
	}

	dataset := fmt.Sprintf("%s/%s", params.Dataset, volname)

	var volSize int64
	if req.GetCapacityRange() != nil && req.GetCapacityRange().RequiredBytes%1024 == 0 {
		volSize = int64(req.GetCapacityRange().RequiredBytes)
	} else {
		return nil, status.Error(codes.OutOfRange, "")
	}

	vol, err := s.TClient.GetVolume(dataset)
	if err != nil {
		if _, ok := err.(*truenasapi.NotFoundError); !ok {
			// Returns error if we get a 500 when getting a volume
			return nil, status.Error(codes.NotFound, err.Error())
		}
	}

	// Returns if volume already exists
	if vol != nil {
		volcontext := &VolumeContext{
			Name:    volname,
			Dataset: vol.ID,
			FSType:  params.FSType,
		}

		return &csi.CreateVolumeResponse{
			Volume: &csi.Volume{
				VolumeId:      volname,
				CapacityBytes: int64(vol.VolBlockSize.Parsed),
				VolumeContext: s.toMap(volcontext),
			},
		}, nil
	}

	volopts := truenasapi.VolumeOpts{
		Name:         dataset,
		Type:         "VOLUME",
		VolSize:      volSize,
		VolBlockSize: "512",
		Sparse:       false,
		Comments:     "TrueNAS CSI Driver",
		Compression:  params.Compression,
	}

	// Creates volume
	volume, err := s.TClient.CreateVolume(volopts)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	volcontext := &VolumeContext{
		Name:    volname,
		Dataset: volume.ID,
		FSType:  params.FSType,
	}

	mapped := s.toMap(volcontext)

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      volume.ID,
			CapacityBytes: int64(volume.VolBlockSize.Parsed),
			VolumeContext: mapped,
		},
	}, nil
}

/*
DeleteVolume
Parameters:
	volume_id
	secrets
Errors:
	FailedPreCondition
	Internal
*/
func (s *service) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	s.DVMux.Lock()
	defer s.DVMux.Unlock()
	s.Logger.Info("DeleteVolume")
	defer s.Logger.Info("DeleteVolume")

	// Full data set Ex: rabbitTank/test/testvol
	volID := req.GetVolumeId()

	//  TODO: Check is snapshots exist
	fmt.Printf("Deleing volume: %s\n", volID)

	s.TClient.DeleteVolume(volID)

	return &csi.DeleteVolumeResponse{}, nil

}

/*
Parameters:
	volume_id
	node_id
	volume_capability
	readonly
	secrets
	volume_context
Errors:
	NotFound
	AlreadyExists
	FailedPreCondition
	ResourceExhuasted
	Internal
*/
func (s *service) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	s.CPVMux.Lock()
	defer s.CPVMux.Unlock()
	s.Logger.Info("ControllerPublishVolume")
	defer s.Logger.Info("ControllerPublishVolume")

	// Full Volume ID Ex: rabbitTank/test/testvol
	volID := req.GetVolumeId()
	// Full Node ID Ex: exegol
	nodeID := req.GetNodeId()

	volctx, err := s.parseVolumeContext(req.GetVolumeContext())
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	// Gets the node definition
	node, err := s.TClient.GetNode(nodeID)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	// Gets volume definition
	volume, err := s.TClient.GetVolume(volID)
	if err != nil {
		switch e := err.(type) {
		case *truenasapi.NotFoundError:
			return nil, status.Error(codes.NotFound, e.Error())
		default:
			return nil, status.Error(codes.Internal, e.Error())
		}
	}

	// ISCSI ID
	iscsidevice, err := s.TClient.GetISCSIDevice(volume.ID)
	if err != nil {
		switch e := err.(type) {
		case *truenasapi.NotFoundError:
			// volume doesnt exists
			iscsiopts := truenasapi.ISCSIDeviceOpts{
				Name:      volctx.Name,
				Node:      node,
				Comment:   volctx.Name,
				Disk:      volume.ID,
				BlockSize: volume.VolBlockSize.Parsed,
				Enabled:   true,
				RO:        false,
				LunID:     1,
			}

			iscsidevice, err := s.TClient.CreateISCSIDevice(iscsiopts)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}

			pubctx := PublishContext{
				IQN:    iscsidevice.IQN,
				Portal: iscsidevice.Portal,
				LunID:  iscsidevice.LunID,
			}

			return &csi.ControllerPublishVolumeResponse{
				PublishContext: s.toMap(pubctx),
			}, nil
		default:
			return nil, status.Error(codes.Internal, e.Error())

		}

	}

	if iscsidevice.Node.ID != node.ID {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	pubctx := PublishContext{
		IQN:    iscsidevice.IQN,
		Portal: iscsidevice.Portal,
		LunID:  iscsidevice.LunID,
	}

	return &csi.ControllerPublishVolumeResponse{
		PublishContext: s.toMap(pubctx),
	}, nil

}

/*
ControllerUnpublishVolume
Parameters:
	volume_id
	node_id
	secrets
Errors:
	NotFound
	Internal
*/
func (s *service) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	volid := req.GetVolumeId()
	if err := s.TClient.DeleteISCSIDevice(volid); err != nil {
		switch err := err.(type) {
		case *truenasapi.NotFoundError:
			return &csi.ControllerUnpublishVolumeResponse{}, nil
		default:
			var _ = err // gets rid of compiler error
			return &csi.ControllerUnpublishVolumeResponse{}, nil
		}
	}

	return &csi.ControllerUnpublishVolumeResponse{}, nil
}

// ControllerGetCapabilities : Reports all the CSI Driver capabilities
func (s *service) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: []*csi.ControllerServiceCapability{
			{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
					},
				},
			},
			{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
					},
				},
			},
		},
	}, nil
}

// Unimplemented

// ListVolumes : returns a list of all known volumes
func (s *service) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	s.Logger.Info("Running ListVolumes")
	return nil, status.Error(codes.Unimplemented, "")
}

// ControllerExpandVolume : Expands a volume
func (s *service) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// CreateSnapshot : Creates a snapshot from a volume
func (s *service) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// DeleteSnapshot : Deletes a snapshot from a volume
func (s *service) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ListSnapshots : lists a known snapshots
func (s *service) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// GetCapacity : returns the total capacity of a volume
func (s *service) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

/*
// ControllerGetVolume : Alpha its just here and does nothing :)
func (s *service) ControllerGetVolume(ctx context.Context, req *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
*/

func (s *service) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
