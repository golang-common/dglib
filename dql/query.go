/**
 * @Author: daipengyuan
 * @Description: 查询抽象
 * @File:  query
 * @Version: 1.0.0
 * @Date: 2021/8/12 11:03
 */

package dql

import (
	"errors"
	"fmt"
	"strings"
)

// Query 查询结构体，用于发送查询请求
// Q 查询主结构
type Query struct {
	Q           string            // 查询主体,展示项以及如何展示不好做抽象，需要用户自己定义
	Pager       *Pager            `json:"pager"`
	Recurse     *Recurse          `json:"recurse"`
	Sorter      []Sorter          `json:"sorter"`
	RootFilter  *Filter           `json:"root_filter"`
	PredFilter  map[string]Filter `json:"pred_filter"`  // 注意,key=谓词名
	FacetFilter map[string]Filter `json:"facet_filter"` // 注意,key=谓词名
}

func (q Query) Parse() (string, error) {
	const (
		rppager   = "$pager"
		rpsorter  = "$sorter"
		rprecurse = "$recurse"
	)
	var (
		r       = q.Q
		pager   string
		sorter  string
		recurse string
	)
	// 分页器解析
	if strings.Contains(r, rppager) {
		if q.Pager == nil || q.Pager.String() == "" {
			return "", errors.New("query has $pager but pager parse failed")
		}
		pager = q.Pager.String()
	}
	// 排序器解析
	if strings.Contains(r, rpsorter) {
		if q.Recurse == nil || q.Recurse.String() == "" {
			return "", errors.New("query has $recurse but recurse parse failed")
		}
	}
}

// Recurse 递归，Depth 递归深度，Loop 是否循环自身
type Recurse struct {
	Depth int
	Loop  bool
}

func (d Recurse) String() string {
	var r string
	if d.Depth > 0 {
		r = fmt.Sprintf("@recurse(depth:%d,loop:%t)", d.Depth, d.Loop)
	}
	return r
}

// Pager 分页器 在after之后偏移offset个结果取first个结果
type Pager struct {
	First  int
	Offset int
	After  int64
}

func (d Pager) String() string {
	var r []string
	if d.First > 0 {
		r = append(r, fmt.Sprintf("first: %d", d.First))
	}
	if d.Offset > 0 {
		r = append(r, fmt.Sprintf("offset: %d", d.Offset))
	}
	if d.After > 0 {
		r = append(r, fmt.Sprintf("after: %d", d.After))
	}
	return strings.Join(r, ",")
}

// Sorter 排序器，Order=排序方向，Orderby=排序的目标谓词
type Sorter struct {
	Order   string // orderasc or orderdesc
	Orderby string
}

func (d Sorter) Parse() (string, error) {
	if d.Order != "orderasc" && d.Order != "orderdesc" {
		return "", errors.New("unsupport order func,need [orderasc] or [orderdesc]")
	}
	pred, ok := PredMap[d.Orderby]
	if !ok {
		return "", errors.New("sort predicate not exist " + d.Orderby)
	}
	ptattr, _ := TypeAttrMap[pred.Type]
	if ptattr.Ts == nil {
		return "", errors.New(fmt.Sprintf("predicate [%s] type [%s],do not have sortable index", d.Orderby, pred.Type))
	}
	var sortable bool
	for k, v := range ptattr.Ts {
		for _, tk := range pred.Tokenizer {
			if k == tk && v.Stb == true {
				sortable = true
			}
		}
	}
	if !sortable {
		return "", errors.New(fmt.Sprintf("predicate [%s] type [%s],do not have sortable index", d.Orderby, pred.Type))
	}
	return fmt.Sprintf("%s:%s", d.Order, d.Orderby), nil
}
