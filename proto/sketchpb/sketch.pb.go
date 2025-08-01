// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.21.12
// source: sketch.proto

package sketchpb

import (
	geometrypb "github.com/hansmi/dossier/proto/geometrypb"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Node features are attachment points on a node. Relative positions on other
// nodes refer to these.
type NodeFeature int32

const (
	NodeFeature_NODE_FEATURE_UNSPECIFIED NodeFeature = 0
	NodeFeature_TOP_LEFT                 NodeFeature = 1
	NodeFeature_TOP_RIGHT                NodeFeature = 2
	NodeFeature_BOTTOM_LEFT              NodeFeature = 3
	NodeFeature_BOTTOM_RIGHT             NodeFeature = 4
)

// Enum value maps for NodeFeature.
var (
	NodeFeature_name = map[int32]string{
		0: "NODE_FEATURE_UNSPECIFIED",
		1: "TOP_LEFT",
		2: "TOP_RIGHT",
		3: "BOTTOM_LEFT",
		4: "BOTTOM_RIGHT",
	}
	NodeFeature_value = map[string]int32{
		"NODE_FEATURE_UNSPECIFIED": 0,
		"TOP_LEFT":                 1,
		"TOP_RIGHT":                2,
		"BOTTOM_LEFT":              3,
		"BOTTOM_RIGHT":             4,
	}
)

func (x NodeFeature) Enum() *NodeFeature {
	p := new(NodeFeature)
	*p = x
	return p
}

func (x NodeFeature) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NodeFeature) Descriptor() protoreflect.EnumDescriptor {
	return file_sketch_proto_enumTypes[0].Descriptor()
}

func (NodeFeature) Type() protoreflect.EnumType {
	return &file_sketch_proto_enumTypes[0]
}

func (x NodeFeature) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NodeFeature.Descriptor instead.
func (NodeFeature) EnumDescriptor() ([]byte, []int) {
	return file_sketch_proto_rawDescGZIP(), []int{0}
}

// A one-dimensional position relative to a feature on another node.
type RelativePosition1D struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Referenced node identifier and feature.
	Node    string      `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	Feature NodeFeature `protobuf:"varint,2,opt,name=feature,proto3,enum=dossier.sketch.NodeFeature" json:"feature,omitempty"`
	// Shift relative position by the given distance.
	Offset        *geometrypb.Length `protobuf:"bytes,3,opt,name=offset,proto3" json:"offset,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RelativePosition1D) Reset() {
	*x = RelativePosition1D{}
	mi := &file_sketch_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RelativePosition1D) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RelativePosition1D) ProtoMessage() {}

func (x *RelativePosition1D) ProtoReflect() protoreflect.Message {
	mi := &file_sketch_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RelativePosition1D.ProtoReflect.Descriptor instead.
func (*RelativePosition1D) Descriptor() ([]byte, []int) {
	return file_sketch_proto_rawDescGZIP(), []int{0}
}

func (x *RelativePosition1D) GetNode() string {
	if x != nil {
		return x.Node
	}
	return ""
}

func (x *RelativePosition1D) GetFeature() NodeFeature {
	if x != nil {
		return x.Feature
	}
	return NodeFeature_NODE_FEATURE_UNSPECIFIED
}

func (x *RelativePosition1D) GetOffset() *geometrypb.Length {
	if x != nil {
		return x.Offset
	}
	return nil
}

// A two-dimensional position relative to a feature on another node.
type RelativePosition2D struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Referenced node identifier and feature.
	Node    string      `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	Feature NodeFeature `protobuf:"varint,2,opt,name=feature,proto3,enum=dossier.sketch.NodeFeature" json:"feature,omitempty"`
	// Shift relative position by the given distances.
	Offset        *geometrypb.Size `protobuf:"bytes,3,opt,name=offset,proto3" json:"offset,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RelativePosition2D) Reset() {
	*x = RelativePosition2D{}
	mi := &file_sketch_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RelativePosition2D) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RelativePosition2D) ProtoMessage() {}

