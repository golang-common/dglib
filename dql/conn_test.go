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
	"testing"
	"time"
)

func TestConn(t *testing.T) {
	c, err := NewClient(Config{
		Targets:     []string{"localhost:9080"},
		DialTimeout: 2 * time.Second,
		OptTimeout:  2 * time.Second,
		Username:    "groot",
		Password:    "password",
		Tls: Tls{
			ServeName:  "crane",
			CaCert:     "/Users/lyonsdpy/Data/dgraph/tls/ca.crt",
			ClientCert: "/Users/lyonsdpy/Data/dgraph/tls/client.crane.crt",
			ClientKey:  "/Users/lyonsdpy/Data/dgraph/tls/client.crane.key",
		},
	})
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
