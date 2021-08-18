/**
 * @Author: daipengyuan
 * @Description:
 * @File:  conn_test
 * @Version: 1.0.0
 * @Date: 2021/8/13 17:05
 */

package dql

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

var DgConfig = Config{
	Targets:     []string{"localhost:9080"},
	DialTimeout: 2 * time.Second,
	//OptTimeout:  2 * time.Second,
	Username: "groot",
	Password: "password",
	Tls: Tls{
		ServeName:  "crane",
		CaCert:     "/Users/lyonsdpy/Data/dgraph/tls/ca.crt",
		ClientCert: "/Users/lyonsdpy/Data/dgraph/tls/client.crane.crt",
		ClientKey:  "/Users/lyonsdpy/Data/dgraph/tls/client.crane.key",
	},
}

func TestGetSchema(t *testing.T) {
	c, err := NewClient(DgConfig)
	if err != nil {
		t.Fatal(err)
	}
	txn := c.Txn()
	schema, err := txn.GetSchema()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", schema))
}

var Pred1 = Pred{
	Predicate: "name",
	Type:      "string",
	Index:     true,
	Tokenizer: []string{TokenTerm},
	Upsert:    true,
}

var Pred2 = Pred{
	Predicate: "age",
	Type:      "int",
	Index:     false,
}

var Pred3 = Pred{
	Predicate: "friend",
	Type:      "uid",
	Reverse:   true,
	Count:     true,
	List:      true,
}

func TestSetPred(t *testing.T) {
	c, err := NewClient(DgConfig)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetPred(Pred1)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetPred(Pred2)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetPred(Pred3)
}

var Type1 = Type{
	Name: "Person",
	Fields: []Field{
		{Name: "name"},
		{Name: "age"},
		{Name: "friend"},
		{Name: "~friend"},
	},
}

func TestSetType(t *testing.T) {
	c, err := NewClient(DgConfig)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetType(Type1)
	if err != nil {
		t.Fatal(err)
	}
}

func Test2(t *testing.T) {
	var a = "dddd"
	t.Log(strings.Split(a, "@"))
}

//type Person struct {
//	Uid      string   `json:"uid" db:"uid,string" dtype:"Person"`
//	Name     string   `json:"name" db:"name,string" index:"index,exact,id"`
//	Age      int      `json:"age" db:"age,int" index:"index,int"`
//	Friend   []Person `json:"friend" db:"friend,uid" index:"reverse,count,list"`
//	FriendOf []Person `json:"friend_of" db:"~friend,uid"`
//}
//
//func TestSetVal1(t *testing.T) {
//	var a = Person{
//		Name: "DP2",
//		Age:  22,
//	}
//	c, err := NewClient(DgConfig)
//	if err != nil {
//		t.Fatal(err)
//	}
//	txn := c.Txn()
//	resp, err := txn.AddNode(a)
//	defer txn.CommitOrAbort(err)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log(IndentJson(resp))
//}

func TestSetType2(t *testing.T) {
	c, err := NewClient(DgConfig)
	if err != nil {
		t.Fatal(err)
	}
	typ, preds, err := UnmarshalSchema(Person{})
	if err != nil {
		t.Fatal(err)
	}
	for _, pred := range preds {
		err = c.SetPred(pred)
		if err != nil {
			t.Fatal(err)
		}
	}
	err = c.SetType(*typ)
	if err != nil {
		t.Fatal(err)
	}
}
