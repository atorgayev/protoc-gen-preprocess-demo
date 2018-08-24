package plugin

import (
	"bytes"
	"text/template"

	prep "github.com/atorgayev/protoc-gen-preprocess/options"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

type rule map[string]map[string][]string

type preprocessor struct {
	*generator.Generator
	generator.PluginImports
	rules       rule
	messageName string
	fieldName   string
}

func NewPreprocessor() *preprocessor {
	p := &preprocessor{}
	return p
}

func (p *preprocessor) Name() string {
	return "preprocessor"
}

func (p *preprocessor) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *preprocessor) Generate(file *generator.FileDescriptor) {
	for _, message := range file.Messages() {
		p.messageName = generator.CamelCaseSlice(message.TypeName())
		for _, field := range message.Field {
			p.fieldName = field.GetName()
			options := getFieldOptions(field)
			if options == nil {
				continue
			}
			if options.GetText().GetTrimSpace() {
				p.TextTrimSpace()
			}
		}
	}
	generateFromTemplate(p)
}

func (p *preprocessor) TextTrimSpace() {
	for p.rules[p.messageName][p.fieldName] == nil {
		switch {
		case p.rules == nil:
			p.rules = make(map[string]map[string][]string)
		case p.rules[p.messageName] == nil:
			p.rules[p.messageName] = make(map[string][]string)
		case p.rules[p.messageName][p.fieldName] == nil:
			p.rules[p.messageName][p.fieldName] = make([]string, 0)
		}
	}
	fieldRules := p.rules[p.messageName][p.fieldName]
	p.rules[p.messageName][p.fieldName] = append(fieldRules, "trimSpace")
}

func (p *preprocessor) GenerateImports(file *generator.FileDescriptor) {}

func init() {
	generator.RegisterPlugin(NewPreprocessor())
}

func getFieldOptions(field *descriptor.FieldDescriptorProto) *prep.PreprocessFieldOptions {
	if field.Options == nil {
		return nil
	}
	v, err := proto.GetExtension(field.Options, prep.E_Field)
	if err != nil {
		return nil
	}
	opts, ok := v.(*prep.PreprocessFieldOptions)
	if !ok {
		return nil
	}
	return opts
}

func generateFromTemplate(p *preprocessor) {
	const function = `
func (m *{{.Name}}) Preprocess() {
	{{ with .Fields}}{{ range .}}
		*m.{{.}} = strings.TrimSpace(m.{{.}})
	{{ end }}{{ end }}
}	
`
	var tpl bytes.Buffer
	t := template.New("rules")
	t, err := t.Parse(function)
	if err != nil {
	}

	for mn, m := range p.rules {
		fields := make([]string, 0)
		for fn := range m {
			fields = append(fields, fn)
		}

		data := struct {
			Name   string
			Fields []string
		}{
			Name:   mn,
			Fields: fields,
		}

		t.Execute(&tpl, data)
	}

	p.P(tpl.String())
	p.P()
}
