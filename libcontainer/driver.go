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

var log *zap.Logger

const (
	configJSON = "config.json"
	procJSON   = "process.json"
	rootFS     = "rootfs"
	// DefaultRootDir is the driver root path
	DefaultRootDir  = "/tmp/halcyon"
	bootstrapBinary = "/usr/bin/bootstrap"
)

func init() {
	log, _ = zap.NewProduction()
	defer log.Sync()
}

// Driver provides a minimal client API over libcontainer.Factory and helpers
type Driver struct {
	bundleDir    string
	containerDir string
	libcontainer.Factory
}

// NewDriver initialises a new libcontainer client driver
func NewDriver(rootDir string) (*Driver, error) {

	if len(rootDir) == 0 {
		rootDir = DefaultRootDir
	}
	//os.RemoveAll(rootDir)

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
		factory,
	}, nil
}

// Bootstrap setup initial container config
func (d *Driver) Bootstrap(id string, env, entrypoint []string) error {
	// Bootstrap container bundle
	_, _, err := d.createBundle(id, env, entrypoint)
	return err
}

// Create creates a new container and bundle
func (d *Driver) Create(id string) error {

	// Load bootstrapped config.json
	conf, err := d.loadConfigJSON(id)
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

// Delete container with id
func (d *Driver) Delete(id string) (err error) {
	c, err := d.Load(id)
	if err != nil {
		return
	}

	c.Destroy()
	os.RemoveAll(filepath.Join(d.containerDir, id))
	os.RemoveAll(d.BundleDir(id))
	return
}

// Exec runs a command in container
func (d *Driver) Exec(id string, cmd []byte) (stdout, stderr []byte, err error) {

	// Load
	var c libcontainer.Container
	if c, err = d.Factory.Load(id); err != nil {
		return
	}

	// Pipes
	rErr, wErr := io.Pipe()
	rOut, wOut := io.Pipe()
	rIn, wIn := io.Pipe()

	go func() {
		// Write cmd to container's stdin
		wIn.Write(cmd)
		defer wIn.Close()
	}()

	// Process
	var p *libcontainer.Process
	if p, err = d.loadProcessJSON(d.BundleDir(id)); err != nil {
		return
	}
	p.Stdin = rIn
	p.Stdout = wOut
	p.Stderr = wErr

	// Run
	go func() {
		if err = c.Start(p); err != nil {
			log.Error("container exec", zap.String("id", id), zap.Error(err))
			return
		}
		p.Wait()
		defer rIn.Close()
		defer wErr.Close()
		defer wOut.Close()
	}()

	// Read stderr result
	defer rErr.Close()
	stderr, err = ioutil.ReadAll(rErr)
	if err != nil {
		return
	}

	// Read stdout result
	defer rOut.Close()
	stdout, err = ioutil.ReadAll(rOut)

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

// createBundle bootstraps a container bundle
func (d *Driver) createBundle(id string, env, entrypoint []string) (*config, *process, error) {

	// Recreate bundle dir
	//err := os.RemoveAll(d.BundleDir(id))
	// if err != nil {
	// 	if !os.IsNotExist(err) {
	// 		return nil, nil, err
	// 	}
	// }
	err := os.MkdirAll(d.BundleRootFS(id), 0755)
	if err != nil {
		if !os.IsExist(err) {
			return nil, nil, err
		}
	}

	conf, proc := d.generateConfig(id, env, entrypoint)

	// Write process.json
	buf, err := proc.MarshalJSON()
	if err != nil {
		return conf, proc, err
	}
	if err = ioutil.WriteFile(filepath.Join(d.BundleDir(id), procJSON), buf, 0644); err != nil {
		return conf, proc, err
	}

	// Write config.json
	buf, err = conf.MarshalJSON()
	if err != nil {
		return conf, proc, err
	}
	err = ioutil.WriteFile(filepath.Join(d.BundleDir(id), configJSON), buf, 0644)
	return conf, proc, err
}

// Generate an OpenContainer config.json and custom process.json
func (d *Driver) generateConfig(id string, env, entrypoint []string) (*config, *process) {

	log.Info("Generating config", zap.String("bundle", d.BundleDir(id)))

	// Create libcontainer config
	conf := defaultConfig(id, d.BundleRootFS(id))

	// Create libcontainer process
	p := &process{
		Args: entrypoint,
		User: "root",
	}

	// Set Func ENVs
	p.Env = append(p.Env, os.Environ()...) // System
	p.Env = append(p.Env, env...)          // User

	return conf, p
}

func (d *Driver) loadConfigJSON(id string) (*configs.Config, error) {

	buf, err := ioutil.ReadFile(filepath.Join(d.BundleDir(id), configJSON))
	if err != nil {
		return nil, err
	}

	c := &config{}
	c.UnmarshalJSON(buf)
	// Convert internal config type to libcontainer type
	return c.asConfig(), nil
}

func (d *Driver) loadProcessJSON(id string) (*libcontainer.Process, error) {

	buf, err := ioutil.ReadFile(filepath.Join(d.BundleDir(id), procJSON))
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
