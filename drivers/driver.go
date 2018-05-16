// Package drivers provides halcyon client APIs with cherry-picked functionality over container drivers
package drivers

import (
	"github.com/holygits/halcyon-engine/drivers/containerd"
	"github.com/holygits/halcyon-engine/drivers/libcontainer"
	"github.com/holygits/halcyon-engine/structs"

	"go.uber.org/zap"
)

// TODO: create a worker pool of containers and load balance execution

const (
	// TODO: Implement with option funcs
	ContainerdSock   = "/run/containerd/containerd.sock"
	DefaultNamespace = "halcyon"
	LibContainerRoot = "/tmp/halcyon"
)

type Driver struct {
	libcontainer *libcontainer.Driver
	containerd   *containerd.Driver
	logger       *zap.Logger
}

// NewDriver returns a new ContainerDriver
func NewDriver(logger *zap.Logger) (*Driver, error) {
	// init logger
	//log, _ = zap.NewProduction()
	//defer log.Sync()

	l, err := libcontainer.NewDriver(LibContainerRoot, logger)
	if err != nil {
		return nil, err
	}

	c, err := containerd.NewDriver(DefaultNamespace, ContainerdSock, logger)
	if err != nil {
		return nil, err
	}

	return &Driver{
		libcontainer: l,
		containerd:   c,
	}, nil

}

// Bootstrap performs a one-time initalisation of function resources, call prior to Exec
func (d *Driver) Bootstrap(f *structs.FuncExecution) error {

	bID := "TODO: get from f"
	if err := d.libcontainer.NewBundle(bID, nil); err != nil {
		return err
	}
	image, err := d.containerd.PullImage(bID)
	if err != nil {
		return err
	}
	mounts, err := d.containerd.NewViewSnapshot(bID, image)
	if err != nil {
		return err
	}
	if err = d.containerd.MountSnapshot(d.libcontainer.BundleRootFS(bID), mounts); err != nil {
		return err
	}
	return d.libcontainer.NewContainer(string(f.ID), bID)
}

// Destory removes all function resources
func (d *Driver) Destory(f *structs.FuncExecution) error {
	bID := "TODO: get from f"
	if err := d.containerd.DestroySnapshot(bID); err != nil {
		d.logger.Error("destroy snapshot",
			zap.String("bundle", bID),
			zap.String("id", string(f.ID)),
			zap.Error(err))
	}
	if err = d.libcontainer.DestroyBundle(bID); err != nil {
		d.logger.Error("destroy bundle",
			zap.String("bundle", bID),
			zap.String("id", string(f.ID)),
			zap.Error(err))
	}
	if err = d.libcontainer.DestroyContainer(string(f.ID)); err != nil {
		d.logger.Error("destroy container",
			zap.String("bundle", bID),
			zap.String("id", string(f.ID)),
			zap.Error(err))
	}
	return err
}

// Exec runs the function
func (d *Driver) Exec(f *structs.FuncExecution) {
	d.libcontainer.Exec(id, bundle, cmd, ctx)
}
