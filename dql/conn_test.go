/**
 * @Author: daipengyuan
 * @Description:
 * @File:  conn_test
 * @Version: 1.0.0
 * @Date: 2021/8/11 16:42
 */

package dql

import (
	"fmt"
	"testing"
)

func TestConn1(t *testing.T) {
	dql, err := New("127.0.0.1:9080")
	if err != nil {
		t.Fatal(err)
	}
	pl, tl, err := dql.GetSchemaAll()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(IndentJson(pl))
	fmt.Println(IndentJson(tl))
}

func TestConn_SetPred(t *testing.T) {
	dql, err := New("127.0.0.1:9080")
	if err != nil {
		t.Fatal(err)
	}
	err = dql.SetPred(Pred{
		Predicate: "testset1",
		Type:      TypeString,
		Index:     true,
		Tokenizer: []string{IndexHash,IndexTrigram},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestConn_DropPred(t *testing.T) {
	dql, err := New("127.0.0.1:9080")
	if err != nil {
		t.Fatal(err)
	}
	err = dql.DropPred("testset1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestConn_SetType(t *testing.T) {
	dql, err := New("127.0.0.1:9080")
	if err != nil {
		t.Fatal(err)
	}
	err = dql.SetType(Type{Name: "testset2", Fields: []string{"h33", "dd"}})
	if err != nil {
		t.Fatal(err)
	}
}

func TestConn_DropType(t *testing.T) {
	dql, err := New("127.0.0.1:9080")
	if err != nil {
		t.Fatal(err)
	}
	err = dql.DropType("testset2")
	if err != nil {
		t.Fatal(err)
	}
}
