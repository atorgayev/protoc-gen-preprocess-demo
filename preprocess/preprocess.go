package plugin

import (
	"fmt"

	prep "github.com/atorgayev/protoc-gen-preprocess/options"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

type preprocessor struct {
	*generator.Generator
	generator.PluginImports
}

func NewPreprocessor() *preprocessor {
	return &preprocessor{}
}

func (p *preprocessor) Name() string {
	return "preprocessor"
}

func (p *preprocessor) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *preprocessor) Generate(file *generator.FileDescriptor) {
	for _, message := range file.Messages() {
		messageName := generator.CamelCaseSlice(message.TypeName())
		for _, field := range message.Field {
			//options := getFieldOptions(field)
			/*v, err := proto.GetExtension(field.Options, preprocess.E_Field)
			if err != nil {
				p.Error(err, "In get extenstion")
			}
			opts, ok := v.(*preprocess.PreprocessorFieldOptions)
			if ok != true {
				p.Error(errors.New("shit!"))
			}*/
			test, err := proto.GetExtension(field.Options, prep.E_Field)
			p.P(fmt.Sprintf("// %v %v %v", messageName, test, err))
		}
	}
}

func getFieldOptions(field *descriptor.FieldDescriptorProto) *prep.PreprocessFieldOptions {
	if field.Options == nil {
		return nil
	}
	v, err := proto.GetExtension(field.Options, prep.E_Field)
	if err != nil {
		panic(err)
	}
	opts, ok := v.(*prep.PreprocessFieldOptions)
	if !ok {
		return nil
	}
	return opts
}

func (p *preprocessor) GenerateImports(file *generator.FileDescriptor) {}

func init() {
	generator.RegisterPlugin(NewPreprocessor())
}
