/**
 * @Author: daipengyuan
 * @Description: dgraph dql使用grpc连接
 * @File:  conn
 * @Version: 1.0.0
 * @Date: 2021/8/9 10:33
 */

package dql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
	"strings"
)

// NewClient 新建dgraph连接
// target格式为 192.168.1.100:9080
// 第一个返回值为dgraph操作对象
func NewClient(targets []string) (*Client, error) {
	var clients []api.DgraphClient
	for _, target := range targets {
		grpcConn, err := grpc.DialContext(context.Background(), target, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		client := api.NewDgraphClient(grpcConn)
		clients = append(clients, client)
	}

	dgraph := dgo.NewDgraphClient(clients...)
	return &Client{client: dgraph}, nil
}

type Schema struct {
	Preds []Pred `json:"schema"`
	Types []struct {
		Name   string `json:"name"`
		Fields []struct {
			Name string `json:"name"`
		} `json:"fields"`
	} `json:"types"`
}

func (s Schema) ListType() []Type {
	var r []Type
	for _, v := range s.Types {
		if strings.HasPrefix(v.Name, "dgraph.") {
			continue
		}
		var tp = Type{Name: v.Name}
		for _, field := range v.Fields {
			tp.Fields = append(tp.Fields, field.Name)
		}
		r = append(r, tp)
	}
	return r
}

func (s Schema) ListPred() []Pred {
	var r []Pred
	for _, p := range s.Preds {
		if strings.HasPrefix(p.Predicate, "dgraph.") {
			continue
		}
		r = append(r, p)
	}
	return r
}

type Client struct {
	client *dgo.Dgraph
}

func (d *Client) Txn(ReadOnly ...bool) *Txn {
	if len(ReadOnly) > 0 && ReadOnly[0] == true {
		return &Txn{Txn: d.client.NewReadOnlyTxn(), Ctx: context.Background(), Readonly: ReadOnly[0]}
	}
	return &Txn{Txn: d.client.NewTxn(), Ctx: context.Background()}
}

func (d *Client) SetPred(pred Pred) error {
	err := d.client.Alter(context.Background(), &api.Operation{
		Schema: pred.Rdf(),
	})
	return err
}

func (d *Client) DropPred(name string) error {
	err := d.client.Alter(context.Background(), &api.Operation{
		DropValue: name,
		DropOp:    api.Operation_ATTR,
	})
	return err
}

func (d *Client) SetType(tp Type) error {
	err := d.client.Alter(context.Background(), &api.Operation{
		Schema: tp.Rdf(),
	})
	return err
}

func (d *Client) DropType(name string) error {
	err := d.client.Alter(context.Background(), &api.Operation{
		DropValue:       name,
		DropOp:          api.Operation_TYPE,
		RunInBackground: false,
	})
	return err
}

type Txn struct {
	Txn      *dgo.Txn
	Ctx      context.Context
	Readonly bool
	Closed   bool
}

func (d *Txn) GetSchemaAll() ([]Pred, []Type, error) {
	const q = `schema{}`
	var r Schema
	resp, err := d.Txn.Query(context.Background(), q)
	if err != nil {
		return nil, nil, err
	}
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, nil, err
	}
	return r.ListPred(), r.ListType(), err
}

func (d *Txn) GetPred(pred string) (*Pred, error) {
	const q = `schema(pred: %s){}`
	var r Schema

	resp, err := d.Txn.Query(d.Ctx, fmt.Sprintf(q, pred))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}
	if len(r.Preds) == 0 {
		return nil, errors.New("nothing found")
	}
	p := r.Preds[0]
	return &p, nil
}

func (d *Txn) GetType(tp string) (*Type, error) {
	const q = `schema(type: %s){}`
	var r Schema
	resp, err := d.Txn.Query(d.Ctx, fmt.Sprintf(q, tp))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}
	if len(r.Types) == 0 {
		return nil, errors.New("nothing found")
	}
	p := r.ListType()[0]
	return &p, nil
}
