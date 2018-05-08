// Package broker provides interface and connectivity to message broker services
package broker

type Broker interface {
	// QFunc adds a a new func execution request to the broker
	QFunc() error
	// QPipeline adds a a new pipeline execution request to the broker
	QPipeline() error
	// QFuncBuild adds a func build job to the broker
	QFuncBuild() error
	// QFuncDel adds a func delete job to the broker
	QFuncDel() error
	QFuncResp() error
	QPipelineResp() error
	// Poll starts polling the broker for messages
	Poll()
}
