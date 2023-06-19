// Code generated by protoc-gen-go.
// source: proto_session.proto
// DO NOT EDIT!

package ProtoExample

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type InitSessionReq struct {
	SessionID int64 `protobuf:"varint,1,opt,name=sessionID" json:"sessionID,omitempty"`
}

func (m *InitSessionReq) Reset()                    { *m = InitSessionReq{} }
func (m *InitSessionReq) String() string            { return proto.CompactTextString(m) }
func (*InitSessionReq) ProtoMessage()               {}
func (*InitSessionReq) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *InitSessionReq) GetSessionID() int64 {
	if m != nil {
		return m.SessionID
	}
	return 0
}

type InitSessionAck struct {
	Error     ErrorCode `protobuf:"varint,1,opt,name=error,enum=ProtoExample.ErrorCode" json:"error,omitempty"`
	SessionID int64     `protobuf:"varint,2,opt,name=sessionID" json:"sessionID,omitempty"`
}

func (m *InitSessionAck) Reset()                    { *m = InitSessionAck{} }
func (m *InitSessionAck) String() string            { return proto.CompactTextString(m) }
func (*InitSessionAck) ProtoMessage()               {}
func (*InitSessionAck) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func (m *InitSessionAck) GetError() ErrorCode {
	if m != nil {
		return m.Error
	}
	return ErrorCode_OK
}

func (m *InitSessionAck) GetSessionID() int64 {
	if m != nil {
		return m.SessionID
	}
	return 0
}

func init() {
	proto.RegisterType((*InitSessionReq)(nil), "ProtoExample.InitSessionReq")
	proto.RegisterType((*InitSessionAck)(nil), "ProtoExample.InitSessionAck")
}

func init() { proto.RegisterFile("proto_session.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 139 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2e, 0x28, 0xca, 0x2f,
	0xc9, 0x8f, 0x2f, 0x4e, 0x2d, 0x2e, 0xce, 0xcc, 0xcf, 0xd3, 0x03, 0xf3, 0x84, 0x78, 0x02, 0x40,
	0x94, 0x6b, 0x45, 0x62, 0x6e, 0x41, 0x4e, 0xaa, 0x94, 0x20, 0x44, 0x49, 0x6a, 0x51, 0x51, 0x7e,
	0x11, 0x44, 0x81, 0x92, 0x1e, 0x17, 0x9f, 0x67, 0x5e, 0x66, 0x49, 0x30, 0x44, 0x57, 0x50, 0x6a,
	0xa1, 0x90, 0x0c, 0x17, 0x27, 0xd4, 0x0c, 0x4f, 0x17, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xe6, 0x20,
	0x84, 0x80, 0x52, 0x2c, 0x8a, 0x7a, 0xc7, 0xe4, 0x6c, 0x21, 0x5d, 0x2e, 0x56, 0xb0, 0x81, 0x60,
	0xb5, 0x7c, 0x46, 0xe2, 0x7a, 0xc8, 0x56, 0xea, 0xb9, 0x82, 0xa4, 0x9c, 0xf3, 0x53, 0x52, 0x83,
	0x20, 0xaa, 0x50, 0x8d, 0x67, 0x42, 0x33, 0x3e, 0x89, 0x0d, 0xec, 0x2a, 0x63, 0x40, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x87, 0x3e, 0x77, 0xa4, 0xcd, 0x00, 0x00, 0x00,
}