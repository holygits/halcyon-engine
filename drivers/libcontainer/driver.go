// Package libcontainer provides a minimal API over runc/libcontainer
// Halcyon only requires container create, exec, delete operations
package libcontainer

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/configs"

	"go.uber.org/zap"
)

const (
	configJSON      = "config.json"
	procJSON        = "process.json"
	rootFS          = "rootfs"
	bootstrapBinary = "/usr/bin/bootstrap"
)

// Driver provides a minimal client API over libcontainer.Factory and helpers
type Driver struct {
	bundleDir    string
	containerDir string
	logger       *zap.Logger
	libcontainer.Factory
}

// NewDriver initialises a new libcontainer client driver
func NewDriver(rootDir string, logger *zap.Logger) (*Driver, error) {

	// The dir for container "images"
	bundleDir := filepath.Join(rootDir, "bundle")
	if err := os.MkdirAll(bundleDir, 0755); err != nil {
		return nil, err
	}

	// The tmpfs dir for created/running containers
	containerDir := filepath.Join(rootDir, "run")
	if err := os.MkdirAll(containerDir, 0755); err != nil {
		return nil, err
	}

	// Init libcontainer.Factory
	factory, err := libcontainer.New(
		containerDir,
		libcontainer.InitArgs(bootstrapBinary, "init"),
	)
	if err != nil {
		return nil, err
	}

	return &Driver{
		bundleDir,
		containerDir,
		logger,
		factory,
	}, nil
}

// NewBundle sets up an initial container config (bundle)
func (d *Driver) NewBundle(id string, entrypoint []string) error {
	d.logger.Debug("new bundle", zap.String("id", id))
	err := os.MkdirAll(d.BundleRootFS(id), 0755)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	conf, proc := d.generateConfig(id, []string{"PATH=/bin"}, entrypoint)

	// Write process.json
	buf, err := proc.MarshalJSON()
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filepath.Join(d.BundleDir(id), procJSON), buf, 0644); err != nil {
		return err
	}

	// Write config.json
	buf, err = conf.MarshalJSON()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(d.BundleDir(id), configJSON), buf, 0644)
}

// NewContainer creates a new container with id using bundle
func (d *Driver) NewContainer(id, bundle string) error {
	d.logger.Debug("create container", zap.String("id", id), zap.String("bundle", bundle))
	// Load bootstrapped config.json
	conf, err := d.loadConfigJSON(bundle)
	if err != nil {
		return err
	}

	// Create container
	c, err := d.Factory.Create(id, conf)
	if err != nil {
		return err
	}

	p := &libcontainer.Process{
		Args: []string{"echo", "init"},
		Env:  []string{"PATH=/bin"},
	}

	// Init
	return c.Start(p)
}

// DestroyBundle cleans up bundle with id
func (d *Driver) DestroyBundle(id string) (err error) {
	d.logger.Debug("Delete bundle", zap.String("bundle", d.BundleDir(id)))
	if err = os.RemoveAll(d.BundleDir(id)); err != nil {
		d.logger.Error("conatiner Destroy()", zap.Error(err))
	}
	if err = os.RemoveAll(filepath.Join(d.containerDir, id)); err != nil {
		d.logger.Error("conatiner Destroy()", zap.Error(err))
	}
	return
}

// DestroyContainer deletes container with id
func (d *Driver) DestroyContainer(id string) (err error) {
	d.logger.Debug("delete container", zap.String("id", id))
	c, err := d.Load(id)
	if err != nil {
		return
	}
	return c.Destroy()
}

// Exec runs a command in container with id and cmd
func (d *Driver) Exec(id, bundle string, cmd []string, ctx []byte) (stdout, stderr []byte, err error) {
	d.logger.Debug("exec container", zap.String("id", id), zap.String("bundle", bundle))
	// Load
	var c libcontainer.Container
	if c, err = d.Factory.Load(id); err != nil {
		return
	}

	// Stdio pipes
	rErr, wErr := io.Pipe()
	rOut, wOut := io.Pipe()
	rIn, wIn := io.Pipe()

	go func() {
		// Write ctx to container's stdin
		wIn.Write(ctx)
		defer wIn.Close()
	}()

	// Process
	var p *libcontainer.Process
	if p, err = d.loadProcessJSON(bundle); err != nil {
		return
	}
	p.Args = append(p.Args, cmd...)
	p.Stdin = rIn
	p.Stderr = wErr
	p.Stdout = wOut
	// TODO: In kubernetes env will populate with useful cluster service address etc.
	// p.Env = append(p.Env, os.Environ())

	// Run
	go func() {
		if err = c.Start(p); err != nil {
			d.logger.Error("container exec", zap.String("id", id), zap.String("bundle", id), zap.Error(err))
			return
		}
		p.Wait()
		// Clean up stdio pipes
		rIn.Close()
		wOut.Close()
		wErr.Close()
	}()

	// Read stdout result
	defer rOut.Close()
	stdout, err = ioutil.ReadAll(rOut)
	if err != nil {
		return
	}

	// Read stderr result
	defer rErr.Close()
	stderr, err = ioutil.ReadAll(rErr)
	return
}

// BundleDir returns a path to a bundle by id
func (d *Driver) BundleDir(id string) string {
	return filepath.Join(d.bundleDir, id)
}

// BundleRootFS returns a path to a bundle's RootFS by id
func (d *Driver) BundleRootFS(id string) string {
	return filepath.Join(d.BundleDir(id), rootFS)
}

// generateConfig generates an OpenContainer config.json and custom process.json
func (d *Driver) generateConfig(bundle string, env, entrypoint []string) (*config, *process) {
	d.logger.Debug("Generating config", zap.String("bundle", d.BundleDir(bundle)))

	// Create libcontainer config
	conf := defaultConfig(bundle, d.BundleRootFS(bundle))

	// Create libcontainer process
	return conf, &process{
		Args: entrypoint,
		User: "root",
		Env:  env,
	}
}

func (d *Driver) loadConfigJSON(bundle string) (*configs.Config, error) {

	buf, err := ioutil.ReadFile(filepath.Join(d.BundleDir(bundle), configJSON))
	if err != nil {
		return nil, err
	}

	c := &config{}
	c.UnmarshalJSON(buf)
	// Convert internal config type to libcontainer type
	return c.asConfig(), nil
}

func (d *Driver) loadProcessJSON(bundle string) (*libcontainer.Process, error) {

	buf, err := ioutil.ReadFile(filepath.Join(d.BundleDir(bundle), procJSON))
	if err != nil {
		return nil, err
	}

	p := &process{}
	if err = p.UnmarshalJSON(buf); err != nil {
		return nil, err
	}

	// Convert internal process type to libcontainer type
	return &libcontainer.Process{
		Args: p.Args,
		Env:  p.Env,
		User: p.User,
	}, nil
}