func (x *RelativePosition2D) ProtoReflect() protoreflect.Message {
	mi := &file_sketch_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RelativePosition2D.ProtoReflect.Descriptor instead.
func (*RelativePosition2D) Descriptor() ([]byte, []int) {
	return file_sketch_proto_rawDescGZIP(), []int{1}
}

func (x *RelativePosition2D) GetNode() string {
	if x != nil {
		return x.Node
	}
	return ""
}

func (x *RelativePosition2D) GetFeature() NodeFeature {
	if x != nil {
		return x.Feature
	}
	return NodeFeature_NODE_FEATURE_UNSPECIFIED
}

func (x *RelativePosition2D) GetOffset() *geometrypb.Size {
	if x != nil {
		return x.Offset
	}
	return nil
}

// FlexRect describes an abstract rectangle. The four edges (lines) can be
// specified as absolute or relative positions or via vertices (corners) and/or
// the rectangle size. Each edge may only be specified through one method.
// Positions are relative to the top left corner of a page.
//
//	Top left             Top right
//	vertex                  vertex
//	┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓---
//	┃          Top edge          ┃  ^
//	┃                            ┃  ┆
//	┃                            ┃  ┆
//	┃ Left                 Right ┃  ┆
//	┃ edge                  edge ┃  ┆ Height
//	┃                            ┃  ┆
//	┃                            ┃  ┆
//	┃        Bottom edge         ┃  v
//	┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛---
//	| Bottom left   Bottom right |
//	| vertex              vertex |
//	|                            |
//	|<┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄>|
//	            Width
//
// TODO: Allow specifying a different page corner for absolute positions, e.g.
// bottom left.
type FlexRect struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TopLeft       *FlexRect_Vertex       `protobuf:"bytes,1,opt,name=top_left,json=topLeft,proto3" json:"top_left,omitempty"`
	TopRight      *FlexRect_Vertex       `protobuf:"bytes,2,opt,name=top_right,json=topRight,proto3" json:"top_right,omitempty"`
	BottomLeft    *FlexRect_Vertex       `protobuf:"bytes,3,opt,name=bottom_left,json=bottomLeft,proto3" json:"bottom_left,omitempty"`
	BottomRight   *FlexRect_Vertex       `protobuf:"bytes,4,opt,name=bottom_right,json=bottomRight,proto3" json:"bottom_right,omitempty"`
	Top           *FlexRect_Edge         `protobuf:"bytes,5,opt,name=top,proto3" json:"top,omitempty"`
	Right         *FlexRect_Edge         `protobuf:"bytes,6,opt,name=right,proto3" json:"right,omitempty"`
	Bottom        *FlexRect_Edge         `protobuf:"bytes,7,opt,name=bottom,proto3" json:"bottom,omitempty"`
	Left          *FlexRect_Edge         `protobuf:"bytes,8,opt,name=left,proto3" json:"left,omitempty"`
	Width         *geometrypb.Length     `protobuf:"bytes,9,opt,name=width,proto3" json:"width,omitempty"`
	Height        *geometrypb.Length     `protobuf:"bytes,10,opt,name=height,proto3" json:"height,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FlexRect) Reset() {
	*x = FlexRect{}
	mi := &file_sketch_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FlexRect) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlexRect) ProtoMessage() {}

func (x *FlexRect) ProtoReflect() protoreflect.Message {
	mi := &file_sketch_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlexRect.ProtoReflect.Descriptor instead.
func (*FlexRect) Descriptor() ([]byte, []int) {
	return file_sketch_proto_rawDescGZIP(), []int{2}
}

func (x *FlexRect) GetTopLeft() *FlexRect_Vertex {
	if x != nil {
		return x.TopLeft
	}
	return nil
}

func (x *FlexRect) GetTopRight() *FlexRect_Vertex {
	if x != nil {
		return x.TopRight
	}
	return nil
}

func (x *FlexRect) GetBottomLeft() *FlexRect_Vertex {
	if x != nil {
		return x.BottomLeft
	}
	return nil
}

func (x *FlexRect) GetBottomRight() *FlexRect_Vertex {
	if x != nil {
		return x.BottomRight
	}
	return nil
}

func (x *FlexRect) GetTop() *FlexRect_Edge {
	if x != nil {
		return x.Top
	}
	return nil
}

func (x *FlexRect) GetRight() *FlexRect_Edge {
	if x != nil {
		return x.Right
	}
	return nil
}

func (x *FlexRect) GetBottom() *FlexRect_Edge {
	if x != nil {
		return x.Bottom
	}
	return nil
}

func (x *FlexRect) GetLeft() *FlexRect_Edge {
	if x != nil {
		return x.Left
	}
	return nil
}

func (x *FlexRect) GetWidth() *geometrypb.Length {
	if x != nil {
		return x.Width
	}
	return nil
}

func (x *FlexRect) GetHeight() *geometrypb.Length {
	if x != nil {
		return x.Height
	}
	return nil
}

type Node struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unique node identifier.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Rectangles in which content should be matched. Multiple may be specified.
	SearchAreas []*FlexRect `protobuf:"bytes,100,rep,name=search_areas,json=searchAreas,proto3" json:"search_areas,omitempty"`
	// Types that are valid to be assigned to Matcher:
	//
	//	*Node_BlockText
	//	*Node_LineText
	Matcher isNode_Matcher `protobuf_oneof:"matcher"`
	// Tags are arbitrary non-empty, unique strings.
	Tags          []string `protobuf:"bytes,15,rep,name=tags,proto3" json:"tags,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Node) Reset() {
	*x = Node{}
	mi := &file_sketch_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_sketch_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_sketch_proto_rawDescGZIP(), []int{3}
}

