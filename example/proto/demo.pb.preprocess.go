// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: demo.proto

package demo

import strings "strings"
import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/atorgayev/protoc-gen-preprocess/options"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (m *Demo) Preprocess() error {
	m.S = strings.TrimSpace(m.S)
	return nil
}

func (m *DemoReq) Preprocess() error {
	return nil
}

func (m *DemoRes) Preprocess() error {
	return nil
}
