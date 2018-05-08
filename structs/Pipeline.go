package structs

// Pipeline is an ordered "pipeline" of functions
type Pipeline struct {
	ID      []byte
	Stages  [][]byte // ordered list of func IDs
	Author  string
	Name    string
	Version string
}

// PipelineExecution represents a single invocation of a pipeline
type PipelineExecution struct {
	ID         []byte
	PipelineID []byte
	Ctx        string   // JSON encoded string
	Stages     [][]byte // ordered list of func execution IDs
	Start      uint32
	End        uint32
	Duration   uint32
}
