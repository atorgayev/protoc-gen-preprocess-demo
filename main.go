package main

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	preprocessor_plugin "github.com/atorgayev/protoc-gen-preprocess/plugin"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

func main() {
	g := generator.New()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err)
	}

	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err)
	}

	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}

	useGogoImport := false
	// Match parsing algorithm from Generator.CommandLineParameters
	for _, parameter := range strings.Split(g.Request.GetParameter(), ",") {
		kvp := strings.SplitN(parameter, "=", 2)
		// We only care about key-value pairs where the key is "gogoimport"
		if len(kvp) != 2 || kvp[0] != "gogoimport" {
			continue
		}
		useGogoImport, err = strconv.ParseBool(kvp[1])
		if err != nil {
			g.Error(err, "parsing gogoimport option")
		}
	}

	g.CommandLineParameters(g.Request.GetParameter())

	g.WrapTypes()
	g.SetPackageNames()
	g.BuildTypeNameMap()
	g.GeneratePlugin(preprocessor_plugin.NewPlugin(useGogoImport))

	for i := 0; i < len(g.Response.File); i++ {
		g.Response.File[i].Name = proto.String(strings.Replace(*g.Response.File[i].Name, ".pb.go", ".preprocess.pb.go", -1))
	}

	// Send back the results.
	data, err = proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}

}
