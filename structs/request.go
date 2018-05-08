package structs

//NewFunc describes a client API request to create a new function
type NewFunc struct {
	Name, Runtime, Source, Version string
	Packages, OSPackages           []string
}

//NewPipeline describes a client API request to create a new pipeline
type NewPipeline struct {
	Name     string // Human-friendly name for the function
	Version  string // Version number like docker image tag e.g. `:latest`
	Author   string
	pipeline [][]byte // Ordered list of func IDs
}
