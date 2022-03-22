// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.1
// source: api/types/job.proto

package types

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Job struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Steps     []*Step           `protobuf:"bytes,1,rep,name=steps,proto3" json:"steps,omitempty"`
	RunsOn    string            `protobuf:"bytes,2,opt,name=runs_on,json=runsOn,proto3" json:"runs_on,omitempty"`
	Container *Job_Container    `protobuf:"bytes,3,opt,name=container,proto3" json:"container,omitempty"`
	Name      string            `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Outputs   map[string]string `protobuf:"bytes,5,rep,name=outputs,proto3" json:"outputs,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Env       map[string]string `protobuf:"bytes,6,rep,name=env,proto3" json:"env,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Job) Reset() {
	*x = Job{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_types_job_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Job) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Job) ProtoMessage() {}

func (x *Job) ProtoReflect() protoreflect.Message {
	mi := &file_api_types_job_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Job.ProtoReflect.Descriptor instead.
func (*Job) Descriptor() ([]byte, []int) {
	return file_api_types_job_proto_rawDescGZIP(), []int{0}
}

func (x *Job) GetSteps() []*Step {
	if x != nil {
		return x.Steps
	}
	return nil
}

func (x *Job) GetRunsOn() string {
	if x != nil {
		return x.RunsOn
	}
	return ""
}

func (x *Job) GetContainer() *Job_Container {
	if x != nil {
		return x.Container
	}
	return nil
}

func (x *Job) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Job) GetOutputs() map[string]string {
	if x != nil {
		return x.Outputs
	}
	return nil
}

func (x *Job) GetEnv() map[string]string {
	if x != nil {
		return x.Env
	}
	return nil
}

type Job_Container struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Image string `protobuf:"bytes,1,opt,name=image,proto3" json:"image,omitempty"`
}

func (x *Job_Container) Reset() {
	*x = Job_Container{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_types_job_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Job_Container) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Job_Container) ProtoMessage() {}

func (x *Job_Container) ProtoReflect() protoreflect.Message {
	mi := &file_api_types_job_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Job_Container.ProtoReflect.Descriptor instead.
func (*Job_Container) Descriptor() ([]byte, []int) {
	return file_api_types_job_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Job_Container) GetImage() string {
	if x != nil {
		return x.Image
	}
	return ""
}

var File_api_types_job_proto protoreflect.FileDescriptor

var file_api_types_job_proto_rawDesc = []byte{
	0x0a, 0x13, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x2e,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x1a, 0x14, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2f, 0x73, 0x74, 0x65, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9e, 0x03, 0x0a, 0x03,
	0x4a, 0x6f, 0x62, 0x12, 0x2a, 0x0a, 0x05, 0x73, 0x74, 0x65, 0x70, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2e, 0x53, 0x74, 0x65, 0x70, 0x52, 0x05, 0x73, 0x74, 0x65, 0x70, 0x73, 0x12,
	0x17, 0x0a, 0x07, 0x72, 0x75, 0x6e, 0x73, 0x5f, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x72, 0x75, 0x6e, 0x73, 0x4f, 0x6e, 0x12, 0x3b, 0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x74,
	0x61, 0x69, 0x6e, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x73, 0x65,
	0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4a, 0x6f, 0x62,
	0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x74,
	0x61, 0x69, 0x6e, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3a, 0x0a, 0x07, 0x6f, 0x75, 0x74,
	0x70, 0x75, 0x74, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x73, 0x65, 0x71,
	0x75, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4a, 0x6f, 0x62, 0x2e,
	0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x6f, 0x75,
	0x74, 0x70, 0x75, 0x74, 0x73, 0x12, 0x2e, 0x0a, 0x03, 0x65, 0x6e, 0x76, 0x18, 0x06, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2e, 0x4a, 0x6f, 0x62, 0x2e, 0x45, 0x6e, 0x76, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x03, 0x65, 0x6e, 0x76, 0x1a, 0x21, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e,
	0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x1a, 0x3a, 0x0a, 0x0c, 0x4f, 0x75, 0x74, 0x70,
	0x75, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x1a, 0x36, 0x0a, 0x08, 0x45, 0x6e, 0x76, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x27, 0x5a, 0x25,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x72, 0x61, 0x6e, 0x74,
	0x6a, 0x63, 0x2f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_types_job_proto_rawDescOnce sync.Once
	file_api_types_job_proto_rawDescData = file_api_types_job_proto_rawDesc
)

func file_api_types_job_proto_rawDescGZIP() []byte {
	file_api_types_job_proto_rawDescOnce.Do(func() {
		file_api_types_job_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_types_job_proto_rawDescData)
	})
	return file_api_types_job_proto_rawDescData
}

var file_api_types_job_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_types_job_proto_goTypes = []interface{}{
	(*Job)(nil),           // 0: sequence.types.Job
	(*Job_Container)(nil), // 1: sequence.types.Job.Container
	nil,                   // 2: sequence.types.Job.OutputsEntry
	nil,                   // 3: sequence.types.Job.EnvEntry
	(*Step)(nil),          // 4: sequence.types.Step
}
var file_api_types_job_proto_depIdxs = []int32{
	4, // 0: sequence.types.Job.steps:type_name -> sequence.types.Step
	1, // 1: sequence.types.Job.container:type_name -> sequence.types.Job.Container
	2, // 2: sequence.types.Job.outputs:type_name -> sequence.types.Job.OutputsEntry
	3, // 3: sequence.types.Job.env:type_name -> sequence.types.Job.EnvEntry
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_api_types_job_proto_init() }
func file_api_types_job_proto_init() {
	if File_api_types_job_proto != nil {
		return
	}
	file_api_types_step_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_types_job_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Job); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_types_job_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Job_Container); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_types_job_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_types_job_proto_goTypes,
		DependencyIndexes: file_api_types_job_proto_depIdxs,
		MessageInfos:      file_api_types_job_proto_msgTypes,
	}.Build()
	File_api_types_job_proto = out.File
	file_api_types_job_proto_rawDesc = nil
	file_api_types_job_proto_goTypes = nil
	file_api_types_job_proto_depIdxs = nil
}
