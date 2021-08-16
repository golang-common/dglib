/**
 * @Author: daipengyuan
 * @Description:
 * @File:  mutation_test
 * @Version: 1.0.0
 * @Date: 2021/8/15 16:05
 */

package dql

import (
	"github.com/twpayne/go-geom"
	"reflect"
	"testing"
)

type TA struct {
	A string `json:"a"`
	B string `json:"b"`
	TB
}

type TB struct {
	C string   `json:"c"`
	D []string `json:"d"`
}

var A = TA{
	A:  "mya",
	B:  "myb",
	TB: TB{C: "myc", D: []string{"myd1", "myd2"}},
}

func TestStruct(t *testing.T) {
	tp := reflect.TypeOf(A)
	val := reflect.ValueOf(A)
	fi := val.NumField()
	for i := 0; i < fi; i++ {
		t.Log(tp.Field(i).Tag.Get("json"))
		t.Log(tp.Field(i).Name)
	}

}

func TestField(t *testing.T) {
	var a = geom.NewMultiPoint(geom.Layout(2))
	t.Log(reflect.TypeOf(a))
	t.Log(reflect.ValueOf(a).Kind())
}
