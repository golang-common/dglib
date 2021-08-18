/**
 * @Author: daipengyuan
 * @Description:
 * @File:  mutation_test
 * @Version: 1.0.0
 * @Date: 2021/8/18 09:16
 */

package dql

import (
	"fmt"
	"testing"
)

type Person struct {
	Uid      string   `json:"uid" db:"uid,string" dtype:"Person"`
	Name     string   `json:"name" db:"name,string,id" index:"index" token:"exact"`
	Age      int      `json:"age" db:"age,int" index:"index" token:"int"`
	Friend   []string `json:"friend" db:"friend,uid" index:"reverse,count,list"`
	FriendOf []string `json:"friend_of" db:"~friend,uid"`
}

func TestTxn_Add(t *testing.T) {
	my := Person{
		Name: "dpy1",
		Age:  33,
	}
	c, err := NewClient(DgConfig)
	if err != nil {
		t.Fatal(err)
	}
	txn := c.Txn()
	resp, err := txn.Add(my)
	defer txn.CommitOrAbort(err)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(IndentJson(resp))
}

func TestTxn_Update(t *testing.T) {
	my := Person{
		Uid:  "0xa",
		Name: "dpy1",
		Age:  3,
	}
	c, err := NewClient(DgConfig)
	if err != nil {
		t.Fatal(err)
	}
	txn := c.Txn()
	resp, err := txn.Update(my)
	defer txn.CommitOrAbort(err)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(IndentJson(resp))
}

func TestTxn_Merge(t *testing.T) {
	my := Person{
		Uid:    "0xc",
		Name:   "dpy1",
		Age:    99,
		Friend: []string{"0xd"},
	}
	c, err := NewClient(DgConfig)
	if err != nil {
		t.Fatal(err)
	}
	txn := c.Txn()
	resp, err := txn.Merge(my)
	defer txn.CommitOrAbort(err)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(IndentJson(resp))
}

func TestTxn_Delete(t *testing.T) {
	my := Person{
		Uid: "0xa",
		Age: 3,
	}
	c, err := NewClient(DgConfig)
	if err != nil {
		t.Fatal(err)
	}
	txn := c.Txn()
	resp, err := txn.Delete(my)
	defer txn.CommitOrAbort(err)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(IndentJson(resp))
}

func TestTxn_DelNode(t *testing.T) {
	my := Person{
		Uid: "0xa",
		Age: 3,
	}
	c, err := NewClient(DgConfig)
	if err != nil {
		t.Fatal(err)
	}
	txn := c.Txn()
	resp, err := txn.DelNode(my)
	defer txn.CommitOrAbort(err)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(IndentJson(resp))
}

func TestTemp(t *testing.T) {
	fmt.Printf("%x", 2333)
}
