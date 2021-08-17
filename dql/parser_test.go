/**
 * @Author: daipengyuan
 * @Description:
 * @File:  parser_test
 * @Version: 1.0.0
 * @Date: 2021/8/17 15:46
 */

package dql

import (
	"github.com/dgraph-io/dgo/v200/protos/api"
	"testing"
	"time"
)

//
//type NQuad struct {
//	Subject     string   `protobuf:"bytes,1,opt,name=subject,proto3" json:"subject,omitempty"`
//	Predicate   string   `protobuf:"bytes,2,opt,name=predicate,proto3" json:"predicate,omitempty"`
//	ObjectId    string   `protobuf:"bytes,3,opt,name=object_id,json=objectId,proto3" json:"object_id,omitempty"`
//	ObjectValue *Value   `protobuf:"bytes,4,opt,name=object_value,json=objectValue,proto3" json:"object_value,omitempty"`
//	Lang        string   `protobuf:"bytes,6,opt,name=lang,proto3" json:"lang,omitempty"`
//	Facets      []*Facet `protobuf:"bytes,7,rep,name=facets,proto3" json:"facets,omitempty"`
//	Namespace   uint64   `protobuf:"varint,8,opt,name=namespace,proto3" json:"namespace,omitempty"`
//}

func TestPar(t *testing.T) {
	c, err := NewClient(DgConfig)
	if err != nil {
		t.Fatal(err)
	}
	txn := c.Txn()
	b, _ := time.Now().MarshalBinary()
	resp, err := txn.Txn.Mutate(c.Ctx(), &api.Mutation{
		Set: []*api.NQuad{
			{
				Subject:     "_:a",
				Predicate:   "name",
				ObjectValue: &api.Value{Val: &api.Value_StrVal{StrVal: "呆"}},
				Lang:        "cn",
				Facets: []*api.Facet{
					{Key: "kick", Value: b, ValType: api.Facet_DATETIME, Alias: "ff"},
				},
			},
		},
	})
	txn.CommitOrAbort(err)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(IndentJson(resp))
}