func (x *Node) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Node) GetSearchAreas() []*FlexRect {
	if x != nil {
		return x.SearchAreas
	}
	return nil
}

func (x *Node) GetMatcher() isNode_Matcher {
	if x != nil {
		return x.Matcher
	}
	return nil
}

func (x *Node) GetBlockText() *Node_TextMatch {
	if x != nil {
		if x, ok := x.Matcher.(*Node_BlockText); ok {
			return x.BlockText
		}
	}
	return nil
}

func (x *Node) GetLineText() *Node_TextMatch {
	if x != nil {
		if x, ok := x.Matcher.(*Node_LineText); ok {
			return x.LineText
		}
	}
	return nil
}

func (x *Node) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

type isNode_Matcher interface {
	isNode_Matcher()
}

type Node_BlockText struct {
	// Match over blocks of text. A block contains one or more lines.
	BlockText *Node_TextMatch `protobuf:"bytes,10,opt,name=block_text,json=blockText,proto3,oneof"`
}

type Node_LineText struct {
	// Match over single lines of text.
	LineText *Node_TextMatch `protobuf:"bytes,11,opt,name=line_text,json=lineText,proto3,oneof"`
}

func (*Node_BlockText) isNode_Matcher() {}

func (*Node_LineText) isNode_Matcher() {}

// A sketch is an abstract description of where information on a page is to be
// found. Sketches have no concept of multiple pages. If code needs to make
// a distinction between pages the following approaches may be useful:
//
//  1. Define nodes for all pages in a single sketch. Look up nodes in the
//     analysis results depending on the current page number and ignore all
//     other nodes.
//
// 2) Define one sketch per page type and apply them as necessary.
type Sketch struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Nodes []*Node                `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	// Tags are arbitrary non-empty, unique strings.
	Tags          []string `protobuf:"bytes,15,rep,name=tags,proto3" json:"tags,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Sketch) Reset() {
	*x = Sketch{}
	mi := &file_sketch_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Sketch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sketch) ProtoMessage() {}

