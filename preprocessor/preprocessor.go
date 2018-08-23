package preprocessor

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

type preprocessor struct {
	*generator.Generator
	generator.PluginImports
	atleastOne bool
	localName  string
	overwrite  bool
}

func NewPreprocessor() *preprocessor {
	return &preprocessor{}
}

func (p *preprocessor) Name() string {
	return "preprocessor"
}

func (p *preprocessor) Overwrite() {
	p.overwrite = true
}

func (p *preprocessor) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *preprocessor) Generate(file *generator.FileDescriptor) {
	for _, message := range file.Messages() {
		ccTypeName := generator.CamelCaseSlice(message.TypeName())
		p.P(`func (this *`, ccTypeName, `) Preprocess() string {`)
		p.In()
		p.P(`if this == nil {`)
		p.In()
		p.P(`return "nil"`)
		p.Out()
		p.P(`}`)
		/*for _, field := range message.Field {
			p.P(fmt.Sprintf(field.String()))
		}*/
		p.P(`}`)
	}
}

func (p *preprocessor) GenerateOld(file *generator.FileDescriptor) {
	proto3 := gogoproto.IsProto3(file.FileDescriptorProto)
	p.PluginImports = generator.NewPluginImports(p.Generator)
	p.atleastOne = false

	p.localName = generator.FileName(file)

	fmtPkg := p.NewImport("fmt")
	stringsPkg := p.NewImport("strings")
	protoPkg := p.NewImport("github.com/gogo/protobuf/proto")
	if !gogoproto.ImportsGoGoProto(file.FileDescriptorProto) {
		protoPkg = p.NewImport("github.com/golang/protobuf/proto")
	}
	sortPkg := p.NewImport("sort")
	strconvPkg := p.NewImport("strconv")
	reflectPkg := p.NewImport("reflect")
	sortKeysPkg := p.NewImport("github.com/gogo/protobuf/sortkeys")

	extensionTopreprocessorUsed := false
	for _, message := range file.Messages() {
		if !p.overwrite && !gogoproto.HasGoString(file.FileDescriptorProto, message.DescriptorProto) {
			continue
		}
		if message.DescriptorProto.GetOptions().GetMapEntry() {
			continue
		}
		p.atleastOne = true
		packageName := file.GoPackageName()

		ccTypeName := generator.CamelCaseSlice(message.TypeName())
		p.P(`func (this *`, ccTypeName, `) Preprocessor() string {`)
		p.In()
		p.P(`if this == nil {`)
		p.In()
		p.P(`return "nil"`)
		p.Out()
		p.P(`}`)

		p.P(`s := make([]string, 0, `, strconv.Itoa(len(message.Field)+4), `)`)
		p.P(`s = append(s, "&`, packageName, ".", ccTypeName, `{")`)

		oneofs := make(map[string]struct{})
		for _, field := range message.Field {
			nullable := gogoproto.IsNullable(field)
			repeated := field.IsRepeated()
			fieldname := p.GetFieldName(message, field)
			oneof := field.OneofIndex != nil
			if oneof {
				if _, ok := oneofs[fieldname]; ok {
					continue
				} else {
					oneofs[fieldname] = struct{}{}
				}
				p.P(`if this.`, fieldname, ` != nil {`)
				p.In()
				p.P(`s = append(s, "`, fieldname, `: " + `, fmtPkg.Use(), `.Sprintf("%#v", this.`, fieldname, `) + ",\n")`)
				p.Out()
				p.P(`}`)
			} else if p.IsMap(field) {
				m := p.GoMapType(nil, field)
				mapgoTyp, keyField, keyAliasField := m.GoType, m.KeyField, m.KeyAliasField
				keysName := `keysFor` + fieldname
				keygoTyp, _ := p.GoType(nil, keyField)
				keygoTyp = strings.Replace(keygoTyp, "*", "", 1)
				keygoAliasTyp, _ := p.GoType(nil, keyAliasField)
				keygoAliasTyp = strings.Replace(keygoAliasTyp, "*", "", 1)
				keyCapTyp := generator.CamelCase(keygoTyp)
				p.P(keysName, ` := make([]`, keygoTyp, `, 0, len(this.`, fieldname, `))`)
				p.P(`for k, _ := range this.`, fieldname, ` {`)
				p.In()
				if keygoAliasTyp == keygoTyp {
					p.P(keysName, ` = append(`, keysName, `, k)`)
				} else {
					p.P(keysName, ` = append(`, keysName, `, `, keygoTyp, `(k))`)
				}
				p.Out()
				p.P(`}`)
				p.P(sortKeysPkg.Use(), `.`, keyCapTyp, `s(`, keysName, `)`)
				mapName := `mapStringFor` + fieldname
				p.P(mapName, ` := "`, mapgoTyp, `{"`)
				p.P(`for _, k := range `, keysName, ` {`)
				p.In()
				if keygoAliasTyp == keygoTyp {
					p.P(mapName, ` += fmt.Sprintf("%#v: %#v,", k, this.`, fieldname, `[k])`)
				} else {
					p.P(mapName, ` += fmt.Sprintf("%#v: %#v,", k, this.`, fieldname, `[`, keygoAliasTyp, `(k)])`)
				}
				p.Out()
				p.P(`}`)
				p.P(mapName, ` += "}"`)
				p.P(`if this.`, fieldname, ` != nil {`)
				p.In()
				p.P(`s = append(s, "`, fieldname, `: " + `, mapName, `+ ",\n")`)
				p.Out()
				p.P(`}`)
			} else if (field.IsMessage() && !gogoproto.IsCustomType(field) && !gogoproto.IsStdTime(field) && !gogoproto.IsStdDuration(field)) || p.IsGroup(field) {
				if nullable || repeated {
					p.P(`if this.`, fieldname, ` != nil {`)
					p.In()
				}
				if nullable {
					p.P(`s = append(s, "`, fieldname, `: " + `, fmtPkg.Use(), `.Sprintf("%#v", this.`, fieldname, `) + ",\n")`)
				} else if repeated {
					if nullable {
						p.P(`s = append(s, "`, fieldname, `: " + `, fmtPkg.Use(), `.Sprintf("%#v", this.`, fieldname, `) + ",\n")`)
					} else {
						goTyp, _ := p.GoType(message, field)
						goTyp = strings.Replace(goTyp, "[]", "", 1)
						p.P("vs := make([]*", goTyp, ", len(this.", fieldname, "))")
						p.P("for i := range vs {")
						p.In()
						p.P("vs[i] = &this.", fieldname, "[i]")
						p.Out()
						p.P("}")
						p.P(`s = append(s, "`, fieldname, `: " + `, fmtPkg.Use(), `.Sprintf("%#v", vs) + ",\n")`)
					}
				} else {
					p.P(`s = append(s, "`, fieldname, `: " + `, stringsPkg.Use(), `.Replace(this.`, fieldname, `.preprocessor()`, ",`&`,``,1)", ` + ",\n")`)
				}
				if nullable || repeated {
					p.Out()
					p.P(`}`)
				}
			} else {
				if !proto3 && (nullable || repeated) {
					p.P(`if this.`, fieldname, ` != nil {`)
					p.In()
				}
				if field.IsEnum() {
					if nullable && !repeated && !proto3 {
						goTyp, _ := p.GoType(message, field)
						p.P(`s = append(s, "`, fieldname, `: " + valueTopreprocessor`, p.localName, `(this.`, fieldname, `,"`, generator.GoTypeToName(goTyp), `"`, `) + ",\n")`)
					} else {
						p.P(`s = append(s, "`, fieldname, `: " + `, fmtPkg.Use(), `.Sprintf("%#v", this.`, fieldname, `) + ",\n")`)
					}
				} else {
					if nullable && !repeated && !proto3 {
						goTyp, _ := p.GoType(message, field)
						p.P(`s = append(s, "`, fieldname, `: " + valueTopreprocessor`, p.localName, `(this.`, fieldname, `,"`, generator.GoTypeToName(goTyp), `"`, `) + ",\n")`)
					} else {
						p.P(`s = append(s, "`, fieldname, `: " + `, fmtPkg.Use(), `.Sprintf("%#v", this.`, fieldname, `) + ",\n")`)
					}
				}
				if !proto3 && (nullable || repeated) {
					p.Out()
					p.P(`}`)
				}
			}
		}
		if message.DescriptorProto.HasExtension() {
			if gogoproto.HasExtensionsMap(file.FileDescriptorProto, message.DescriptorProto) {
				p.P(`s = append(s, "XXX_InternalExtensions: " + extensionTopreprocessor`, p.localName, `(this) + ",\n")`)
				extensionTopreprocessorUsed = true
			} else {
				p.P(`if this.XXX_extensions != nil {`)
				p.In()
				p.P(`s = append(s, "XXX_extensions: " + `, fmtPkg.Use(), `.Sprintf("%#v", this.XXX_extensions) + ",\n")`)
				p.Out()
				p.P(`}`)
			}
		}
		if gogoproto.HasUnrecognized(file.FileDescriptorProto, message.DescriptorProto) {
			p.P(`if this.XXX_unrecognized != nil {`)
			p.In()
			p.P(`s = append(s, "XXX_unrecognized:" + `, fmtPkg.Use(), `.Sprintf("%#v", this.XXX_unrecognized) + ",\n")`)
			p.Out()
			p.P(`}`)
		}

		p.P(`s = append(s, "}")`)
		p.P(`return `, stringsPkg.Use(), `.Join(s, "")`)
		p.Out()
		p.P(`}`)

		//Generate preprocessor methods for oneof fields
		for _, field := range message.Field {
			oneof := field.OneofIndex != nil
			if !oneof {
				continue
			}
			ccTypeName := p.OneOfTypeName(message, field)
			p.P(`func (this *`, ccTypeName, `) preprocessor() string {`)
			p.In()
			p.P(`if this == nil {`)
			p.In()
			p.P(`return "nil"`)
			p.Out()
			p.P(`}`)
			fieldname := p.GetOneOfFieldName(message, field)
			outStr := strings.Join([]string{
				"s := ",
				stringsPkg.Use(), ".Join([]string{`&", packageName, ".", ccTypeName, "{` + \n",
				"`", fieldname, ":` + ", fmtPkg.Use(), `.Sprintf("%#v", this.`, fieldname, `)`,
				" + `}`",
				`}`,
				`,", "`,
				`)`}, "")
			p.P(outStr)
			p.P(`return s`)
			p.Out()
			p.P(`}`)
		}
	}

	if !p.atleastOne {
		return
	}

	p.P(`func valueTopreprocessor`, p.localName, `(v interface{}, typ string) string {`)
	p.In()
	p.P(`rv := `, reflectPkg.Use(), `.ValueOf(v)`)
	p.P(`if rv.IsNil() {`)
	p.In()
	p.P(`return "nil"`)
	p.Out()
	p.P(`}`)
	p.P(`pv := `, reflectPkg.Use(), `.Indirect(rv).Interface()`)
	p.P(`return `, fmtPkg.Use(), `.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)`)
	p.Out()
	p.P(`}`)

	if extensionTopreprocessorUsed {
		if !gogoproto.ImportsGoGoProto(file.FileDescriptorProto) {
			fmt.Fprintf(os.Stderr, "The preprocessor plugin for messages with extensions requires importing gogoprotobuf. Please see file %s", file.GetName())
			os.Exit(1)
		}
		p.P(`func extensionTopreprocessor`, p.localName, `(m `, protoPkg.Use(), `.Message) string {`)
		p.In()
		p.P(`e := `, protoPkg.Use(), `.GetUnsafeExtensionsMap(m)`)
		p.P(`if e == nil { return "nil" }`)
		p.P(`s := "proto.NewUnsafeXXX_InternalExtensions(map[int32]proto.Extension{"`)
		p.P(`keys := make([]int, 0, len(e))`)
		p.P(`for k := range e {`)
		p.In()
		p.P(`keys = append(keys, int(k))`)
		p.Out()
		p.P(`}`)
		p.P(sortPkg.Use(), `.Ints(keys)`)
		p.P(`ss := []string{}`)
		p.P(`for _, k := range keys {`)
		p.In()
		p.P(`ss = append(ss, `, strconvPkg.Use(), `.Itoa(k) + ": " + e[int32(k)].preprocessor())`)
		p.Out()
		p.P(`}`)
		p.P(`s+=`, stringsPkg.Use(), `.Join(ss, ",") + "})"`)
		p.P(`return s`)
		p.Out()
		p.P(`}`)
	}
}

func (p *preprocessor) GenerateImports(file *generator.FileDescriptor) {}

func init() {
	generator.RegisterPlugin(NewPreprocessor())
}
