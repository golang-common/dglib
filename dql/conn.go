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
	if config.Username != "" && config.Password != "" {
		err := dgraph.Login(ctx, config.Username, config.Password)
		if err != nil {
			return nil, err
		}
	}
	return &Client{client: dgraph, optTimeout: config.OptTimeout}, nil
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
	Types []Type `json:"types"`
}

// SkipSysSchema 忽略dgraph系统自身schema
func (s *Schema) SkipSysSchema() Schema {
	var (
		r     Schema
		preds []Pred
		types []Type
	)
	for _, p := range s.Preds {
		if strings.HasPrefix(p.Predicate, "dgraph.") {
			continue
		}
		preds = append(preds, p)
	}
	for _, v := range s.Types {
		if strings.HasPrefix(v.Name, "dgraph.") {
			continue
		}
		types = append(types, v)
	}
	r.Preds = preds
	r.Types = types
	return r
}

type Client struct {
	client     *dgo.Dgraph
	optTimeout time.Duration
}

func (d *Client) Txn(ReadOnly ...bool) *Txn {
	if len(ReadOnly) > 0 && ReadOnly[0] == true {
		return &Txn{Txn: d.client.NewReadOnlyTxn(), Readonly: ReadOnly[0], Timeout: d.optTimeout}
	}
	return &Txn{Txn: d.client.NewTxn(), Timeout: d.optTimeout}
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
		Schema: tp.Schema(),
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
	Timeout  time.Duration
	Readonly bool
	cancel   context.CancelFunc
}

func (d *Txn) Ctx() context.Context {
	var (
		r = context.Background()
		c context.CancelFunc
	)
	// 清除之前的ctx资源
	if d.cancel != nil {
		d.cancel()
	}
	if d.Timeout > 0 {
		r, c = context.WithTimeout(r, d.Timeout)
	}
	d.cancel = c
	return r
}

func (d *Txn) Cancel() {
	if d.cancel != nil {
		d.cancel()
	}
}

func (d *Txn) CommitOrAbort(err error) error {
	defer d.Cancel()
	if err != nil {
		return d.Txn.Discard(d.Ctx())
	}
	return d.Txn.Commit(d.Ctx())
}

// GetSchema 获取dgraph所有谓词和类型
func (d *Txn) GetSchema() (*Schema, error) {
	const q = `schema{}`
	var res Schema
	defer d.Cancel()
	resp, err := d.Txn.Query(d.Ctx(), q)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		return nil, err
	}
	r := res.SkipSysSchema()
	return &r, err
}

// FindPred 查找特定谓词结构,如果不存在则报错
func (d *Txn) FindPred(pred string) (*Pred, error) {
	const q = `schema(pred: %s){}`
	var res Schema
	defer d.Cancel()
	resp, err := d.Txn.Query(d.Ctx(), fmt.Sprintf(q, pred))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		return nil, err
	}
	if len(res.Preds) == 0 {
		return nil, errors.New("not found")
	}
	p := res.Preds[0]
	return &p, nil
}

// FindType 查找特定类型,如果不存在则报错
func (d *Txn) FindType(tp string) (*Type, error) {
	const q = `schema(type: %s){}`
	var res Schema
	defer d.Cancel()
	resp, err := d.Txn.Query(d.Ctx(), fmt.Sprintf(q, tp))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		return nil, err
	}
	if len(res.Types) == 0 {
		return nil, errors.New("not found")
	}
	p := res.Types[0]
	return &p, nil
}
//
//// FindTypePreds 输入特定类型,查找其下所有关联谓词(第二个参数指示是否忽略反向谓词)
//func (d *Txn) FindTypePreds(tp Type, ignoreReverse ...bool) ([]Pred, error) {
//	var ignore bool
//	if len(ignoreReverse) > 0 && ignoreReverse[0] == true {
//		ignore = true
//	}
//
//}
