package structs

import (
	flatbuffers "github.com/google/flatbuffers/go"

	messages "github.com/holygits/halcyon-engine/messages/structs"
)

const (
	TestUser = "test"
)

/*
  (Un)Marshal structs to internal messaging format (flatbuffers)
*/

// Marshal serialize Func to byte stream
func (f *Func) Marshal() ([]byte, error) {
	// TODO: get / re-use buffer from some pool
	b := flatbuffers.NewBuilder()
	messages.FuncStart(b)
	messages.FuncAddId(f.ID)
	messages.FuncAddRuntime(f.Runtime)
	messages.FuncAddAuthor([]byte(TestUser))
	messages.FuncAddVersion(f.Version)
	messages.FuncAddName(f.Name)
	messages.FuncAddSource(f.Source)
	messages.FuncAddPackages(f.Packages)
	messages.FuncAddOsPackages(f.OSPackages)
	messages.FuncEnd(b)
	b.Finish()
	return b.FinishedBytes(), nil
}

// Unmarshal deserialize Func from byte stream
func (f *Func) Unmarshal(buf []byte) error {
	return messages.GetRootAsFunc(buf, 0), nil
}

// Marshal serialize FuncExecution to byte stream
func (f *FuncExecution) Marshal() ([]byte, error) {
	// TODO: get / re-use buffer from some pool
	b := flatbuffers.NewBuilder()
	messages.FuncExecutionStart(b)
	messages.FuncExecutionAddId(f.ID)
	messages.FuncExecutionAddFuncId(f.FuncID)
	messages.FuncExecutionAddStdOut(f.StdOut)
	messages.FuncExecutionAddStdErr(f.StdErr)
	messages.FuncExecutionAddUser([]byte(TestUser))
	messages.FuncExecutionAddContext(f.Ctx)
	messages.FuncExecutionAddStart(f.Start)
	messages.FuncExecutionAddEnd(f.End)
	messages.FuncExecutionAddDuration(f.Duration)
	messages.FuncExecutionEnd(b)
	b.Finish()
	return b.FinishedBytes(), nil
}

// Unmarshal deserialize FuncExecution from byte stream
func (f *FuncExecution) Unmarshal(buf []byte) error {
	return messages.GetRootAsFuncExecution(buf, 0), nil
}

// Marshal serialize Pipeline to byte stream
func (p *Pipeline) Marshal() ([]byte, error) {
	b := flatbuffers.NewBuilder()
	messages.PipelineStart(b)
	messages.PipelineAddId(p.ID)
	messages.PipelineAddAuthor([]byte(TestUser))
	messages.PipelineAddVersion(p.Version)
	messages.PipelineAddName(p.Name)
	messages.PipelineAddSource(p.Source)
	messages.PipelineEnd(b)
	b.Finish()
	return b.FinishedBytes(), nil
}

// Unmarshal deserialize Pipeline from byte stream
func (p *Pipeline) Unmarshal(buf []byte) error {
	return messages.GetRootAsPipeline(buf, 0), nil
}

// Marshal serialize PipelineExecution to byte stream
func (p *PipelineExecution) Marshal() ([]byte, error) {
	b := flatbuffers.NewBuilder()
	messages.PipelineExecutionStart(b)
	messages.PipelineExecutionAddId(p.ID)
	messages.PipelineExecutionAddPipelineId(p.PipelineID)
	messages.PipelineExecutionAddUser([]byte(TestUser))
	messages.PipelineExecutionAddStages(p.Stages)
	messages.PipelineExecutionAddStart(p.Start)
	messages.PipelineExecutionAddEnd(p.End)
	messages.PipelineExecutionAddDuration(p.Duration)
	messages.PipelineExecutionEnd(b)
	b.Finish()
	return b.FinishedBytes(), nil
}

// Unmarshal deserialize PipelineExecution from byte stream
func (p *PipelineExecution) Unmarshal(buf []byte) error {
	return messages.GetRootAsPipelineExecution(buf, 0), nil
}
