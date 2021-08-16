/**
 * @Author: daipengyuan
 * @Description:
 * @File:  type
 * @Version: 1.0.0
 * @Date: 2021/8/11 12:18
 */

package dql

import (
	"errors"
	"fmt"
	"strings"
)

var TypeMap = map[string]Type{}

type Type struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields,omitempty"`
}

type Field struct {
	Name string `json:"name"`
}

func (t Type) Schema() string {
	var (
		preds []string
		r     string
	)
	for _, p := range t.Fields {
		var pr = p.Name
		if strings.HasPrefix(p.Name, "~") {
			pr = fmt.Sprintf("<~%s>", p.Name)
		}
		preds = append(preds, pr)
	}
	if len(preds) > 0 {
		r = fmt.Sprintf("type %s{\n\t%s\n}", t.Name, strings.Join(preds, "\n\t"))
	}
	return r
}

func UnmarshalType(s string) (*Type, error) {
	var (
		name      string
		plist     []string
		rplist    []string
		tbody     string
		errFormat = errors.New("error type format")
	)
	if len(s) < 12 {
		return nil, errFormat
	}
	s = strings.Trim(s, " \n")
	bsIndex := strings.Index(s, "{")
	beIndex := strings.Index(s, "}")
	if bsIndex == -1 || beIndex == -1 || bsIndex >= beIndex {
		return nil, errFormat
	}
	name = s[5:bsIndex]
	if name == "" {
		return nil, errFormat
	}
	tbody = s[bsIndex+1 : beIndex]
	tbody = strings.Trim(tbody, " \n\t")
	if len(tbody) == 0 {
		return nil, errFormat
	}
	for _, line := range strings.Split(tbody, "\n") {
		line = strings.Trim(line, " \t")
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "<~") && len(line) > 3 {
			rplist = append(rplist, line[1:len(line)-1])
		}
		plist = append(plist, line)
	}
	if len(plist) == 0 {
		return nil, errFormat
	}
	t := new(Type)
	t.Name = name
	for _, p := range plist {
		t.Fields = append(t.Fields, Field{Name: p})
	}
	return t, nil
}
