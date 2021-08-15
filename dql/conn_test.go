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
		Targets:     []string{"localhost:7080"},
		DialTimeout: 3 * time.Second,
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
	pds, tps, err := c.Txn().GetSchemaAll()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", pds))
	t.Log(fmt.Sprintf("%+v", tps))
}
