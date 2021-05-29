package service

import (
	"context"

	"github.com/arturoguerra/truenas-csi/pkg/truenasapi"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-lib-iscsi/iscsi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	s.Logger.Info("Getting NodeInfo for")

	topology := new(csi.Topology)

	return &csi.NodeGetInfoResponse{
		NodeId:             s.NodeID,
		AccessibleTopology: topology,
	}, nil
}

/*
NodeStageVolume
Parameters:
	volume_id
	publish_content
	staging_target_path
	volume_capability
	secrets
	volume_context
Errors:
	NotFound
	AlreadyExists
	FailedPreCondition
	Internal
*/
func (s *service) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	//olID := req.GetVolumeId()
	stagingTargetPath := req.GetStagingTargetPath()

	publishContext, err := s.parsePublishContext(req.GetPublishContext())
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	volumeContext, err := s.parseVolumeContext(req.GetVolumeContext())
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	conn := iscsi.Connector{
		TargetIqn:     publishContext.IQN,
		TargetPortals: []string{publishContext.Portal},
		Lun:           int32(publishContext.LunID),
		AuthType:      "none",
		Multipath:     false,
		RetryCount:    11,
		CheckInterval: 1,
	}

	// Connects iscsi device to node
	device, err := iscsi.Connect(conn)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	// Mouting ISCSI Drive to staging path
	opts := []string{}
	if err := s.Mounter.Mount(stagingTargetPath, device, volumeContext.FSType, opts); err != nil {
		s.Logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

/*
NodeUnstageVolume
Parameters:
	volume_id
	staging_target_path
Errors:
	NotFound
	Internal
*/
func (s *service) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	volID := req.GetVolumeId()
	stagingTargetPath := req.GetStagingTargetPath()

	if err := s.Mounter.Unmount(stagingTargetPath); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	iscsid := s.TClient.GetISCSIID(volID)
	iscsidevice, err := s.TClient.GetISCSIDevice(iscsid)
	if err != nil {
		if e, ok := err.(*truenasapi.NotFoundError); !ok {
			return nil, status.Error(codes.Internal, e.Error())
		}

		return nil, status.Error(codes.NotFound, "ISCSI Device not found")
	}

	iscsi.Disconnect(iscsidevice.IQN, []string{iscsidevice.Portal})
	return nil, nil
}

/*
NodePublishVolume Bind Mounts from staging to pod mount path
Parameters:
	volume_id
	publish_context
	staging_target_path
	target_path
	volume_capability
	readonly
	secrets
	volume_context
Errors:
	NotFound
	AlreadyExists
	FailedPreCondition
	Internal
*/
func (s *service) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	s.Logger.Info("Running NodePublishVolume")
	stagingTargetPath := req.GetStagingTargetPath()
	targetPath := req.GetTargetPath()

	/* Check if target is a path and creates it if its not */
	if pathExists, err := s.Mounter.PathExists(targetPath); !pathExists {
		if err = s.Mounter.MakeDir(targetPath); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	opts := []string{
		"bind",
	}

	if err := s.Mounter.Mount(stagingTargetPath, targetPath, "auto", opts); err != nil {
		s.Logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	/*
		opts := s.Mounter.BindMount(stagingTargetPath, targetPath, "auto", opts); err != nil {
			s.Logger.Error(err)
			return nil, status.Error(codes.Internal, err.Error())
		}
	*/

	s.Logger.Infof("Bind Mounted Volume: (%s) to Path: (%s)", stagingTargetPath, targetPath)
	return &csi.NodePublishVolumeResponse{}, nil
}

/*
NodeUnpublishVolume unmounts the bind mount between the pod and common directory
Parameters:
	volume_id
	node_id
	secrets
Errors:
	NotFound
	Internal
*/
func (s *service) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	s.Logger.Info("Running NodeUnpublishVolume")
	targetPath := req.GetTargetPath()
	if len(targetPath) == 0 {
		return nil, status.Error(codes.InvalidArgument, "NodeUnpublishVolume Target Path must be provided")
	}

	s.Logger.Infof("Unmounting: %s", targetPath)

	if err := s.Mounter.Unmount(targetPath); err != nil {
		s.Logger.Error(err)
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (s *service) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{
			{
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
					},
				},
			},
		},
	}, nil
}

func (s *service) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (s *service) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
