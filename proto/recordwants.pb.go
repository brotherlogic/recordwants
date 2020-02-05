// Code generated by protoc-gen-go. DO NOT EDIT.
// source: recordwants.proto

package recordwants

import (
	fmt "fmt"
	godiscogs "github.com/brotherlogic/godiscogs"
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MasterWant_Level int32

const (
	MasterWant_UNKNOWN MasterWant_Level = 0
	MasterWant_ANYTIME MasterWant_Level = 1
	MasterWant_LIST    MasterWant_Level = 2
	MasterWant_ALWAYS  MasterWant_Level = 3
)

var MasterWant_Level_name = map[int32]string{
	0: "UNKNOWN",
	1: "ANYTIME",
	2: "LIST",
	3: "ALWAYS",
}

var MasterWant_Level_value = map[string]int32{
	"UNKNOWN": 0,
	"ANYTIME": 1,
	"LIST":    2,
	"ALWAYS":  3,
}

func (x MasterWant_Level) String() string {
	return proto.EnumName(MasterWant_Level_name, int32(x))
}

func (MasterWant_Level) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_033e63fc596348d2, []int{2, 0}
}

type Config struct {
	Wants                []*MasterWant          `protobuf:"bytes,1,rep,name=wants,proto3" json:"wants,omitempty"`
	Budget               int32                  `protobuf:"varint,2,opt,name=budget,proto3" json:"budget,omitempty"`
	Spends               map[int32]*RecordSpend `protobuf:"bytes,3,rep,name=spends,proto3" json:"spends,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	LastSpendUpdate      int64                  `protobuf:"varint,4,opt,name=last_spend_update,json=lastSpendUpdate,proto3" json:"last_spend_update,omitempty"`
	LastPush             int64                  `protobuf:"varint,5,opt,name=last_push,json=lastPush,proto3" json:"last_push,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_033e63fc596348d2, []int{0}
}

func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (m *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(m, src)
}
func (m *Config) XXX_Size() int {
	return xxx_messageInfo_Config.Size(m)
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

func (m *Config) GetWants() []*MasterWant {
	if m != nil {
		return m.Wants
	}
	return nil
}

func (m *Config) GetBudget() int32 {
	if m != nil {
		return m.Budget
	}
	return 0
}

func (m *Config) GetSpends() map[int32]*RecordSpend {
	if m != nil {
		return m.Spends
	}
	return nil
}

func (m *Config) GetLastSpendUpdate() int64 {
	if m != nil {
		return m.LastSpendUpdate
	}
	return 0
}

func (m *Config) GetLastPush() int64 {
	if m != nil {
		return m.LastPush
	}
	return 0
}

type RecordSpend struct {
	Cost                 int32    `protobuf:"varint,1,opt,name=cost,proto3" json:"cost,omitempty"`
	DateAdded            int64    `protobuf:"varint,2,opt,name=date_added,json=dateAdded,proto3" json:"date_added,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RecordSpend) Reset()         { *m = RecordSpend{} }
func (m *RecordSpend) String() string { return proto.CompactTextString(m) }
func (*RecordSpend) ProtoMessage()    {}
func (*RecordSpend) Descriptor() ([]byte, []int) {
	return fileDescriptor_033e63fc596348d2, []int{1}
}

func (m *RecordSpend) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RecordSpend.Unmarshal(m, b)
}
func (m *RecordSpend) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RecordSpend.Marshal(b, m, deterministic)
}
func (m *RecordSpend) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RecordSpend.Merge(m, src)
}
func (m *RecordSpend) XXX_Size() int {
	return xxx_messageInfo_RecordSpend.Size(m)
}
func (m *RecordSpend) XXX_DiscardUnknown() {
	xxx_messageInfo_RecordSpend.DiscardUnknown(m)
}

var xxx_messageInfo_RecordSpend proto.InternalMessageInfo

func (m *RecordSpend) GetCost() int32 {
	if m != nil {
		return m.Cost
	}
	return 0
}

func (m *RecordSpend) GetDateAdded() int64 {
	if m != nil {
		return m.DateAdded
	}
	return 0
}

type MasterWant struct {
	Release              *godiscogs.Release `protobuf:"bytes,1,opt,name=release,proto3" json:"release,omitempty"`
	DateAdded            int64              `protobuf:"varint,2,opt,name=date_added,json=dateAdded,proto3" json:"date_added,omitempty"`
	Staged               bool               `protobuf:"varint,3,opt,name=staged,proto3" json:"staged,omitempty"`
	Active               bool               `protobuf:"varint,4,opt,name=active,proto3" json:"active,omitempty"`
	Demoted              bool               `protobuf:"varint,5,opt,name=demoted,proto3" json:"demoted,omitempty"`
	Superwant            bool               `protobuf:"varint,6,opt,name=superwant,proto3" json:"superwant,omitempty"`
	Level                MasterWant_Level   `protobuf:"varint,7,opt,name=level,proto3,enum=recordwants.MasterWant_Level" json:"level,omitempty"`
	DatePurchased        int64              `protobuf:"varint,8,opt,name=date_purchased,json=datePurchased,proto3" json:"date_purchased,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *MasterWant) Reset()         { *m = MasterWant{} }
func (m *MasterWant) String() string { return proto.CompactTextString(m) }
func (*MasterWant) ProtoMessage()    {}
func (*MasterWant) Descriptor() ([]byte, []int) {
	return fileDescriptor_033e63fc596348d2, []int{2}
}

func (m *MasterWant) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MasterWant.Unmarshal(m, b)
}
func (m *MasterWant) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MasterWant.Marshal(b, m, deterministic)
}
func (m *MasterWant) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MasterWant.Merge(m, src)
}
func (m *MasterWant) XXX_Size() int {
	return xxx_messageInfo_MasterWant.Size(m)
}
func (m *MasterWant) XXX_DiscardUnknown() {
	xxx_messageInfo_MasterWant.DiscardUnknown(m)
}

