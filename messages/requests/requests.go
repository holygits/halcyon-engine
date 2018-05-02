// Package requests provides definitions for external API request messages
package requests

/*
  (Un)Marshal request structs to external messaging format (JSON)
*/

//NewFuncRequest describes a client API request to create a new function
type NewFunc struct {
	Name, Runtime, Source, Version string
	Packages, OSPackages           []string
}

//NewPipelineRequest describes a client API request to create a new pipeline
type NewPipeline struct {
	Source  *Pipe
	Name    string // Human-friendly name for the function
	Version string // Version number like docker image tag e.g. `:latest`
}
