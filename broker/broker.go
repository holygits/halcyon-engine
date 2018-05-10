// Package broker provides interface and connectivity to message broker services
package broker

const (
  // default broker queue names
  PipelineExecQ = "pExec"
  PipelineRespQ = "pResp"
  FuncBuildQ = "fBuild"
  FuncBuildRespQ = "fBuildResp"
)

// Message contains a brokered message payload and associated handling errors
type Mesage struct {
  Body []byte
  Error error
}

// Broker provides a simple PubSub interface for halcyon message queueing
type Broker interface {
  // Publish pushes a new message to a queue
  Publish(q string, msg *Message) error

  // Subscribe listens for new messages on a queue, provides updates over a channel
  Subscribe(q string) (<-chan *Message, error)
}
