// Package containerd wraps continerd/containerd  providing a minimal API layer
// halcyon only requires image distribution, snapshotting functionality, and container create/exec/delete
package containerd

import (
	"context"
	"errors"
	"log"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/platforms"

	"github.com/opencontainers/image-spec/identity"
)

// TODO: Create a high-level driver package which unifies containerd & libcontainer

// Docker Image repositories
const (
	BuilderRepo  = "halcyo/builder"
	EngineRepo   = "halcyo/engine"
	ResolverRepo = "halcyo/%s-resolver"
	RuntimeRepo  = "halcyo/%s-runtime"
	WorkerRepo   = "halcyo/worker"

	// RunForeverCmd a command which keeps a container running indefinitley
	ContainerdSock   = "/run/containerd/containerd.sock"
	DefaultNamespace = "halcyon"
	RunForeverCmd    = "tail -f /dev/null"
)

// Containerd provides a convenient client API over containerd.Client
type Containerd struct {
	*containerd.Client
}

// New instantiates a new Containerd instance
func New() (*Containerd, error) {
	client, err := containerd.New(ContainerdSock, containerd.WithDefaultNamespace(DefaultNamespace))
	if err != nil {
		return nil, err
	}

	// init Snapshot service
	client.SnapshotService(containerd.DefaultSnapshotter)

	return &Containerd{
		client,
	}, err
}

// Pull pulls an image with id
func (c *Containerd) Pull(image string) (containerd.Image, error) {
	ctx := context.Background()
	return c.Client.Pull(ctx, image, containerd.WithPullUnpack)
}

// GetImage returns an image as images.Image
func (c *Containerd) GetImage(image string) (images.Image, error) {
	ctx := context.Background()
	return c.Client.ImageService().Get(ctx, image)
}

// NewSnapshot creates a snapshot RootFS from an image
func (c *Containerd) NewSnapshot(id string, i images.Image) ([]mount.Mount, error) {
	ctx := context.Background()
	// get diff IDs from image
	diffs, err := i.RootFS(ctx, c.ContentStore(), platforms.Default())
	if err != nil {
		return nil, err
	}
	parent := identity.ChainID(diffs).String()
	return c.SnapshotService(containerd.DefaultSnapshotter).Prepare(ctx, id, parent)
}

// DeleteSnapshot creates a snapshot RootFS from an image
func (c *Containerd) DeleteSnapshot(id string) error {
	ctx := context.Background()
	return c.SnapshotService(containerd.DefaultSnapshotter).Remove(ctx, id)
}

// MountSnapshot mounts the snapshot files at the given rootfs path
func (c *Containerd) MountSnapshot(rootfs string, mounts []mount.Mount) (err error) {

	// Clean up if mounting fails
	defer func() {
		if err != nil {
			if err2 := c.UnmountSnapshot(rootfs); err != nil {
				log.Printf("MountSnapshot failed clean up: %s\n", err2)
			}
		}
	}()

	// Mount layers as RootFS at rootfs
	for _, m := range mounts {
		if err = m.Mount(rootfs); err != nil {
			return errors.New("failed to mount rootfs component")
		}
	}

	return
}

// UnmountSnapshot unmounts all snapshots from rootfs
func (c *Containerd) UnmountSnapshot(rootfs string) error {
	return mount.UnmountAll(rootfs, 0)
}

// NewContainer creates a new container
func (c *Containerd) NewContainer(id string, image containerd.Image) (*containerd.Task, error) {
	ctx := context.Background()

	con, err := c.Client.NewContainer(ctx, id,
		containerd.WithNewSnapshot(image.Name()+"-rootfs", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)

	// create a new task
	task, err := con.NewTask(ctx, cio.NewCreator(cio.WithStdio))

	// the task is now running and has a pid that can be use to setup networking
	// or other runtime settings outside of containerd
	// pid := task.Pid()

	// Let clients clean up Task
	// defer task.Delete(ctx)

	// wait for the task to exit and get the exit status
	// status, err := task.Wait(ctx)

	// start the process inside the container
	err = task.Start(ctx)
	return &task, err
}
