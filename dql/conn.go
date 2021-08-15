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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"strings"
	"time"
)

type Config struct {
	Targets     []string      `json:"targets"`
	Username    string        `json:"username,omitempty"`
	Password    string        `json:"password,omitempty"`
	DialTimeout time.Duration `json:"dial_timeout,omitempty"`
	OptTimeout  time.Duration `json:"opt_timeout,omitempty"`
	Tls         Tls           `json:"tls"`
}

type Tls struct {
	ServeName  string `json:"tls_serve_name"`
	CaCert     string `json:"tls_ca_file"`
	ClientCert string `json:"tls_client_file"`
	ClientKey  string `json:"tls_client_key"`
}

// NewClient 新建dgraph连接
// target格式为 192.168.1.100:9080
// 第一个返回值为dgraph操作对象
func NewClient(config Config) (*Client, error) {
	var (
		clients []api.DgraphClient
		ctx     = context.Background()
		cancel  context.CancelFunc
		opts    []grpc.DialOption
	)
	if len(config.Targets) == 0 {
		return nil, errors.New("no target given")
	}
	if config.DialTimeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, config.DialTimeout)
		defer cancel()
	}
	if config.Tls == (Tls{}) {
		opts = append(opts, grpc.WithInsecure())
	} else {
		cred, err := newTlsCred(config.Tls)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}
	for _, target := range config.Targets {
		grpcConn, err := grpc.DialContext(ctx, target, opts...)
		if err != nil {
			return nil, err
		}
		client := api.NewDgraphClient(grpcConn)
		clients = append(clients, client)
	}
	dgraph := dgo.NewDgraphClient(clients...)
	return &Client{client: dgraph}, nil
}

func newTlsCred(ts Tls) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(ts.ClientCert, ts.ClientKey)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(ts.CaCert)
	if err != nil {
		return nil, err
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, errors.New("failed to append ca certs")
	}
	creds := credentials.NewTLS(&tls.Config{
		ServerName:   ts.ServeName,
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	})
	return creds, nil
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
