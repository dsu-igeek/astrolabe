package kubernetes

import (
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
	"k8s.io/apimachinery/pkg/runtime"
)

type AstrolabeVolumeSnapshotterGetter struct {
	
}

func (a AstrolabeVolumeSnapshotterGetter) GetVolumeSnapshotter(name string) (velero.VolumeSnapshotter, error) {
	panic("implement me")
}

type AstrolabeVolumeSnapshotter struct {
	
}

func (a AstrolabeVolumeSnapshotter) Init(config map[string]string) error {
	panic("implement me")
}

func (a AstrolabeVolumeSnapshotter) CreateVolumeFromSnapshot(snapshotID, volumeType, volumeAZ string, iops *int64) (volumeID string, err error) {
	panic("implement me")
}

func (a AstrolabeVolumeSnapshotter) GetVolumeID(pv runtime.Unstructured) (string, error) {
	panic("implement me")
}

func (a AstrolabeVolumeSnapshotter) SetVolumeID(pv runtime.Unstructured, volumeID string) (runtime.Unstructured, error) {
	panic("implement me")
}

func (a AstrolabeVolumeSnapshotter) GetVolumeInfo(volumeID, volumeAZ string) (string, *int64, error) {
	panic("implement me")
}

func (a AstrolabeVolumeSnapshotter) CreateSnapshot(volumeID, volumeAZ string, tags map[string]string) (snapshotID string, err error) {
	panic("implement me")
}

func (a AstrolabeVolumeSnapshotter) DeleteSnapshot(snapshotID string) error {
	panic("implement me")
}
