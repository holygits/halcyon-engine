package structs

import (
	flatbuffers "github.com/google/flatbuffers/go"

	messages "github.com/holygits/halcyon-engine/messages/system"
)

var testUserID []byte

func init() {
	testUserID = []byte("test")
}

/*
  (Un)Marshal structs to internal messaging format (flatbuffers)
*/

// Marshal serialize Func to byte stream
func (f *Func) Marshal(b *flatbuffers.Builder) []byte {
	b.Reset()

	messages.FuncStart(b)
	messages.FuncAddId(b, b.CreateByteString(f.ID))
	messages.FuncAddRuntime(b, b.CreateString(f.Runtime.String()))
	messages.FuncAddAuthor(b, b.CreateString(f.Author))
	messages.FuncAddVersion(b, b.CreateString(f.Version))
	messages.FuncAddName(b, b.CreateString(f.Name))
	messages.FuncAddCode(b, b.CreateString(f.Code))

	// Build Packages vector
	messages.FuncStartPackagesVector(b, len(f.Packages))
	for _, p := range f.Packages {
		b.PrependUOffsetT(b.CreateString(p))
	}
	packages := b.EndVector(len(f.Packages))
	messages.FuncAddPackages(b, packages)

	// Build OS Packages vector
	messages.FuncStartOsPackagesVector(b, len(f.OSPackages))
	for _, p := range f.OSPackages {
		b.PrependUOffsetT(b.CreateString(p))
	}
	os := b.EndVector(len(f.OSPackages))
	messages.FuncAddOsPackages(b, os)
	b.Finish(messages.FuncEnd(b))

	return b.FinishedBytes()
}

// Unmarshal deserialize Func from byte stream
func (f *Func) Unmarshal(buf []byte) *messages.Func {
	return messages.GetRootAsFunc(buf, 0)
}

// Marshal serialize FuncExecution to byte stream
func (f *FuncExecution) Marshal(b *flatbuffers.Builder) []byte {
	b.Reset()

	messages.FuncExecutionStart(b)
	messages.FuncExecutionAddId(b, b.CreateByteString(f.ID))
	messages.FuncExecutionAddFuncId(b, b.CreateByteString(f.FuncID))
	messages.FuncExecutionAddStdout(b, b.CreateString(f.StdOut))
	messages.FuncExecutionAddStderr(b, b.CreateString(f.StdErr))
	messages.FuncExecutionAddUserId(b, b.CreateByteString(testUserID))
	messages.FuncExecutionAddContext(b, b.CreateString(f.Ctx))
	messages.FuncExecutionAddStart(b, f.Start)
	messages.FuncExecutionAddEnd(b, f.End)
	messages.FuncExecutionAddDuration(b, f.Duration)
	b.Finish(messages.FuncExecutionEnd(b))

	return b.FinishedBytes()
}

// Unmarshal deserialize FuncExecution from byte stream
func (f *FuncExecution) Unmarshal(buf []byte) *messages.FuncExecution {
	return messages.GetRootAsFuncExecution(buf, 0)
}

// Marshal serialize Pipeline to byte stream
func (p *Pipeline) Marshal(b *flatbuffers.Builder) []byte {
	b.Reset()

	messages.PipelineStart(b)
	messages.PipelineAddId(b, b.CreateByteString(p.ID))
	messages.PipelineAddAuthor(b, b.CreateByteString(testUserID))
	messages.PipelineAddVersion(b, b.CreateString(p.Version))
	messages.PipelineAddName(b, b.CreateString(p.Name))

	// Build pipeline vector
	messages.PipelineStartStagesVector(b, len(p.Stages))
	for _, id := range p.Stages {
		b.PrependUOffsetT(b.CreateByteString(id))
	}
	funcIDs := b.EndVector(len(p.Stages))
	messages.PipelineAddStages(b, funcIDs)

	b.Finish(messages.PipelineEnd(b))

	return b.FinishedBytes()
}

// Unmarshal deserialize Pipeline from byte stream
func (p *Pipeline) Unmarshal(buf []byte) *messages.Pipeline {
	return messages.GetRootAsPipeline(buf, 0)
}

// Marshal serialize PipelineExecution to byte stream
func (p *PipelineExecution) Marshal(b *flatbuffers.Builder) []byte {
	b.Reset()

	messages.PipelineExecutionStart(b)
	messages.PipelineExecutionAddId(b, b.CreateByteString(p.ID))
	messages.PipelineExecutionAddPipelineId(b, b.CreateByteString(p.PipelineID))
	messages.PipelineExecutionAddUserId(b, b.CreateByteString(testUserID))

	// Build pipeline vector
	messages.PipelineExecutionStartStagesVector(b, len(p.Stages))
	for _, id := range p.Stages {
		b.PrependUOffsetT(b.CreateByteString(id))
	}
	funcIDs := b.EndVector(len(p.Stages))
	messages.PipelineExecutionAddStages(b, funcIDs)

	messages.PipelineExecutionAddStart(b, p.Start)
	messages.PipelineExecutionAddEnd(b, p.End)
	messages.PipelineExecutionAddDuration(b, p.Duration)

	b.Finish(messages.PipelineExecutionEnd(b))

	return b.FinishedBytes()
}

// Unmarshal deserialize PipelineExecution from byte stream
func (p *PipelineExecution) Unmarshal(buf []byte) *messages.PipelineExecution {
	return messages.GetRootAsPipelineExecution(buf, 0)
}
