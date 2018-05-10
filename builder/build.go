package builder

import (
	"errors"
	"os"
  "fmt"
	"path"
	"strconv"
  "strings"
  "text/template"
  "sync"

	//"github.com/cespare/xxhash"

	"github.com/holygits/halcyon-engine/containerd"
	"github.com/holygits/halcyon-engine/structs"
)

type runtimeBuilder struct {
	buildPath string
	c       *containerd.Containerd
}

type Package struct {
  Path string
}

func init() {
  
  // TODO: needs to be moved to correct location
  templateDir := os.Getenv("HALCYON_TEMPLATE_DIR")
  if len(templateDir) == 0 {
    // Running in development mode
    templateDir = "./template"
  }
  templateDir = strings.TrimSuffix(templateDir, "/")

  // Load Templates
  var err error
  images, err = template.ParseGlob(path.Join(templateDir, "docker", "*.Dockerfile"))
  if err != nil {
    log.Fatalf("Failed loading docker templates in: %s", err)
  }
  funcfiles, err = template.ParseGlob(path.Join(templateDir, "func", "*"))
  if err != nil {
    log.Fatalf("Failed loading func templates in: %s", err)
  }

}

func newRuntimeBuilder(buildPath string) (*runtimeBuilder, error) {

	if err := os.Mkdir(buildPath, 0644); err != nil {
		// TODO: handle OSExists error?
		return nil, err
	}

	c, err := containerd.NewContainerd()
	if err != nil {
		return nil, err
	}

  // Pull runtime resolver images
  var wg sync.WaitGroup
  wg.Add(len(structs.Runtimes))
  for _, r := range structs.Runtimes {
    id := fmt.Sprintf(containerd.ResolverRepo, r, ":latest")
    go func() {
      defer wg.Done()

      img, err := c.Pull(id)
      if err != nil {
        return nil, err
      }

      // TODO: start resolver container
      redis, err := c.NewContainer(c.Context, img)
      if err != nil {
        return nil, err
      }
      task, err := redis.NewTask(context, )
      if err != nil {
        return nil, err
      }
    }
  }
  wg.Wait()


	return &runtimeBuilder{
		buildPath: buildPath,
		c:       c,
	}, nil
}

// BuildPackage creates a dependency package for a function e.g. venv, npm package
func (r *runtimeBuilder) BuildPackage(f *structs.Func, s *buildStrategy) error {
	// No build needed
	// TODO: Move upstream to API, pre-request
	// TODO: Also check packages hash against existing build, skip if equal
	if len(f.Packages) == 0 {
		return nil
	}

	depsPath := path.Join(r.buildPath, string(f.ID)+" "+f.Version)
	if err := os.Mkdir(depsPath, 0644); err != nil {
		return err
	}

	// TODO: Run container with target dir mount, command pulls packages into some target dir
	r.cli.ContainerService().Create(ctx, container)

	return nil
}

// Publish builds func image and pushes to registry
func Publish(id, pack string) (err error) {

  // Build package as Docker Image & Push to repository.
  if err = docker.ImageBuild(id, pack); err != nil {
    log.Errorf("Failed to build image: %s", err)
    return
  }

  if err = docker.ImagePush(id); err != nil {
    return
  }

  log.Infof("Published func: %s", id)
  retur

type buildStrategy interface {
	Command() []string
}

func build(f *structs.Func) error {
	/* cases:
	    1. No existing image (initial case)
	    2. Image exists
	       - A: packages are the same
	       - B: packages have changed
	       - C: OSPackages are the same
	       - D: OSPackages have changed

	   Every build should:
	   1. Resolve runtime dependencies
	      - skip if "packages" is empty or "packages" is unchanged from a previous build
	   2. Template a function file with function source code
     3. Generate a dockerfile
	   3. Build the dockerfile and export as tar
	*/

	switch f.Runtime {
	case structs.RuntimeJavascript:
    // e.g. yarn add <packages> --modules-folder=<<func-id-version>
	case structs.RuntimePython:
		// TODO: submit packages to the runtime builder container
    // e.g. pip install --install-option="--prefix=<func-id-version>" -f
	case structs.RuntimeR:
    // e.g. install.packages(c('p1', 'p2')lib='./<func-id-version>')
	default:
		return errors.New("Invalid Runtime: " + strconv.Itoa(int(f.Runtime)))
	}

	return nil
}

/*
Build container process
1. Request received, deserialized as structs.Func
2. Invoke an image build
   - Installs runtime dependencies in a build container, outputs as a dir
   - Generate Dockerfile with `ADD dir`
   - Build docker image
   - Export image as tar blob
   - Save tar blob to datastore
*/
