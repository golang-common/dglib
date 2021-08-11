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

type Type struct {
	Name   string   `json:"name"`
	Fields []string `json:"fields,omitempty"`
}

func (t Type) Rdf() string {
	var (
		preds []string
		r     string
	)
	for _, pred := range t.Fields {
		if strings.HasPrefix(pred, "~") {
			pred = fmt.Sprintf("<%s>", pred)
		}
		preds = append(preds, pred)
	}
	if len(preds) > 0 {
		r = fmt.Sprintf("type %s{\n\t%s\n}", t.Name, strings.Join(preds, "\n\t"))
	}
	return r
}

func (t *Type) Unmarshal(s string) error {
	var (
		name      string
		plist     []string
		rplist    []string
		tbody     string
		errFormat = errors.New("error type format")
	)
	if len(s) < 12 {
		return errFormat
	}
	s = strings.Trim(s, " \n")
	bsIndex := strings.Index(s, "{")
	beIndex := strings.Index(s, "}")
	if bsIndex == -1 || beIndex == -1 || bsIndex >= beIndex {
		return errFormat
	}
	name = s[5:bsIndex]
	if name == "" {
		return errFormat
	}
	tbody = s[bsIndex+1 : beIndex]
	tbody = strings.Trim(tbody, " \n\t")
	if len(tbody) == 0 {
		return errFormat
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
		return errFormat
	}
	t.Name = name
	t.Fields = append(plist, rplist...)
	return nil
}
