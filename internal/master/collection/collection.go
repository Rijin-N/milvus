package collection

import (
	"time"

	"github.com/golang/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
	"github.com/zilliztech/milvus-distributed/internal/proto/etcdpb"
	"github.com/zilliztech/milvus-distributed/internal/proto/schemapb"
	"github.com/zilliztech/milvus-distributed/internal/util/typeutil"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type UniqueID = typeutil.UniqueID
type Timestamp = typeutil.Timestamp

type Collection struct {
	ID         UniqueID    `json:"id"`
	Name       string      `json:"name"`
	CreateTime Timestamp   `json:"creat_time"`
	Schema     []FieldMeta `json:"schema"`
	//	ExtraSchema       []FieldMeta `json:"extra_schema"`
	SegmentIDs        []UniqueID `json:"segment_ids"`
	PartitionTags     []string   `json:"partition_tags"`
	GrpcMarshalString string     `json:"grpc_marshal_string"`
}

type FieldMeta struct {
	FieldName string            `json:"field_name"`
	Type      schemapb.DataType `json:"type"`
	DIM       int64             `json:"dimension"`
}

func GrpcMarshal(c *Collection) *Collection {
	if c.GrpcMarshalString != "" {
		c.GrpcMarshalString = ""
	}
	pbSchema := &schemapb.CollectionSchema{
		Fields: []*schemapb.FieldSchema{},
	}
	schemaSlice := []*schemapb.FieldSchema{}
	for _, v := range c.Schema {
		newpbMeta := &schemapb.FieldSchema{
			Name:     v.FieldName,
			DataType: schemapb.DataType(v.Type), //czs_tag
		}
		schemaSlice = append(schemaSlice, newpbMeta)
	}
	pbSchema.Fields = schemaSlice
	grpcCollection := &etcdpb.CollectionMeta{
		ID:            c.ID,
		Schema:        pbSchema,
		CreateTime:    c.CreateTime,
		SegmentIDs:    c.SegmentIDs,
		PartitionTags: c.PartitionTags,
	}
	out := proto.MarshalTextString(grpcCollection)
	c.GrpcMarshalString = out
	return c
}

func NewCollection(id UniqueID, name string, createTime time.Time,
	schema []*schemapb.FieldSchema, sIDs []UniqueID, ptags []string) Collection {

	segementIDs := []UniqueID{}
	newSchema := []FieldMeta{}
	for _, v := range schema {
		newSchema = append(newSchema, FieldMeta{FieldName: v.Name, Type: v.DataType, DIM: 16})
	}
	segementIDs = append(segementIDs, sIDs...)
	return Collection{
		ID:            id,
		Name:          name,
		CreateTime:    Timestamp(createTime.Unix()),
		Schema:        newSchema,
		SegmentIDs:    segementIDs,
		PartitionTags: ptags,
	}
}

func Collection2JSON(c Collection) (string, error) {
	b, err := json.Marshal(&c)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func JSON2Collection(s string) (*Collection, error) {
	var c Collection
	err := json.Unmarshal([]byte(s), &c)
	if err != nil {
		return &Collection{}, err
	}
	return &c, nil
}
