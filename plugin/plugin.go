package plugin

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gogo/protobuf/vanity"
)

type plugin struct {
	*generator.Generator
	generator.PluginImports

	regexPkg        generator.Single
	fmtPkg          generator.Single
	protoPkg        generator.Single
	preprocessorPkg generator.Single
	useGogoImport   bool
}

func (p *plugin) Name() string {
	return "preprocessor"
}

func (p *plugin) Init(g *generator.Generator) {
	p.Generator = g
}

func NewPlugin(useGogoImport bool) generator.Plugin {
	return &plugin{useGogoImport: useGogoImport}
}

func (p *plugin) Generate(file *generator.FileDescriptor) {
	if !p.useGogoImport {
		vanity.TurnOffGogoImport(file.FileDescriptorProto)
	}
	p.PluginImports = generator.NewPluginImports(p.Generator)
	p.regexPkg = p.NewImport("regexp")
	p.fmtPkg = p.NewImport("fmt")
	p.preprocessorPkg = p.NewImport("github.com/atorgayev/protoc-gen-preprocess")

	for _, msg := range file.Messages() {
		if msg.DescriptorProto.GetOptions().GetMapEntry() {
			continue
		}
		p.generateProto3Message(file, msg)

	}
}