func (x *Sketch) ProtoReflect() protoreflect.Message {
	mi := &file_sketch_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sketch.ProtoReflect.Descriptor instead.
func (*Sketch) Descriptor() ([]byte, []int) {
	return file_sketch_proto_rawDescGZIP(), []int{4}
}

func (x *Sketch) GetNodes() []*Node {
	if x != nil {
		return x.Nodes
	}
	return nil
}

func (x *Sketch) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

type FlexRect_Vertex struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Method:
	//
	//	*FlexRect_Vertex_Abs
	//	*FlexRect_Vertex_Rel
	Method        isFlexRect_Vertex_Method `protobuf_oneof:"method"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FlexRect_Vertex) Reset() {
	*x = FlexRect_Vertex{}
	mi := &file_sketch_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FlexRect_Vertex) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlexRect_Vertex) ProtoMessage() {}

func (x *FlexRect_Vertex) ProtoReflect() protoreflect.Message {
	mi := &file_sketch_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlexRect_Vertex.ProtoReflect.Descriptor instead.
func (*FlexRect_Vertex) Descriptor() ([]byte, []int) {
	return file_sketch_proto_rawDescGZIP(), []int{2, 0}
}

func (x *FlexRect_Vertex) GetMethod() isFlexRect_Vertex_Method {
	if x != nil {
		return x.Method
	}
	return nil
}

func (x *FlexRect_Vertex) GetAbs() *geometrypb.Point {
	if x != nil {
		if x, ok := x.Method.(*FlexRect_Vertex_Abs); ok {
			return x.Abs
		}
	}
	return nil
}

func (x *FlexRect_Vertex) GetRel() *RelativePosition2D {
	if x != nil {
		if x, ok := x.Method.(*FlexRect_Vertex_Rel); ok {
			return x.Rel
		}
	}
	return nil
}

type isFlexRect_Vertex_Method interface {
	isFlexRect_Vertex_Method()
}

type FlexRect_Vertex_Abs struct {
	Abs *geometrypb.Point `protobuf:"bytes,1,opt,name=abs,proto3,oneof"`
}

type FlexRect_Vertex_Rel struct {
	Rel *RelativePosition2D `protobuf:"bytes,2,opt,name=rel,proto3,oneof"`
}

func (*FlexRect_Vertex_Abs) isFlexRect_Vertex_Method() {}

func (*FlexRect_Vertex_Rel) isFlexRect_Vertex_Method() {}

type FlexRect_Edge struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Method:
	//
	//	*FlexRect_Edge_Abs
	//	*FlexRect_Edge_Rel
	Method        isFlexRect_Edge_Method `protobuf_oneof:"method"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FlexRect_Edge) Reset() {
	*x = FlexRect_Edge{}
	mi := &file_sketch_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FlexRect_Edge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlexRect_Edge) ProtoMessage() {}

