// Code generated by protoc-gen-go. DO NOT EDIT.
// source: presslabs/dashboard/projects/v1/project.proto

package v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import empty "github.com/golang/protobuf/ptypes/empty"
import field_mask "google.golang.org/genproto/protobuf/field_mask"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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

// Projects represents an project within the presslabs dashboard
// context
type Project struct {
	// The fully qualified project name in the form proj/{project_name}.
	// The `project_name` is a valid DNS label (RFC 1123) with maximum
	// length of 48 characters.
	// The name is read-only
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The display_name for the project
	DisplayName string `protobuf:"bytes,2,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	// The organization name. The organization is immutable.
	// This field is read-only
	Organization         string   `protobuf:"bytes,3,opt,name=organization,proto3" json:"organization,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Project) Reset()         { *m = Project{} }
func (m *Project) String() string { return proto.CompactTextString(m) }
func (*Project) ProtoMessage()    {}
func (*Project) Descriptor() ([]byte, []int) {
	return fileDescriptor_project_60b175e306b49eb9, []int{0}
}
func (m *Project) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Project.Unmarshal(m, b)
}
func (m *Project) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Project.Marshal(b, m, deterministic)
}
func (dst *Project) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Project.Merge(dst, src)
}
func (m *Project) XXX_Size() int {
	return xxx_messageInfo_Project.Size(m)
}
func (m *Project) XXX_DiscardUnknown() {
	xxx_messageInfo_Project.DiscardUnknown(m)
}

var xxx_messageInfo_Project proto.InternalMessageInfo

func (m *Project) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Project) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

func (m *Project) GetOrganization() string {
	if m != nil {
		return m.Organization
	}
	return ""
}

