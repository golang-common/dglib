/**
 * @Author: daipengyuan
 * @Description: 谓词
 * @File:  pred
 * @Version: 1.0.0
 * @Date: 2021/8/11 12:18
 */

package dql

import (
	"fmt"
	"strings"
)

var PredMap = map[string]Pred{}

// Pred 类型/边
type Pred struct {
	Predicate string   `json:"predicate"`
	Type      string   `json:"type"`
	Index     bool     `json:"index"`
	Tokenizer []string `json:"tokenizer"`
	Reverse   bool     `json:"reverse"`
	Count     bool     `json:"count"`
	List      bool     `json:"list"`
	Upsert    bool     `json:"upsert"`
}

func (p Pred) String() string {
	return p.Predicate
}

// Rdf 转换为rdf格式
func (p Pred) Rdf() string {
	var (
		model   = `$name: $type $indices .`
		ptype   = p.Type
		indices []string
	)
	if p.Index && len(p.Tokenizer) > 0 {
		indices = append(indices, fmt.Sprintf("@index(%s)", strings.Join(p.Tokenizer, ",")))
	}
	if p.Reverse {
		indices = append(indices, "@reverse")
	}
	if p.Count {
		indices = append(indices, "@count")
	}
	if p.List {
		ptype = fmt.Sprintf("[%s]", ptype)
	}
	if p.Upsert {
		indices = append(indices, "@upsert")
	}

	replacer := strings.NewReplacer(
		"$name", p.Predicate,
		"$type", ptype,
		"$indices", strings.Join(indices, " "),
	)
	return replacer.Replace(model)
}
