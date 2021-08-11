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
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
)

// New 新建dgraph连接
// target格式为 192.168.1.100:9080
// 第一个返回值为dgraph操作对象
// 第二个返回值为取消/关闭对象的方法
func New(target string) (*Dql, error) {
	grpcConn, err := grpc.DialContext(context.Background(), target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	dgClient := api.NewDgraphClient(grpcConn)
	dgraph := dgo.NewDgraphClient(dgClient)
	return &Dql{client: dgraph, cancel: func() error {
		return grpcConn.Close()
	}}, nil
}

type Dql struct {
	client *dgo.Dgraph
	cancel func() error
}
