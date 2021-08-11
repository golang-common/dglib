/**
 * @Author: daipengyuan
 * @Description: <subject> <predicate> <object> 三元组
 * @File:  spo
 * @Version: 1.0.0
 * @Date: 2021/8/11 09:51
 */

package dql

import (
	"fmt"
	"strings"
)

// Facet 面，即边的属性
type Facet struct {
	Key string
	Val interface{} // 支持string,bool,int,float,datetime
}

// Nquad 即spo三元组
type Nquad struct {
	S      string      // subject,即uid
	P      Pred        // predicate
	O      interface{} // object
	Lang   string      // 语言
	Facets []Facet     // 面所属的key-value对
}

// Rdf 转换为rdf格式
func (n Nquad) Rdf() string {
	var (
		model    = "<$subject> <$predicate> $object$lang $facets ."
		facets   []string
		facetStr string
		lang     string
		obj      string
	)
	if n.Lang != "" {
		lang = fmt.Sprintf("@%s", n.Lang)
	}
	switch n.P.Type {
	case TUid:
		obj = fmt.Sprintf("<%s>", n.O)
	default:
		obj = fmt.Sprintf(`"%s"`, n.O)
	}
	for _, facet := range n.Facets {
		facets = append(facets, fmt.Sprintf("%s=%s", facet.Key, facet.Val))
	}
	if len(facets) > 0 {
		facetStr = fmt.Sprintf("(%s)", strings.Join(facets, ","))
	}
	replacer := strings.NewReplacer(
		"$subject", n.S,
		"$predicate", n.P.Name,
		"$object", obj,
		"$lang", lang,
		"$facets", facetStr,
	)
	return replacer.Replace(model)
}

// Node 一个节点,至少包含一个或多个有向边
type Node struct {
	Uid     string
	Type    string
	Nquards []Nquad
}

func (n Node) Rdf() string {
	var rdfList []string
	for _, nq := range n.Nquards {
		nq.S = n.Uid
		rdfList = append(rdfList, nq.Rdf())
	}
	if n.Type != "" {
		rdfList = append(rdfList, fmt.Sprintf(`%s <dgraph.type> "%s" .`, n.Uid, n.Type))
	}
	return fmt.Sprintf("%s", strings.Join(rdfList, "\n"))
}
