package docker

const (
	BuilderRepo  = "halcyo/builder"
	EngineRepo   = "halcyo/engine"
	ResolverRepo = "halcyo/%s-resolver"
	RuntimeRepo  = "halcyo/%s-runtime"
	WorkerRepo   = "halcyo/worker"

	// RunForeverCmd a command which keeps a container running indefinitley
	RunForeverCmd = "tail -f /dev/null"
)
