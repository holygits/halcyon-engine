package builder

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/holygits/halcyon-engine/structs"
)

type RuntimeBuilder struct {
	buildPath string
}

// BuildPackagecreates a dependency package for a function
func (r *RuntimeBuilder) BuildPackage(f *structs.Func, s *buildStrategy) error {
	if err := os.Mkdir(r.buildPath, 0644); err != nil {
		return err
	}
	// TODO: run container with target dir mount, command pulls packages into some target dir
}

type buildStrategy interface {
	Command() []string
}

func init() {
	// TODO: Pull latest build containers for each runtime
}

func build(f *structs.Func) error {
	switch f.Runtime {
	case structs.Runtime.Go:
	case structs.Runtime.Javascript:
	case structs.Runtime.Python:
		// TODO: submit packages to the runtime builder container

	case structs.Runtime.R:
	default:
		return errors.New("Invalid Runtime: " + strconv.Itoa(f.Runtime))
	}
}

/*
Build container process
1. Request received, deserialized as structs.Func
2. Invoke an image build
   - Installs runtime dependencies in a build container, outputs as dir
   - Generate Dockerfile with `ADD dir`
   - Build docker image
   - Export image as tar blob
   - Save tar blob to datastore
*/

/* Pandas / Numpy builds are unacceptably slow, however the initial build will be amortised when used again
   two solutions:
     - Eagerly build common slow packages e.g. {numpy, pandas}? not a generic solution, focuses on edge cases
     - Switch to a debian distro for such packages
     - Host a private pypi with pre-compiled wheels (hosting costs?)
     - Upload packages to pypi
*/
