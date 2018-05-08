package structs

import (
	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/holygits/halcyon-engine/messages/requests"
)

// Marshal serialize FuncExecution to byte stream
func (f *NewFunc) Marshal(b *flatbuffers.Builder) []byte {
	b.Reset()

	requests.NewFuncAddRuntime(b, b.CreateString(f.Runtime))
	requests.NewFuncAddAuthor(b, b.CreateByteString(testUserID))
	requests.NewFuncAddVersion(b, b.CreateString(f.Version))
	requests.NewFuncAddName(b, b.CreateString(f.Name))
	requests.NewFuncAddSource(b, b.CreateString(f.Source))

	// Build Packages vector
	requests.NewFuncStartPackagesVector(b, len(f.Packages))
	for _, p := range f.Packages {
		b.PrependUOffsetT(b.CreateString(p))
	}
	packages := b.EndVector(len(f.Packages))
	requests.NewFuncAddPackages(b, packages)

	// Build OS Packages vector
	requests.NewFuncStartOsPackagesVector(b, len(f.OSPackages))
	for _, p := range f.OSPackages {
		b.PrependUOffsetT(b.CreateString(p))
	}
	os := b.EndVector(len(f.OSPackages))
	requests.NewFuncAddOsPackages(b, os)

	b.Finish(requests.NewFuncEnd(b))

	return b.FinishedBytes()
}

// Unmarshal deserialize NewFunc from byte stream
func (f *NewFunc) Unmarshal(buf []byte) *requests.NewFunc {
	return requests.GetRootAsNewFunc(buf, 0)
}

// Marshal serialize Pipeline to byte stream
func (p *NewPipeline) Marshal(b *flatbuffers.Builder) []byte {
	b.Reset()

	requests.NewPipelineStart(b)
	requests.NewPipelineAddAuthor(b, b.CreateByteString(testUserID))
	requests.NewPipelineAddVersion(b, b.CreateString(p.Version))
	requests.NewPipelineAddName(b, b.CreateString(p.Name))

	// Build pipeline vector
	requests.NewPipelineStartPipelineVector(b, len(p.pipeline))
	for _, f := range p.pipeline {
		b.PrependUOffsetT(b.CreateByteString(f))
	}
	funcIDs := b.EndVector(len(p.pipeline))
	requests.NewPipelineAddPipeline(b, funcIDs)

	b.Finish(requests.NewPipelineEnd(b))

	return b.FinishedBytes()
}

// Unmarshal deserialize NewPipeline from byte stream
func (p *NewPipeline) Unmarshal(buf []byte) *requests.NewPipeline {
	return requests.GetRootAsNewPipeline(buf, 0)
}