func (x *FlexRect_Edge) ProtoReflect() protoreflect.Message {
	mi := &file_sketch_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlexRect_Edge.ProtoReflect.Descriptor instead.
func (*FlexRect_Edge) Descriptor() ([]byte, []int) {
	return file_sketch_proto_rawDescGZIP(), []int{2, 1}
}

func (x *FlexRect_Edge) GetMethod() isFlexRect_Edge_Method {
	if x != nil {
		return x.Method
	}
	return nil
}

func (x *FlexRect_Edge) GetAbs() *geometrypb.Length {
	if x != nil {
		if x, ok := x.Method.(*FlexRect_Edge_Abs); ok {
			return x.Abs
		}
	}
	return nil
}

func (x *FlexRect_Edge) GetRel() *RelativePosition1D {
	if x != nil {
		if x, ok := x.Method.(*FlexRect_Edge_Rel); ok {
			return x.Rel
		}
	}
	return nil
}

type isFlexRect_Edge_Method interface {
	isFlexRect_Edge_Method()
}

type FlexRect_Edge_Abs struct {
	Abs *geometrypb.Length `protobuf:"bytes,1,opt,name=abs,proto3,oneof"`
}

type FlexRect_Edge_Rel struct {
	Rel *RelativePosition1D `protobuf:"bytes,2,opt,name=rel,proto3,oneof"`
}

func (*FlexRect_Edge_Abs) isFlexRect_Edge_Method() {}

func (*FlexRect_Edge_Rel) isFlexRect_Edge_Method() {}

type Node_TextMatch struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Regular expression to look for. If a source document has been processed
	// via OCR the expression may have to be written in a more flexible form to
	// accept lookalike-characters (e.g. German umlauts may be recognized as A,
	// O and U without diacritic).
	//
	// Example: "(?i)^\\s*Destination\\s+address\\s*:"
	//
	// Syntax: https://pkg.go.dev/regexp/syntax
	Regex string `protobuf:"bytes,1,opt,name=regex,proto3" json:"regex,omitempty"`
	// By default the reported bounds of valid matches are those of the block
	// or line containing the matched text. By setting "bounds_from_match" the
	// exact bounds of the matched text are used instead.
	BoundsFromMatch bool `protobuf:"varint,2,opt,name=bounds_from_match,json=boundsFromMatch,proto3" json:"bounds_from_match,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *Node_TextMatch) Reset() {
	*x = Node_TextMatch{}
	mi := &file_sketch_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Node_TextMatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node_TextMatch) ProtoMessage() {}

func (x *Node_TextMatch) ProtoReflect() protoreflect.Message {
	mi := &file_sketch_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node_TextMatch.ProtoReflect.Descriptor instead.
func (*Node_TextMatch) Descriptor() ([]byte, []int) {
	return file_sketch_proto_rawDescGZIP(), []int{3, 0}
}

func (x *Node_TextMatch) GetRegex() string {
	if x != nil {
		return x.Regex
	}
	return ""
}

func (x *Node_TextMatch) GetBoundsFromMatch() bool {
	if x != nil {
		return x.BoundsFromMatch
	}
	return false
}

var File_sketch_proto protoreflect.FileDescriptor

const file_sketch_proto_rawDesc = "" +
	"\n" +
	"\fsketch.proto\x12\x0edossier.sketch\x1a\x0egeometry.proto\"\x91\x01\n" +
	"\x12RelativePosition1D\x12\x12\n" +
	"\x04node\x18\x01 \x01(\tR\x04node\x125\n" +
	"\afeature\x18\x02 \x01(\x0e2\x1b.dossier.sketch.NodeFeatureR\afeature\x120\n" +
	"\x06offset\x18\x03 \x01(\v2\x18.dossier.geometry.LengthR\x06offset\"\x8f\x01\n" +
	"\x12RelativePosition2D\x12\x12\n" +
	"\x04node\x18\x01 \x01(\tR\x04node\x125\n" +
	"\afeature\x18\x02 \x01(\x0e2\x1b.dossier.sketch.NodeFeatureR\afeature\x12.\n" +
	"\x06offset\x18\x03 \x01(\v2\x16.dossier.geometry.SizeR\x06offset\"\xad\x06\n" +
	"\bFlexRect\x12:\n" +
	"\btop_left\x18\x01 \x01(\v2\x1f.dossier.sketch.FlexRect.VertexR\atopLeft\x12<\n" +
	"\ttop_right\x18\x02 \x01(\v2\x1f.dossier.sketch.FlexRect.VertexR\btopRight\x12@\n" +
	"\vbottom_left\x18\x03 \x01(\v2\x1f.dossier.sketch.FlexRect.VertexR\n" +
	"bottomLeft\x12B\n" +
	"\fbottom_right\x18\x04 \x01(\v2\x1f.dossier.sketch.FlexRect.VertexR\vbottomRight\x12/\n" +
	"\x03top\x18\x05 \x01(\v2\x1d.dossier.sketch.FlexRect.EdgeR\x03top\x123\n" +
	"\x05right\x18\x06 \x01(\v2\x1d.dossier.sketch.FlexRect.EdgeR\x05right\x125\n" +
	"\x06bottom\x18\a \x01(\v2\x1d.dossier.sketch.FlexRect.EdgeR\x06bottom\x121\n" +
	"\x04left\x18\b \x01(\v2\x1d.dossier.sketch.FlexRect.EdgeR\x04left\x12.\n" +
	"\x05width\x18\t \x01(\v2\x18.dossier.geometry.LengthR\x05width\x120\n" +
	"\x06height\x18\n" +
	" \x01(\v2\x18.dossier.geometry.LengthR\x06height\x1aw\n" +
	"\x06Vertex\x12+\n" +
	"\x03abs\x18\x01 \x01(\v2\x17.dossier.geometry.PointH\x00R\x03abs\x126\n" +
	"\x03rel\x18\x02 \x01(\v2\".dossier.sketch.RelativePosition2DH\x00R\x03relB\b\n" +
	"\x06method\x1av\n" +
	"\x04Edge\x12,\n" +
	"\x03abs\x18\x01 \x01(\v2\x18.dossier.geometry.LengthH\x00R\x03abs\x126\n" +
	"\x03rel\x18\x02 \x01(\v2\".dossier.sketch.RelativePosition1DH\x00R\x03relB\b\n" +
	"\x06method\"\xc5\x02\n" +
	"\x04Node\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12;\n" +
	"\fsearch_areas\x18d \x03(\v2\x18.dossier.sketch.FlexRectR\vsearchAreas\x12?\n" +
	"\n" +
	"block_text\x18\n" +
	" \x01(\v2\x1e.dossier.sketch.Node.TextMatchH\x00R\tblockText\x12=\n" +
	"\tline_text\x18\v \x01(\v2\x1e.dossier.sketch.Node.TextMatchH\x00R\blineText\x12\x12\n" +
	"\x04tags\x18\x0f \x03(\tR\x04tags\x1aM\n" +
	"\tTextMatch\x12\x14\n" +
	"\x05regex\x18\x01 \x01(\tR\x05regex\x12*\n" +
	"\x11bounds_from_match\x18\x02 \x01(\bR\x0fboundsFromMatchB\t\n" +
	"\amatcher\"H\n" +
	"\x06Sketch\x12*\n" +
	"\x05nodes\x18\x01 \x03(\v2\x14.dossier.sketch.NodeR\x05nodes\x12\x12\n" +
	"\x04tags\x18\x0f \x03(\tR\x04tags*k\n" +
	"\vNodeFeature\x12\x1c\n" +
	"\x18NODE_FEATURE_UNSPECIFIED\x10\x00\x12\f\n" +
	"\bTOP_LEFT\x10\x01\x12\r\n" +
	"\tTOP_RIGHT\x10\x02\x12\x0f\n" +
	"\vBOTTOM_LEFT\x10\x03\x12\x10\n" +
	"\fBOTTOM_RIGHT\x10\x04B*Z(github.com/hansmi/dossier/proto/sketchpbb\x06proto3"

var (
	file_sketch_proto_rawDescOnce sync.Once
	file_sketch_proto_rawDescData []byte
)

func file_sketch_proto_rawDescGZIP() []byte {
	file_sketch_proto_rawDescOnce.Do(func() {
		file_sketch_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_sketch_proto_rawDesc), len(file_sketch_proto_rawDesc)))
	})
	return file_sketch_proto_rawDescData
}

var file_sketch_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_sketch_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_sketch_proto_goTypes = []any{
	(NodeFeature)(0),           // 0: dossier.sketch.NodeFeature
	(*RelativePosition1D)(nil), // 1: dossier.sketch.RelativePosition1D
	(*RelativePosition2D)(nil), // 2: dossier.sketch.RelativePosition2D
	(*FlexRect)(nil),           // 3: dossier.sketch.FlexRect
	(*Node)(nil),               // 4: dossier.sketch.Node
	(*Sketch)(nil),             // 5: dossier.sketch.Sketch
	(*FlexRect_Vertex)(nil),    // 6: dossier.sketch.FlexRect.Vertex
	(*FlexRect_Edge)(nil),      // 7: dossier.sketch.FlexRect.Edge
	(*Node_TextMatch)(nil),     // 8: dossier.sketch.Node.TextMatch
	(*geometrypb.Length)(nil),  // 9: dossier.geometry.Length
	(*geometrypb.Size)(nil),    // 10: dossier.geometry.Size
	(*geometrypb.Point)(nil),   // 11: dossier.geometry.Point
}
var file_sketch_proto_depIdxs = []int32{
	0,  // 0: dossier.sketch.RelativePosition1D.feature:type_name -> dossier.sketch.NodeFeature
	9,  // 1: dossier.sketch.RelativePosition1D.offset:type_name -> dossier.geometry.Length
	0,  // 2: dossier.sketch.RelativePosition2D.feature:type_name -> dossier.sketch.NodeFeature
	10, // 3: dossier.sketch.RelativePosition2D.offset:type_name -> dossier.geometry.Size
	6,  // 4: dossier.sketch.FlexRect.top_left:type_name -> dossier.sketch.FlexRect.Vertex
	6,  // 5: dossier.sketch.FlexRect.top_right:type_name -> dossier.sketch.FlexRect.Vertex
	6,  // 6: dossier.sketch.FlexRect.bottom_left:type_name -> dossier.sketch.FlexRect.Vertex
	6,  // 7: dossier.sketch.FlexRect.bottom_right:type_name -> dossier.sketch.FlexRect.Vertex
	7,  // 8: dossier.sketch.FlexRect.top:type_name -> dossier.sketch.FlexRect.Edge
	7,  // 9: dossier.sketch.FlexRect.right:type_name -> dossier.sketch.FlexRect.Edge
	7,  // 10: dossier.sketch.FlexRect.bottom:type_name -> dossier.sketch.FlexRect.Edge
	7,  // 11: dossier.sketch.FlexRect.left:type_name -> dossier.sketch.FlexRect.Edge
	9,  // 12: dossier.sketch.FlexRect.width:type_name -> dossier.geometry.Length
	9,  // 13: dossier.sketch.FlexRect.height:type_name -> dossier.geometry.Length
	3,  // 14: dossier.sketch.Node.search_areas:type_name -> dossier.sketch.FlexRect
	8,  // 15: dossier.sketch.Node.block_text:type_name -> dossier.sketch.Node.TextMatch
	8,  // 16: dossier.sketch.Node.line_text:type_name -> dossier.sketch.Node.TextMatch
	4,  // 17: dossier.sketch.Sketch.nodes:type_name -> dossier.sketch.Node
	11, // 18: dossier.sketch.FlexRect.Vertex.abs:type_name -> dossier.geometry.Point
	2,  // 19: dossier.sketch.FlexRect.Vertex.rel:type_name -> dossier.sketch.RelativePosition2D
	9,  // 20: dossier.sketch.FlexRect.Edge.abs:type_name -> dossier.geometry.Length
	1,  // 21: dossier.sketch.FlexRect.Edge.rel:type_name -> dossier.sketch.RelativePosition1D
	22, // [22:22] is the sub-list for method output_type
	22, // [22:22] is the sub-list for method input_type
	22, // [22:22] is the sub-list for extension type_name
	22, // [22:22] is the sub-list for extension extendee
	0,  // [0:22] is the sub-list for field type_name
}

func init() { file_sketch_proto_init() }
func file_sketch_proto_init() {
	if File_sketch_proto != nil {
		return
	}
	file_sketch_proto_msgTypes[3].OneofWrappers = []any{
		(*Node_BlockText)(nil),
		(*Node_LineText)(nil),
	}
	file_sketch_proto_msgTypes[5].OneofWrappers = []any{
		(*FlexRect_Vertex_Abs)(nil),
		(*FlexRect_Vertex_Rel)(nil),
	}
	file_sketch_proto_msgTypes[6].OneofWrappers = []any{
		(*FlexRect_Edge_Abs)(nil),
		(*FlexRect_Edge_Rel)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_sketch_proto_rawDesc), len(file_sketch_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_sketch_proto_goTypes,
		DependencyIndexes: file_sketch_proto_depIdxs,
		EnumInfos:         file_sketch_proto_enumTypes,
		MessageInfos:      file_sketch_proto_msgTypes,
	}.Build()
	File_sketch_proto = out.File
	file_sketch_proto_goTypes = nil
	file_sketch_proto_depIdxs = nil
}
