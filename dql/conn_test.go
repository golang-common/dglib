/**
 * @Author: daipengyuan
 * @Description:
 * @File:  conn_test
 * @Version: 1.0.0
 * @Date: 2021/8/13 17:05
 */

package dql

import (
	"strings"
	"testing"
)

func TestH(t *testing.T) {
	s :=strings.FieldsFunc("(NOT A OR B) AND (C AND NOT (D OR E))", func(r rune) bool {
		if r == '(' || r == ')' {
			return true
		}
		return false
	})
	t.Log(s)
}
