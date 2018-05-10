// Package containerd wraps continerd/containerd providing an API layer
package containerd

import (
  "context"
  "io"

	"github.com/containerd/containerd"
	// "github.com/containerd/containerd/cio"
)

const (
	// Docker Repositoriesß
	BuilderRepo  = "halcyo/builder"
	EngineRepo   = "halcyo/engine"
	ResolverRepo = "halcyo/%s-resolver"
  RuntimeRepo =ß "halcyo/%s-runtime"
	WorkerRepo   = "halcyo/worker"

	// Commands
	RunForeverCmd = "tail -f /dev/null"
)

// Client provides a halcyon API over containerd
type Containerd struct {
  context *context.Context
  *containerd.Containerd
}

// NewContanerd instantiates a new Contanerd instance
func NewContanerd() (*Containerd, error) {
	client, err := containerd.New("/run/containerd/containerd.sock",
                                containerd.WithDefaultNamespace("halcyon"))
	return &Containerd{
    context: context.Background(),
		client,
	}, err
}

// Exec runs a new command in a container
func (c *Containerd) Exec(id string, cmd []string) (stdout, stderr []byte, err error) {

    cli, ctx, cancel, err := commands.NewClient(context)
    if err != nil {
      return
    }
    defer cancel()

    container, err := cli.LoadContainer(ctx, id)
    if err != nil {
      return
    }

    spec, err := container.Spec(ctx)
    if err != nil {
      return
    }

    task, err := container.Task(ctx, nil)
    if err != nil {
      return
    }

    pspec := spec.Process
    pspec.Terminal = false
    pspec.Args = cmd

    // Create Stdin/out/err streams
    rStdout, wStdout := io.Pipe()
    defer rStdout.Close()
    defer wStdout.Close()

    rStdin, wStdin := io.Pipe()
    defer rStdin.Close()
    defer wStdin.Close()

    rStderr, wStderr := io.Pipe()
    defer rStderr.Close()
    defer wStderr.Close()

    ioCreator := cio.NewCreator(cio.WithStreams(rStdin, wStdout, wStderr))
    process, err := task.Exec(ctx, context.String(id), pspec, ioCreator)
    if err != nil {
      return
    }
    defer process.Delete(ctx)

    statusC, err := process.Wait(ctx)
    if err != nil {
      return
    }

    sigc := commands.ForwardAllSignals(ctx, process)
    defer commands.StopCatch(sigc)

    if err := process.Start(ctx); err != nil {
      return
    }
    status := <-statusC
    code, out, err := status.Result()
    if err != nil {
      return
    }
    

    return nil
}