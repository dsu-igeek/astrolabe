package ebs

import (
	"context"
	"github.com/vmware-tanzu/astrolabe/pkg/astrolabe"
	"github.com/aws/aws-sdk-go/service/ec2"

	"io"
)

type EBSProtectedEntity struct {
	epetm *EBSProtectedEntityTypeManager
	id astrolabe.ProtectedEntityID
	ec2 *ec2.EC2
}

func newEBSProtectedEntity(epetm *EBSProtectedEntityTypeManager, id astrolabe.ProtectedEntityID) (EBSProtectedEntity, error) {
	panic("implement me")
}

func (this EBSProtectedEntity) GetInfo(ctx context.Context) (astrolabe.ProtectedEntityInfo, error) {
	panic("implement me")
}

func (this EBSProtectedEntity) GetCombinedInfo(ctx context.Context) ([]astrolabe.ProtectedEntityInfo, error) {
	panic("implement me")
}

func (this EBSProtectedEntity) Snapshot(ctx context.Context) (astrolabe.ProtectedEntitySnapshotID, error) {
	csi := ec2.CreateSnapshotInput{
		Description:       nil,
		DryRun:            nil,
		TagSpecifications: nil,
		VolumeId:          nil,
	}
	this.ec2.CreateSnapshot(&csi)
	panic("implement me")
}

func (this EBSProtectedEntity) ListSnapshots(ctx context.Context) ([]astrolabe.ProtectedEntitySnapshotID, error) {
	panic("implement me")
}

func (this EBSProtectedEntity) DeleteSnapshot(ctx context.Context, snapshotToDelete astrolabe.ProtectedEntitySnapshotID) (bool, error) {
	panic("implement me")
}

func (this EBSProtectedEntity) GetInfoForSnapshot(ctx context.Context, snapshotID astrolabe.ProtectedEntitySnapshotID) (*astrolabe.ProtectedEntityInfo, error) {
	panic("implement me")
}

func (this EBSProtectedEntity) GetComponents(ctx context.Context) ([]astrolabe.ProtectedEntity, error) {
	panic("implement me")
}

func (this EBSProtectedEntity) GetID() astrolabe.ProtectedEntityID {
	panic("implement me")
}

func (this EBSProtectedEntity) GetDataReader(ctx context.Context) (io.ReadCloser, error) {
	panic("implement me")
}

func (this EBSProtectedEntity) GetMetadataReader(ctx context.Context) (io.ReadCloser, error) {
	panic("implement me")
}
