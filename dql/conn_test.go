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
)

func TestH(t *testing.T) {
	s := strings.FieldsFunc("(NOT A OR B) AND (C AND NOT (D OR E))", func(r rune) bool {
		if r == '(' || r == ')' {
			return true
		}
		return false
	})
	t.Log(s)
}

func TestR(t *testing.T) {
	a := "hello world"
	i := strings.Index(a, "hello") + len("hello")
	fmt.Println(i)
	fmt.Println(a[0:i])
	fmt.Println(a[i:])
}
