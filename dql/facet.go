/**
 * @Author: daipengyuan
 * @Description: 面/属性操作，Facet为属性插入的基础类型
 * @File:  facets
 * @Version: 1.0.0
 * @Date: 2021/8/17 15:03
 */

package dql

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"reflect"
	"time"
)

type ToFacet interface {
	Facet() (*Facet, error)
}

type Facet struct {
	Seq          int         `json:"seq"`
	PredWithLang string      `json:"pred"`
	Key          string      `json:"key"`
	Value        interface{} `json:"value"`
}

// Combine 将facet绑定到nquad
// 数据类型检查,谓词检查，序号检查
func (f *Facet) Combine(nq *api.NQuad) error {
	facet, err := f.parse()
	if err != nil {
		return err
	}
	nq.Facets = append(nq.Facets, facet)
	return nil
}

func (f *Facet) parse() (*api.Facet, error) {
	if f.Key == "" {
		return nil, errors.New("nil facet key")
	}
	var facet = &api.Facet{Key: f.Key}
	switch f.Value.(type) {
	case int:
		facet.ValType = api.Facet_INT
		facet.Value = []byte(fmt.Sprintf("%d", f.Value.(int)))
	case float64:
		facet.ValType = api.Facet_FLOAT
		facet.Value = []byte(fmt.Sprintf("%f", f.Value.(float64)))
	case string:
		facet.ValType = api.Facet_STRING
		facet.Value = []byte(f.Value.(string))
	case bool:
		facet.ValType = api.Facet_BOOL
		facet.Value = []byte(fmt.Sprintf("%t", f.Value.(bool)))
	case time.Time:
		facet.ValType = api.Facet_DATETIME
		b, err := f.Value.(time.Time).MarshalBinary()
		if err != nil {
			return nil, err
		}
		facet.Value = b
	default:
		return nil, errors.New(fmt.Sprintf("wrong facet datatype, need one of string/bool/int/float/datetime, get [%s]", reflect.TypeOf(f.Value)))
	}
	return facet, nil
}