type GetProjectRequest struct {
	// The resource name of the project to fetch in the form proj/{project_name}
	// The `project_name` MUST be a valid DNS label (RFC 1123)
	// with maximum length of 48 characters.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetProjectRequest) Reset()         { *m = GetProjectRequest{} }
func (m *GetProjectRequest) String() string { return proto.CompactTextString(m) }
func (*GetProjectRequest) ProtoMessage()    {}
func (*GetProjectRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_project_60b175e306b49eb9, []int{1}
}
func (m *GetProjectRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetProjectRequest.Unmarshal(m, b)
}
func (m *GetProjectRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetProjectRequest.Marshal(b, m, deterministic)
}
func (dst *GetProjectRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetProjectRequest.Merge(dst, src)
}
func (m *GetProjectRequest) XXX_Size() int {
	return xxx_messageInfo_GetProjectRequest.Size(m)
}
func (m *GetProjectRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetProjectRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetProjectRequest proto.InternalMessageInfo

func (m *GetProjectRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type ListProjectsRequest struct {
	// The maximum number of items to return.
	PageSize int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// The next_page_token value returned from a previous List request, if
	// any.
	PageToken            string   `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListProjectsRequest) Reset()         { *m = ListProjectsRequest{} }
func (m *ListProjectsRequest) String() string { return proto.CompactTextString(m) }
func (*ListProjectsRequest) ProtoMessage()    {}
func (*ListProjectsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_project_60b175e306b49eb9, []int{2}
}
func (m *ListProjectsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListProjectsRequest.Unmarshal(m, b)
}
func (m *ListProjectsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListProjectsRequest.Marshal(b, m, deterministic)
}
func (dst *ListProjectsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListProjectsRequest.Merge(dst, src)
}
func (m *ListProjectsRequest) XXX_Size() int {
	return xxx_messageInfo_ListProjectsRequest.Size(m)
}
func (m *ListProjectsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListProjectsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListProjectsRequest proto.InternalMessageInfo

func (m *ListProjectsRequest) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *ListProjectsRequest) GetPageToken() string {
	if m != nil {
		return m.PageToken
	}
	return ""
}

type ListProjectsResponse struct {
	Projects []*Project `protobuf:"bytes,1,rep,name=projects,proto3" json:"projects,omitempty"`
	// Token to retrieve the next page of results, or empty if there are no
	// more results in the list.
	NextPageToken        string   `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListProjectsResponse) Reset()         { *m = ListProjectsResponse{} }
func (m *ListProjectsResponse) String() string { return proto.CompactTextString(m) }
func (*ListProjectsResponse) ProtoMessage()    {}
func (*ListProjectsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_project_60b175e306b49eb9, []int{3}
}
func (m *ListProjectsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListProjectsResponse.Unmarshal(m, b)
}
func (m *ListProjectsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListProjectsResponse.Marshal(b, m, deterministic)
}
func (dst *ListProjectsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListProjectsResponse.Merge(dst, src)
}
func (m *ListProjectsResponse) XXX_Size() int {
	return xxx_messageInfo_ListProjectsResponse.Size(m)
}
func (m *ListProjectsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListProjectsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListProjectsResponse proto.InternalMessageInfo

func (m *ListProjectsResponse) GetProjects() []*Project {
	if m != nil {
		return m.Projects
	}
	return nil
}

func (m *ListProjectsResponse) GetNextPageToken() string {
	if m != nil {
		return m.NextPageToken
	}
	return ""
}

type CreateProjectRequest struct {
	// The parent resource name where the project is to be created
	// The parent is a required parameter
	Parent string `protobuf:"bytes,1,opt,name=parent,proto3" json:"parent,omitempty"`
	// The project name to assign
	ProjectId string `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	// The project resource to create
	Project              *Project `protobuf:"bytes,3,opt,name=project,proto3" json:"project,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateProjectRequest) Reset()         { *m = CreateProjectRequest{} }
func (m *CreateProjectRequest) String() string { return proto.CompactTextString(m) }
func (*CreateProjectRequest) ProtoMessage()    {}
func (*CreateProjectRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_project_60b175e306b49eb9, []int{4}
}
func (m *CreateProjectRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateProjectRequest.Unmarshal(m, b)
}
func (m *CreateProjectRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateProjectRequest.Marshal(b, m, deterministic)
}
func (dst *CreateProjectRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateProjectRequest.Merge(dst, src)
}
func (m *CreateProjectRequest) XXX_Size() int {
	return xxx_messageInfo_CreateProjectRequest.Size(m)
}
func (m *CreateProjectRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateProjectRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateProjectRequest proto.InternalMessageInfo

func (m *CreateProjectRequest) GetParent() string {
	if m != nil {
		return m.Parent
	}
	return ""
}

func (m *CreateProjectRequest) GetProjectId() string {
	if m != nil {
		return m.ProjectId
	}
	return ""
}

func (m *CreateProjectRequest) GetProject() *Project {
	if m != nil {
		return m.Project
	}
	return nil
}

type UpdateProjectRequest struct {
	// The new definition of the Folder. It must include
	// a `name` , `organization` and `display_name` field.
	// The other fields will be ignored.
	Project *Project `protobuf:"bytes,1,opt,name=project,proto3" json:"project,omitempty"`
	// Fields to be updated.
	// Only the `display_name` can be updated.
	UpdateMask           *field_mask.FieldMask `protobuf:"bytes,2,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *UpdateProjectRequest) Reset()         { *m = UpdateProjectRequest{} }
func (m *UpdateProjectRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateProjectRequest) ProtoMessage()    {}
func (*UpdateProjectRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_project_60b175e306b49eb9, []int{5}
}
func (m *UpdateProjectRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateProjectRequest.Unmarshal(m, b)
}
func (m *UpdateProjectRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateProjectRequest.Marshal(b, m, deterministic)
}
func (dst *UpdateProjectRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateProjectRequest.Merge(dst, src)
}
func (m *UpdateProjectRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateProjectRequest.Size(m)
}
func (m *UpdateProjectRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateProjectRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateProjectRequest proto.InternalMessageInfo

func (m *UpdateProjectRequest) GetProject() *Project {
	if m != nil {
		return m.Project
	}
	return nil
}

func (m *UpdateProjectRequest) GetUpdateMask() *field_mask.FieldMask {
	if m != nil {
		return m.UpdateMask
	}
	return nil
}

type DeleteProjectRequest struct {
	// The resource name of the project to delete in the form projs/{project_name}
	// The `project_name` MUST be a valid DNS label (RFC 1123)
	// with maximum length of 48 characters.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteProjectRequest) Reset()         { *m = DeleteProjectRequest{} }
func (m *DeleteProjectRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteProjectRequest) ProtoMessage()    {}
func (*DeleteProjectRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_project_60b175e306b49eb9, []int{6}
}
func (m *DeleteProjectRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteProjectRequest.Unmarshal(m, b)
}
func (m *DeleteProjectRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteProjectRequest.Marshal(b, m, deterministic)
}
func (dst *DeleteProjectRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteProjectRequest.Merge(dst, src)
}
func (m *DeleteProjectRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteProjectRequest.Size(m)
}
func (m *DeleteProjectRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteProjectRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteProjectRequest proto.InternalMessageInfo

func (m *DeleteProjectRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*Project)(nil), "presslabs.dashboard.projects.v1.Project")
	proto.RegisterType((*GetProjectRequest)(nil), "presslabs.dashboard.projects.v1.GetProjectRequest")
	proto.RegisterType((*ListProjectsRequest)(nil), "presslabs.dashboard.projects.v1.ListProjectsRequest")
	proto.RegisterType((*ListProjectsResponse)(nil), "presslabs.dashboard.projects.v1.ListProjectsResponse")
	proto.RegisterType((*CreateProjectRequest)(nil), "presslabs.dashboard.projects.v1.CreateProjectRequest")
	proto.RegisterType((*UpdateProjectRequest)(nil), "presslabs.dashboard.projects.v1.UpdateProjectRequest")
	proto.RegisterType((*DeleteProjectRequest)(nil), "presslabs.dashboard.projects.v1.DeleteProjectRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ProjectsServiceClient is the client API for ProjectsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProjectsServiceClient interface {
	// CreateProject creates a new project
	CreateProject(ctx context.Context, in *CreateProjectRequest, opts ...grpc.CallOption) (*Project, error)
	// GetProject fetches an project by it's name
	GetProject(ctx context.Context, in *GetProjectRequest, opts ...grpc.CallOption) (*Project, error)
	// UpdateProject updates an project details
	UpdateProject(ctx context.Context, in *UpdateProjectRequest, opts ...grpc.CallOption) (*Project, error)
	// DeleteProject deletes an project by it's name
	DeleteProject(ctx context.Context, in *DeleteProjectRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// ListProjects list projects
	ListProjects(ctx context.Context, in *ListProjectsRequest, opts ...grpc.CallOption) (*ListProjectsResponse, error)
}

type projectsServiceClient struct {
	cc *grpc.ClientConn
}

func NewProjectsServiceClient(cc *grpc.ClientConn) ProjectsServiceClient {
	return &projectsServiceClient{cc}
}

func (c *projectsServiceClient) CreateProject(ctx context.Context, in *CreateProjectRequest, opts ...grpc.CallOption) (*Project, error) {
	out := new(Project)
	err := c.cc.Invoke(ctx, "/presslabs.dashboard.projects.v1.ProjectsService/CreateProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) GetProject(ctx context.Context, in *GetProjectRequest, opts ...grpc.CallOption) (*Project, error) {
	out := new(Project)
	err := c.cc.Invoke(ctx, "/presslabs.dashboard.projects.v1.ProjectsService/GetProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) UpdateProject(ctx context.Context, in *UpdateProjectRequest, opts ...grpc.CallOption) (*Project, error) {
	out := new(Project)
	err := c.cc.Invoke(ctx, "/presslabs.dashboard.projects.v1.ProjectsService/UpdateProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) DeleteProject(ctx context.Context, in *DeleteProjectRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/presslabs.dashboard.projects.v1.ProjectsService/DeleteProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) ListProjects(ctx context.Context, in *ListProjectsRequest, opts ...grpc.CallOption) (*ListProjectsResponse, error) {
	out := new(ListProjectsResponse)
	err := c.cc.Invoke(ctx, "/presslabs.dashboard.projects.v1.ProjectsService/ListProjects", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProjectsServiceServer is the server API for ProjectsService service.
type ProjectsServiceServer interface {
	// CreateProject creates a new project
	CreateProject(context.Context, *CreateProjectRequest) (*Project, error)
	// GetProject fetches an project by it's name
	GetProject(context.Context, *GetProjectRequest) (*Project, error)
	// UpdateProject updates an project details
	UpdateProject(context.Context, *UpdateProjectRequest) (*Project, error)
	// DeleteProject deletes an project by it's name
	DeleteProject(context.Context, *DeleteProjectRequest) (*empty.Empty, error)
	// ListProjects list projects
	ListProjects(context.Context, *ListProjectsRequest) (*ListProjectsResponse, error)
}

func RegisterProjectsServiceServer(s *grpc.Server, srv ProjectsServiceServer) {
	s.RegisterService(&_ProjectsService_serviceDesc, srv)
}

func _ProjectsService_CreateProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).CreateProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/presslabs.dashboard.projects.v1.ProjectsService/CreateProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).CreateProject(ctx, req.(*CreateProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_GetProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).GetProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/presslabs.dashboard.projects.v1.ProjectsService/GetProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).GetProject(ctx, req.(*GetProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_UpdateProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).UpdateProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/presslabs.dashboard.projects.v1.ProjectsService/UpdateProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).UpdateProject(ctx, req.(*UpdateProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_DeleteProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).DeleteProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/presslabs.dashboard.projects.v1.ProjectsService/DeleteProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).DeleteProject(ctx, req.(*DeleteProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_ListProjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListProjectsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).ListProjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/presslabs.dashboard.projects.v1.ProjectsService/ListProjects",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).ListProjects(ctx, req.(*ListProjectsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProjectsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "presslabs.dashboard.projects.v1.ProjectsService",
	HandlerType: (*ProjectsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateProject",
			Handler:    _ProjectsService_CreateProject_Handler,
		},
		{
			MethodName: "GetProject",
			Handler:    _ProjectsService_GetProject_Handler,
		},
		{
			MethodName: "UpdateProject",
			Handler:    _ProjectsService_UpdateProject_Handler,
		},
		{
			MethodName: "DeleteProject",
			Handler:    _ProjectsService_DeleteProject_Handler,
		},
		{
			MethodName: "ListProjects",
			Handler:    _ProjectsService_ListProjects_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "presslabs/dashboard/projects/v1/project.proto",
}

func init() {
	proto.RegisterFile("presslabs/dashboard/projects/v1/project.proto", fileDescriptor_project_60b175e306b49eb9)
}

var fileDescriptor_project_60b175e306b49eb9 = []byte{
	// 510 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0x4d, 0x6f, 0xd3, 0x4c,
	0x10, 0x96, 0xdb, 0xbc, 0xfd, 0x98, 0x24, 0xaa, 0xde, 0x25, 0xaa, 0x22, 0x57, 0x88, 0xe0, 0x03,
	0x44, 0x48, 0xd8, 0x4a, 0xa0, 0x27, 0x6e, 0xa5, 0x80, 0x90, 0x00, 0x95, 0x14, 0x2e, 0x1c, 0xb0,
	0x36, 0xf5, 0x34, 0x6c, 0xe3, 0x78, 0x17, 0xef, 0x26, 0xa2, 0xe1, 0xca, 0x1f, 0xe0, 0xc4, 0xbf,
	0xe1, 0xb7, 0x21, 0xef, 0x7a, 0x4d, 0x5d, 0x5b, 0x4a, 0xc2, 0x6d, 0x77, 0xe6, 0x99, 0x8f, 0x67,
	0x9e, 0xd9, 0x85, 0xc7, 0x22, 0x45, 0x29, 0x63, 0x3a, 0x96, 0x41, 0x44, 0xe5, 0x97, 0x31, 0xa7,
	0x69, 0x14, 0x88, 0x94, 0x5f, 0xe1, 0x85, 0x92, 0xc1, 0x62, 0x60, 0xcf, 0xbe, 0x48, 0xb9, 0xe2,
	0xe4, 0x5e, 0x01, 0xf7, 0x0b, 0xb8, 0x6f, 0xe1, 0xfe, 0x62, 0xe0, 0x1e, 0x4d, 0x38, 0x9f, 0xc4,
	0x18, 0x68, 0xf8, 0x78, 0x7e, 0x19, 0xe0, 0x4c, 0xa8, 0x6b, 0x13, 0xed, 0xf6, 0x6e, 0x3b, 0x2f,
	0x19, 0xc6, 0x51, 0x38, 0xa3, 0x72, 0x6a, 0x10, 0x5e, 0x04, 0xbb, 0x67, 0x26, 0x1b, 0x21, 0xd0,
	0x48, 0xe8, 0x0c, 0xbb, 0x4e, 0xcf, 0xe9, 0xef, 0x8f, 0xf4, 0x99, 0xdc, 0x87, 0x56, 0xc4, 0xa4,
	0x88, 0xe9, 0x75, 0xa8, 0x7d, 0x5b, 0xda, 0xd7, 0xcc, 0x6d, 0xef, 0x32, 0x88, 0x07, 0x2d, 0x9e,
	0x4e, 0x68, 0xc2, 0x96, 0x54, 0x31, 0x9e, 0x74, 0xb7, 0x35, 0xa4, 0x64, 0xf3, 0x1e, 0xc2, 0xff,
	0xaf, 0x50, 0xe5, 0x85, 0x46, 0xf8, 0x75, 0x8e, 0xb2, 0xb6, 0x9e, 0xf7, 0x1e, 0xee, 0xbc, 0x61,
	0xd2, 0x22, 0xa5, 0x85, 0x1e, 0xc1, 0xbe, 0xa0, 0x13, 0x0c, 0x25, 0x5b, 0x1a, 0xfc, 0x7f, 0xa3,
	0xbd, 0xcc, 0x70, 0xce, 0x96, 0x48, 0xee, 0x02, 0x68, 0xa7, 0xe2, 0x53, 0x4c, 0xf2, 0x0e, 0x35,
	0xfc, 0x43, 0x66, 0xf0, 0x7e, 0x38, 0xd0, 0x29, 0xe7, 0x94, 0x82, 0x27, 0x12, 0xc9, 0x29, 0xec,
	0xd9, 0x41, 0x76, 0x9d, 0xde, 0x76, 0xbf, 0x39, 0xec, 0xfb, 0x2b, 0xa6, 0xed, 0x5b, 0x0a, 0x45,
	0x24, 0x79, 0x00, 0x07, 0x09, 0x7e, 0x53, 0x61, 0xa5, 0x85, 0x76, 0x66, 0x3e, 0x2b, 0xda, 0xf8,
	0xe9, 0x40, 0xe7, 0x79, 0x8a, 0x54, 0xe1, 0xad, 0x31, 0x1c, 0xc2, 0x8e, 0xa0, 0x29, 0x26, 0x2a,
	0x1f, 0x44, 0x7e, 0xd3, 0xb4, 0x0c, 0x32, 0x64, 0x51, 0x41, 0xcb, 0x58, 0x5e, 0x47, 0xe4, 0x04,
	0x76, 0xf3, 0x8b, 0x9e, 0xf8, 0x26, 0xcd, 0xdb, 0x40, 0xef, 0x97, 0x03, 0x9d, 0x8f, 0x22, 0xaa,
	0xf6, 0x74, 0x23, 0xb9, 0xf3, 0x8f, 0xc9, 0xc9, 0x33, 0x68, 0xce, 0x75, 0x6e, 0xbd, 0x6e, 0x9a,
	0x40, 0x73, 0xe8, 0xfa, 0x66, 0x23, 0x7d, 0xbb, 0x91, 0xfe, 0xcb, 0x6c, 0x23, 0xdf, 0x52, 0x39,
	0x1d, 0x81, 0x81, 0x67, 0x67, 0xef, 0x11, 0x74, 0x4e, 0x31, 0xc6, 0x4a, 0x63, 0x35, 0x3b, 0x33,
	0xfc, 0xdd, 0x80, 0x03, 0x2b, 0xee, 0x39, 0xa6, 0x0b, 0x76, 0x81, 0x44, 0x40, 0xbb, 0x34, 0x6c,
	0x72, 0xbc, 0x92, 0x40, 0x9d, 0x38, 0xee, 0xda, 0xbc, 0xc9, 0x15, 0xc0, 0xdf, 0x15, 0x27, 0xc3,
	0x95, 0x71, 0x95, 0xf7, 0xb0, 0x41, 0x2d, 0x01, 0xed, 0x92, 0x6c, 0x6b, 0xb0, 0xab, 0x93, 0x79,
	0x83, 0x8a, 0x9f, 0xa1, 0x5d, 0xd2, 0x63, 0x8d, 0x8a, 0x75, 0xfa, 0xb9, 0x87, 0x15, 0xfd, 0x5f,
	0x64, 0xdf, 0x15, 0xf9, 0x0e, 0xad, 0x9b, 0x6f, 0x94, 0x3c, 0x5d, 0x99, 0xbe, 0xe6, 0x9b, 0x70,
	0x8f, 0x37, 0x8c, 0x32, 0x1f, 0xc1, 0x49, 0xe3, 0xd3, 0xd6, 0x62, 0x30, 0xde, 0xd1, 0x2d, 0x3d,
	0xf9, 0x13, 0x00, 0x00, 0xff, 0xff, 0x1a, 0x38, 0x8f, 0xe3, 0xa1, 0x05, 0x00, 0x00,
}
