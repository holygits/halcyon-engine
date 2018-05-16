// Package containerd wraps continerd/containerd providing a minimal API layer with cherry-picked snapshotting functionality
package containerd

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/platforms"

	"github.com/opencontainers/image-spec/identity"
)

// Driver provides a minimal client API over containerd.Client
type Driver struct {
	logger *zap.Logger
	*containerd.Client
}

// NewDriver instantiates a new Driver
func NewDriver(namespace, sock string, logger *zap.Logger) (*Driver, error) {
	client, err := containerd.New(sock, containerd.WithDefaultNamespace(namespace))
	if err != nil {
		return nil, err
	}

	// init Snapshot service
	client.SnapshotService(containerd.DefaultSnapshotter)

	return &Driver{
		logger,
		client,
	}, err
}

// PullImage pulls an image using docker uri e.g. docker.io/library/redis:alpine
func (d *Driver) PullImage(image string) (containerd.Image, error) {
	ctx := context.Background()
	return d.Client.Pull(ctx, image, containerd.WithPullUnpack)
}

// extractImageParent converts containerd.Image to parent layer ID string
func (d *Driver) extractImageParent(ctx context.Context, image containerd.Image) (parent string, err error) {
	// implementation details of containerd
	i, err := d.Client.ImageService().Get(ctx, image.Name())
	if err != nil {
		return
	}
	// get diff IDs from image
	diffs, err := i.RootFS(ctx, d.ContentStore(), platforms.Default())
	if err != nil {
		return
	}
	parent = identity.ChainID(diffs).String()
	return
}

// NewSnapshot creates a snapshot RootFS from an image
func (d *Driver) NewSnapshot(id string, image containerd.Image) ([]mount.Mount, error) {
	ctx := context.Background()
	parent, err := d.extractImageParent(ctx, image)
	if err != nil {
		return nil, err
	}
	return d.SnapshotService(containerd.DefaultSnapshotter).Prepare(ctx, id, parent)
}

// NewViewSnapshot creates a new read-only snapshot RootFS from an image
func (d *Driver) NewViewSnapshot(id string, image containerd.Image) ([]mount.Mount, error) {
	ctx := context.Background()
	parent, err := d.extractImageParent(ctx, image)
	if err != nil {
		return nil, err
	}
	return d.SnapshotService(containerd.DefaultSnapshotter).View(ctx, id, parent)
}

// DestroySnapshot removes a snapshot
func (d *Driver) DestroySnapshot(id string) error {
	ctx := context.Background()
	return d.SnapshotService(containerd.DefaultSnapshotter).Remove(ctx, id)
}

// MountSnapshot mounts the snapshot files at the given target path
func (d *Driver) MountSnapshot(target string, mounts []mount.Mount) (err error) {
	// Clean up if mounting fails
	defer func() {
		if err != nil {
			if err2 := d.UnmountSnapshot(target); err2 != nil {
				d.logger.Error("Unmounting snapshot(s) failed", zap.Error(err2))
			}
		}
	}()

	// Mount layers as RootFS at target
	for _, m := range mounts {
		if err = m.Mount(target); err != nil {
			return errors.New("Failed to mount target component")
		}
	}

	return
}

// UnmountSnapshot unmounts all snapshots from target
func (d *Driver) UnmountSnapshot(target string) error {
	return mount.UnmountAll(target, 0)
}