var xxx_messageInfo_MasterWant proto.InternalMessageInfo

func (m *MasterWant) GetRelease() *godiscogs.Release {
	if m != nil {
		return m.Release
	}
	return nil
}

func (m *MasterWant) GetDateAdded() int64 {
	if m != nil {
		return m.DateAdded
	}
	return 0
}

func (m *MasterWant) GetStaged() bool {
	if m != nil {
		return m.Staged
	}
	return false
}

func (m *MasterWant) GetActive() bool {
	if m != nil {
		return m.Active
	}
	return false
}

func (m *MasterWant) GetDemoted() bool {
	if m != nil {
		return m.Demoted
	}
	return false
}

func (m *MasterWant) GetSuperwant() bool {
	if m != nil {
		return m.Superwant
	}
	return false
}

func (m *MasterWant) GetLevel() MasterWant_Level {
	if m != nil {
		return m.Level
	}
	return MasterWant_UNKNOWN
}

func (m *MasterWant) GetDatePurchased() int64 {
	if m != nil {
		return m.DatePurchased
	}
	return 0
}

type Spend struct {
	Month                int32    `protobuf:"varint,1,opt,name=month,proto3" json:"month,omitempty"`
	Spend                int32    `protobuf:"varint,2,opt,name=spend,proto3" json:"spend,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Spend) Reset()         { *m = Spend{} }
func (m *Spend) String() string { return proto.CompactTextString(m) }
func (*Spend) ProtoMessage()    {}
func (*Spend) Descriptor() ([]byte, []int) {
	return fileDescriptor_033e63fc596348d2, []int{3}
}

func (m *Spend) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Spend.Unmarshal(m, b)
}
func (m *Spend) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Spend.Marshal(b, m, deterministic)
}
func (m *Spend) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Spend.Merge(m, src)
}
func (m *Spend) XXX_Size() int {
	return xxx_messageInfo_Spend.Size(m)
}
func (m *Spend) XXX_DiscardUnknown() {
	xxx_messageInfo_Spend.DiscardUnknown(m)
}

var xxx_messageInfo_Spend proto.InternalMessageInfo

func (m *Spend) GetMonth() int32 {
	if m != nil {
		return m.Month
	}
	return 0
}

func (m *Spend) GetSpend() int32 {
	if m != nil {
		return m.Spend
	}
	return 0
}

type SpendingRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SpendingRequest) Reset()         { *m = SpendingRequest{} }
func (m *SpendingRequest) String() string { return proto.CompactTextString(m) }
func (*SpendingRequest) ProtoMessage()    {}
func (*SpendingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_033e63fc596348d2, []int{4}
}

func (m *SpendingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SpendingRequest.Unmarshal(m, b)
}
func (m *SpendingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SpendingRequest.Marshal(b, m, deterministic)
}
func (m *SpendingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SpendingRequest.Merge(m, src)
}
func (m *SpendingRequest) XXX_Size() int {
	return xxx_messageInfo_SpendingRequest.Size(m)
}
func (m *SpendingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SpendingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SpendingRequest proto.InternalMessageInfo

type SpendingResponse struct {
	Spends               []*Spend `protobuf:"bytes,1,rep,name=spends,proto3" json:"spends,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SpendingResponse) Reset()         { *m = SpendingResponse{} }
func (m *SpendingResponse) String() string { return proto.CompactTextString(m) }
func (*SpendingResponse) ProtoMessage()    {}
func (*SpendingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_033e63fc596348d2, []int{5}
}

func (m *SpendingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SpendingResponse.Unmarshal(m, b)
}
func (m *SpendingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SpendingResponse.Marshal(b, m, deterministic)
}
func (m *SpendingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SpendingResponse.Merge(m, src)
}
func (m *SpendingResponse) XXX_Size() int {
	return xxx_messageInfo_SpendingResponse.Size(m)
}
func (m *SpendingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SpendingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SpendingResponse proto.InternalMessageInfo

func (m *SpendingResponse) GetSpends() []*Spend {
	if m != nil {
		return m.Spends
	}
	return nil
}

type UpdateRequest struct {
	Want                 *godiscogs.Release `protobuf:"bytes,1,opt,name=want,proto3" json:"want,omitempty"`
	KeepWant             bool               `protobuf:"varint,2,opt,name=keep_want,json=keepWant,proto3" json:"keep_want,omitempty"`
	Super                bool               `protobuf:"varint,3,opt,name=super,proto3" json:"super,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *UpdateRequest) Reset()         { *m = UpdateRequest{} }
func (m *UpdateRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateRequest) ProtoMessage()    {}
func (*UpdateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_033e63fc596348d2, []int{6}
}

func (m *UpdateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateRequest.Unmarshal(m, b)
}
func (m *UpdateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateRequest.Marshal(b, m, deterministic)
}
func (m *UpdateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateRequest.Merge(m, src)
}
func (m *UpdateRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateRequest.Size(m)
}
func (m *UpdateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateRequest proto.InternalMessageInfo

func (m *UpdateRequest) GetWant() *godiscogs.Release {
	if m != nil {
		return m.Want
	}
	return nil
}

func (m *UpdateRequest) GetKeepWant() bool {
	if m != nil {
		return m.KeepWant
	}
	return false
}

func (m *UpdateRequest) GetSuper() bool {
	if m != nil {
		return m.Super
	}
	return false
}

type UpdateResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateResponse) Reset()         { *m = UpdateResponse{} }
func (m *UpdateResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateResponse) ProtoMessage()    {}
func (*UpdateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_033e63fc596348d2, []int{7}
}

func (m *UpdateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateResponse.Unmarshal(m, b)
}
func (m *UpdateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateResponse.Marshal(b, m, deterministic)
}
func (m *UpdateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateResponse.Merge(m, src)
}
func (m *UpdateResponse) XXX_Size() int {
	return xxx_messageInfo_UpdateResponse.Size(m)
}
func (m *UpdateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateResponse proto.InternalMessageInfo

type AddWantRequest struct {
	ReleaseId            int32    `protobuf:"varint,1,opt,name=release_id,json=releaseId,proto3" json:"release_id,omitempty"`
	Superwant            bool     `protobuf:"varint,2,opt,name=superwant,proto3" json:"superwant,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddWantRequest) Reset()         { *m = AddWantRequest{} }
func (m *AddWantRequest) String() string { return proto.CompactTextString(m) }
func (*AddWantRequest) ProtoMessage()    {}
func (*AddWantRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_033e63fc596348d2, []int{8}
}

func (m *AddWantRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddWantRequest.Unmarshal(m, b)
}
func (m *AddWantRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddWantRequest.Marshal(b, m, deterministic)
}
func (m *AddWantRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddWantRequest.Merge(m, src)
}
func (m *AddWantRequest) XXX_Size() int {
	return xxx_messageInfo_AddWantRequest.Size(m)
}
func (m *AddWantRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddWantRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddWantRequest proto.InternalMessageInfo

func (m *AddWantRequest) GetReleaseId() int32 {
	if m != nil {
		return m.ReleaseId
	}
	return 0
}

func (m *AddWantRequest) GetSuperwant() bool {
	if m != nil {
		return m.Superwant
	}
	return false
}

type AddWantResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddWantResponse) Reset()         { *m = AddWantResponse{} }
func (m *AddWantResponse) String() string { return proto.CompactTextString(m) }
func (*AddWantResponse) ProtoMessage()    {}
func (*AddWantResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_033e63fc596348d2, []int{9}
}

func (m *AddWantResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddWantResponse.Unmarshal(m, b)
}
func (m *AddWantResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddWantResponse.Marshal(b, m, deterministic)
}
func (m *AddWantResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddWantResponse.Merge(m, src)
}
func (m *AddWantResponse) XXX_Size() int {
	return xxx_messageInfo_AddWantResponse.Size(m)
}
func (m *AddWantResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddWantResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddWantResponse proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("recordwants.MasterWant_Level", MasterWant_Level_name, MasterWant_Level_value)
	proto.RegisterType((*Config)(nil), "recordwants.Config")
	proto.RegisterMapType((map[int32]*RecordSpend)(nil), "recordwants.Config.SpendsEntry")
	proto.RegisterType((*RecordSpend)(nil), "recordwants.RecordSpend")
	proto.RegisterType((*MasterWant)(nil), "recordwants.MasterWant")
	proto.RegisterType((*Spend)(nil), "recordwants.Spend")
	proto.RegisterType((*SpendingRequest)(nil), "recordwants.SpendingRequest")
	proto.RegisterType((*SpendingResponse)(nil), "recordwants.SpendingResponse")
	proto.RegisterType((*UpdateRequest)(nil), "recordwants.UpdateRequest")
	proto.RegisterType((*UpdateResponse)(nil), "recordwants.UpdateResponse")
	proto.RegisterType((*AddWantRequest)(nil), "recordwants.AddWantRequest")
	proto.RegisterType((*AddWantResponse)(nil), "recordwants.AddWantResponse")
}

func init() { proto.RegisterFile("recordwants.proto", fileDescriptor_033e63fc596348d2) }

var fileDescriptor_033e63fc596348d2 = []byte{
	// 687 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0xed, 0x4e, 0xdb, 0x48,
	0x14, 0xc5, 0x49, 0x9c, 0x8f, 0x6b, 0x11, 0xc2, 0x68, 0xb5, 0x6b, 0x05, 0xa2, 0x8d, 0x2c, 0xed,
	0x2a, 0x42, 0x6d, 0x50, 0xc3, 0x0f, 0xaa, 0xfe, 0xa8, 0x1a, 0x21, 0x5a, 0xa1, 0x86, 0x14, 0x4d,
	0x40, 0x88, 0x5f, 0x91, 0xe3, 0x99, 0x3a, 0x2e, 0xc1, 0x76, 0x3d, 0xe3, 0x54, 0xbc, 0x50, 0x1f,
	0xa9, 0xcf, 0xd0, 0xc7, 0xa8, 0xe6, 0xce, 0x84, 0xc4, 0x14, 0xd4, 0x7f, 0x73, 0xcf, 0x3d, 0x73,
	0xe6, 0xe4, 0xdc, 0xeb, 0xc0, 0x6e, 0xc6, 0x83, 0x24, 0x63, 0xdf, 0xfc, 0x58, 0x8a, 0x7e, 0x9a,
	0x25, 0x32, 0x21, 0xce, 0x06, 0xd4, 0x7e, 0x15, 0x46, 0x72, 0x9e, 0xcf, 0xfa, 0x41, 0x72, 0x77,
	0x38, 0xcb, 0x12, 0x39, 0xe7, 0xd9, 0x22, 0x09, 0xa3, 0xe0, 0x30, 0x4c, 0x58, 0x24, 0x82, 0x24,
	0x14, 0xeb, 0x93, 0xbe, 0xef, 0x7d, 0x2f, 0x41, 0xf5, 0x24, 0x89, 0x3f, 0x47, 0x21, 0x79, 0x09,
	0x36, 0xca, 0xb8, 0x56, 0xb7, 0xdc, 0x73, 0x06, 0xff, 0xf4, 0x37, 0x5f, 0x3b, 0xf7, 0x85, 0xe4,
	0xd9, 0xb5, 0x1f, 0x4b, 0xaa, 0x59, 0xe4, 0x6f, 0xa8, 0xce, 0x72, 0x16, 0x72, 0xe9, 0x96, 0xba,
	0x56, 0xcf, 0xa6, 0xa6, 0x22, 0xc7, 0x50, 0x15, 0x29, 0x8f, 0x99, 0x70, 0xcb, 0xa8, 0xf3, 0x6f,
	0x41, 0x47, 0xbf, 0xd5, 0x9f, 0x20, 0xe3, 0x34, 0x96, 0xd9, 0x3d, 0x35, 0x74, 0x72, 0x00, 0xbb,
	0x0b, 0x5f, 0xc8, 0x29, 0x96, 0xd3, 0x3c, 0x65, 0xbe, 0xe4, 0x6e, 0xa5, 0x6b, 0xf5, 0xca, 0x74,
	0x47, 0x35, 0xf0, 0xce, 0x15, 0xc2, 0x64, 0x0f, 0x1a, 0xc8, 0x4d, 0x73, 0x31, 0x77, 0x6d, 0xe4,
	0xd4, 0x15, 0x70, 0x91, 0x8b, 0x79, 0x7b, 0x02, 0xce, 0x86, 0x3e, 0x69, 0x41, 0xf9, 0x96, 0xdf,
	0xbb, 0x16, 0xba, 0x54, 0x47, 0xd2, 0x07, 0x7b, 0xe9, 0x2f, 0x72, 0x8e, 0xce, 0x9d, 0x81, 0x5b,
	0x70, 0x48, 0xf1, 0x8c, 0x02, 0x54, 0xd3, 0xde, 0x94, 0x5e, 0x5b, 0xde, 0x3b, 0x70, 0x36, 0x3a,
	0x84, 0x40, 0x25, 0x48, 0x84, 0x34, 0xaa, 0x78, 0x26, 0x1d, 0x00, 0x65, 0x6e, 0xea, 0x33, 0xc6,
	0x19, 0x6a, 0x97, 0x69, 0x43, 0x21, 0x43, 0x05, 0x78, 0x3f, 0x4a, 0x00, 0xeb, 0x18, 0xc9, 0x0b,
	0xa8, 0x65, 0x7c, 0xc1, 0x7d, 0xc1, 0x51, 0xc4, 0x19, 0x90, 0xfe, 0x7a, 0x38, 0x54, 0x77, 0xe8,
	0x8a, 0xf2, 0x07, 0x6d, 0x35, 0x0c, 0x21, 0xfd, 0x90, 0x33, 0xb7, 0xdc, 0xb5, 0x7a, 0x75, 0x6a,
	0x2a, 0x85, 0xfb, 0x81, 0x8c, 0x96, 0x3a, 0xc8, 0x3a, 0x35, 0x15, 0x71, 0xa1, 0xc6, 0xf8, 0x5d,
	0x22, 0x39, 0xc3, 0xf4, 0xea, 0x74, 0x55, 0x92, 0x7d, 0x68, 0x88, 0x3c, 0xe5, 0x99, 0x0a, 0xc3,
	0xad, 0x62, 0x6f, 0x0d, 0x90, 0x23, 0xb0, 0x17, 0x7c, 0xc9, 0x17, 0x6e, 0xad, 0x6b, 0xf5, 0x9a,
	0x83, 0xce, 0x33, 0x3b, 0xd2, 0x1f, 0x29, 0x12, 0xd5, 0x5c, 0xf2, 0x1f, 0x34, 0xd1, 0x7b, 0x9a,
	0x67, 0xc1, 0xdc, 0x17, 0x9c, 0xb9, 0x75, 0xf4, 0xbf, 0xad, 0xd0, 0x8b, 0x15, 0xe8, 0x1d, 0x83,
	0x8d, 0xd7, 0x88, 0x03, 0xb5, 0xab, 0xf1, 0xc7, 0xf1, 0xa7, 0xeb, 0x71, 0x6b, 0x4b, 0x15, 0xc3,
	0xf1, 0xcd, 0xe5, 0xd9, 0xf9, 0x69, 0xcb, 0x22, 0x75, 0xa8, 0x8c, 0xce, 0x26, 0x97, 0xad, 0x12,
	0x01, 0xa8, 0x0e, 0x47, 0xd7, 0xc3, 0x9b, 0x49, 0xab, 0xec, 0x1d, 0x81, 0xad, 0x87, 0xf2, 0x17,
	0xd8, 0x77, 0x49, 0x2c, 0xe7, 0x66, 0x2a, 0xba, 0x50, 0x28, 0xae, 0x94, 0xd9, 0x53, 0x5d, 0x78,
	0xbb, 0xb0, 0x83, 0x97, 0xa2, 0x38, 0xa4, 0xfc, 0x6b, 0xce, 0x85, 0xf4, 0xde, 0x42, 0x6b, 0x0d,
	0x89, 0x34, 0x89, 0x05, 0x27, 0x07, 0x0f, 0xdb, 0xac, 0xbf, 0x0a, 0x52, 0xf8, 0xc5, 0x7a, 0x4b,
	0x0c, 0xc3, 0xfb, 0x02, 0xdb, 0x7a, 0x3d, 0x8d, 0x20, 0xf9, 0x1f, 0x2a, 0x18, 0xe3, 0xf3, 0xf3,
	0xc5, 0xbe, 0xda, 0xe6, 0x5b, 0xce, 0xd3, 0x29, 0x92, 0x4b, 0x98, 0x79, 0x5d, 0x01, 0xb8, 0x27,
	0xca, 0xbe, 0xca, 0xdf, 0x4c, 0x56, 0x17, 0x5e, 0x0b, 0x9a, 0xab, 0xb7, 0xb4, 0x53, 0xef, 0x1c,
	0x9a, 0x43, 0xc6, 0xf0, 0x0b, 0x35, 0xcf, 0x77, 0x00, 0xcc, 0xfa, 0x4c, 0x23, 0x66, 0x32, 0x69,
	0x18, 0xe4, 0xec, 0xd1, 0xa4, 0x4b, 0x8f, 0x26, 0xad, 0xf2, 0x79, 0x90, 0xd3, 0x2f, 0x0c, 0x7e,
	0x5a, 0xe0, 0x28, 0x60, 0xc2, 0xb3, 0x65, 0x14, 0x70, 0x32, 0x02, 0xe7, 0x03, 0x97, 0xab, 0xc8,
	0xc8, 0xfe, 0xef, 0xd1, 0xac, 0xc3, 0x6d, 0x77, 0x9e, 0xe9, 0x1a, 0xf7, 0x5b, 0xe4, 0x04, 0xaa,
	0xe6, 0xe3, 0x6e, 0x17, 0xa8, 0x85, 0x48, 0xdb, 0x7b, 0x4f, 0xf6, 0x1e, 0x44, 0xde, 0x43, 0xcd,
	0xb8, 0x26, 0x45, 0x66, 0x31, 0x9a, 0xf6, 0xfe, 0xd3, 0xcd, 0x95, 0xce, 0xac, 0x8a, 0xff, 0x8e,
	0x47, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0xc1, 0x6a, 0x98, 0x44, 0x72, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// WantServiceClient is the client API for WantService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WantServiceClient interface {
	GetSpending(ctx context.Context, in *SpendingRequest, opts ...grpc.CallOption) (*SpendingResponse, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
	AddWant(ctx context.Context, in *AddWantRequest, opts ...grpc.CallOption) (*AddWantResponse, error)
}

type wantServiceClient struct {
	cc *grpc.ClientConn
}

func NewWantServiceClient(cc *grpc.ClientConn) WantServiceClient {
	return &wantServiceClient{cc}
}

func (c *wantServiceClient) GetSpending(ctx context.Context, in *SpendingRequest, opts ...grpc.CallOption) (*SpendingResponse, error) {
	out := new(SpendingResponse)
	err := c.cc.Invoke(ctx, "/recordwants.WantService/GetSpending", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *wantServiceClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, "/recordwants.WantService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *wantServiceClient) AddWant(ctx context.Context, in *AddWantRequest, opts ...grpc.CallOption) (*AddWantResponse, error) {
	out := new(AddWantResponse)
	err := c.cc.Invoke(ctx, "/recordwants.WantService/AddWant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WantServiceServer is the server API for WantService service.
type WantServiceServer interface {
	GetSpending(context.Context, *SpendingRequest) (*SpendingResponse, error)
	Update(context.Context, *UpdateRequest) (*UpdateResponse, error)
	AddWant(context.Context, *AddWantRequest) (*AddWantResponse, error)
}

func RegisterWantServiceServer(s *grpc.Server, srv WantServiceServer) {
	s.RegisterService(&_WantService_serviceDesc, srv)
}

func _WantService_GetSpending_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SpendingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WantServiceServer).GetSpending(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recordwants.WantService/GetSpending",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WantServiceServer).GetSpending(ctx, req.(*SpendingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WantService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WantServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recordwants.WantService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WantServiceServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WantService_AddWant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddWantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WantServiceServer).AddWant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recordwants.WantService/AddWant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WantServiceServer).AddWant(ctx, req.(*AddWantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _WantService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "recordwants.WantService",
	HandlerType: (*WantServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSpending",
			Handler:    _WantService_GetSpending_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _WantService_Update_Handler,
		},
		{
			MethodName: "AddWant",
			Handler:    _WantService_AddWant_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "recordwants.proto",
}
