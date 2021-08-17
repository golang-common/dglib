/**
 * @Author: daipengyuan
 * @Description: 变更操作基础抽象
 * @File:  mutation
 * @Version: 1.0.0
 * @Date: 2021/8/15 16:00
 */

package dql

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/dgo/v200/protos/api"
	uuid "github.com/satori/go.uuid"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	StarAll   = "_STAR_ALL"
	StructTag = "db" // 解析结构时使用的tag

	tagId   = "id"
	tagMust = "must"

	actionAdd     = "add"
	actionMerge   = "merge"
	actionUpdate  = "update"
	actionDelete  = "delete"
	actionDelNode = "delnode"
)

// AddNode 给入一个结构体，其UID字段可以为空,系统会自动生成UID
// 结构体的dgraph tag中，各字段解释如下
// id:字段，则表示该字段为结构体的唯一键，且不能为空，需要先查重
// must:必备字段，该字段必须不为空，可以修改，但是无法删除和置空
func (d *Txn) AddNode(obj interface{}) (*api.Response, error) {
	var (
		q    string
		cond string
	)
	nqList, idList, err := d.nquad(obj, actionAdd, true)
	if err != nil {
		return nil, err
	}
	if len(idList) > 3 {
		model := `query{a as var(func: eq($pred,$value)) @filter($typefilter) }`
		typeList := idList[2:]
		var tpExprList []string
		for _, tp := range typeList {
			tpExprList = append(tpExprList, fmt.Sprintf("type(%s)", tp))
		}
		tpExprStr := strings.Join(tpExprList, " OR ")
		replacer := strings.NewReplacer(
			"$pred", idList[0],
			"$value", idList[1],
			"$typefilter", tpExprStr,
		)
		q = replacer.Replace(model)
		cond = `@if(eq(uid(a),0))`
	}
	req := &api.Request{
		Query:      q,
		Mutations:  []*api.Mutation{{Set: nqList, Cond: cond}},
		RespFormat: 0,
	}
	defer d.Cancel()
	resp, err := d.Txn.Do(d.Ctx(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateNode 更新节点，节点上的空值也会被更新
// 注意：如果更新了带面的边，则需要带上相应面的值，否则面值会被删除
func (d *Txn) UpdateNode(obj interface{}) (*api.Response, error) {
	var (
		q    string
		cond string
	)
	nqList, idList, err := d.nquad(obj, actionUpdate, false)
	if err != nil {
		return nil, err
	}
	if len(idList) > 3 {
		model := `query{a as var(func: eq($pred,$value)) @filter($typefilter) }`
		typeList := idList[2:]
		var tpExprList []string
		for _, tp := range typeList {
			tpExprList = append(tpExprList, fmt.Sprintf("type(%s)", tp))
		}
		tpExprStr := strings.Join(tpExprList, " OR ")
		replacer := strings.NewReplacer(
			"$pred", idList[0],
			"$value", idList[1],
			"$typefilter", tpExprStr,
		)
		q = replacer.Replace(model)
		cond = `@if(eq(uid(a),0))`
	}
	req := &api.Request{
		Query:      q,
		Mutations:  []*api.Mutation{{Set: nqList, Cond: cond}},
		RespFormat: 0,
	}
	defer d.Cancel()
	resp, err := d.Txn.Do(d.Ctx(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// MergeNode 合并节点，节点中的空值会被忽略
// 注意：如果更新了带面的边，则需要带上相应面的值，否则面值会被删除
func (d *Txn) MergeNode(obj interface{}) (*api.Response, error) {
	var (
		q    string
		cond string
	)
	nqList, idList, err := d.nquad(obj, actionMerge, true)
	if err != nil {
		return nil, err
	}
	if len(idList) > 3 {
		model := `query{a as var(func: eq($pred,$value)) @filter($typefilter) }`
		typeList := idList[2:]
		var tpExprList []string
		for _, tp := range typeList {
			tpExprList = append(tpExprList, fmt.Sprintf("type(%s)", tp))
		}
		tpExprStr := strings.Join(tpExprList, " OR ")
		replacer := strings.NewReplacer(
			"$pred", idList[0],
			"$value", idList[1],
			"$typefilter", tpExprStr,
		)
		q = replacer.Replace(model)
		cond = `@if(eq(uid(a),0))`
	}
	req := &api.Request{
		Query:      q,
		Mutations:  []*api.Mutation{{Set: nqList, Cond: cond}},
		RespFormat: 0,
	}
	defer d.Cancel()
	resp, err := d.Txn.Do(d.Ctx(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DelVal 删除obj中的相应值，但不会删除节点本身
// 删除操作会忽略面操作
func (d *Txn) DelVal(obj interface{}) (*api.Response, error) {
	var (
		q    string
		cond string
	)
	nqList, idList, err := d.nquad(obj, actionDelete, true)
	if err != nil {
		return nil, err
	}
	if len(idList) > 3 {
		model := `query{a as var(func: eq($pred,$value)) @filter($typefilter) }`
		typeList := idList[2:]
		var tpExprList []string
		for _, tp := range typeList {
			tpExprList = append(tpExprList, fmt.Sprintf("type(%s)", tp))
		}
		tpExprStr := strings.Join(tpExprList, " OR ")
		replacer := strings.NewReplacer(
			"$pred", idList[0],
			"$value", idList[1],
			"$typefilter", tpExprStr,
		)
		q = replacer.Replace(model)
		cond = `@if(eq(uid(a),0))`
	}
	req := &api.Request{
		Query:      q,
		Mutations:  []*api.Mutation{{Del: nqList, Cond: cond}},
		RespFormat: 0,
	}
	defer d.Cancel()
	resp, err := d.Txn.Do(d.Ctx(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DelNode 删除该节点的所有出方向边，同时删除所有指向该节点的边
func (d *Txn) DelNode(obj interface{}) (*api.Response, error) {
	var (
		q           string
		muList      []*api.Mutation
		revPredList []string
		uid         string
		nquadList   []*api.NQuad
	)
	structVal, err := d.getStructObj(obj)
	if err != nil {
		return nil, err
	}
	uid = d.getUidFromStructVal(structVal)
	if uid == "" {
		return nil, errors.New(fmt.Sprintf("get uid failed,operation [%s] cannot operate without uid", actionDelNode))
	}
	nquadList = append(nquadList,
		&api.NQuad{Subject: uid, Predicate: StarAll, ObjectValue: &api.Value{Val: &api.Value_DefaultVal{DefaultVal: StarAll}}})
	muList = append(muList, &api.Mutation{Set: nquadList})
	for i := 0; i < structVal.NumField(); i++ {
		field := structVal.Type().Field(i)
		tagList := strings.Split(field.Tag.Get(StructTag), ",")
		if len(tagList) == 0 {
			continue
		}
		pred := tagList[0]
		// 面无法删除
		if strings.Contains(pred, "|") {
			continue
		}
		if pred == "uid" {
			continue
		}
		// 记录反向边
		if strings.HasPrefix(pred, "~") {
			revPredList = append(revPredList, pred)
		}
	}
	if len(revPredList) > 0 {
		var revList []string
		var rmap = make(map[string]struct{})
		q = `query{
	var(func: uid($uid)){
		$revpred
	}
}`
		for _, revv := range revPredList {
			var rvar string
			for {
				rvar = GetRandomString(3)
				if _, ok := rmap[rvar]; !ok {
					break
				}
			}
			rmap[rvar] = struct{}{}
			revList = append(revList, fmt.Sprintf("%s as %s", rvar, revv))
			mu := &api.Mutation{
				Cond: fmt.Sprintf(`@if(gt(uid(%s),0))`, rvar),
				Set: []*api.NQuad{{
					Subject:     uid,
					Predicate:   strings.Trim(revv, "~"),
					ObjectValue: &api.Value{Val: &api.Value_DefaultVal{DefaultVal: StarAll}},
				}},
			}
			muList = append(muList, mu)
		}
		rplc := strings.NewReplacer(
			"$uid", uid,
			"$revpred", strings.Join(revList, "\n\t\t"),
		)
		q = rplc.Replace(q)
	}
	req := &api.Request{
		Query:     q,
		Mutations: muList,
	}
	defer d.Cancel()
	resp, err := d.Txn.Do(d.Ctx(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// nquad 将结构体转换为dgraph操作对象
// action 为操作类型 ，add/merge/update/delete
// omitempty 为忽略空值，通常由具体方法控制
// 如果传入的结构体中有id字段，则需要返回唯一校验 pred/val/type的数组,val只支持 string和int类型
func (d *Txn) nquad(obj interface{}, action string, omitempty ...bool) ([]*api.NQuad, []string, error) {
	var (
		subjectId    string
		idArr        []string
		dtype        []string
		singleNquad  = make(map[string]*api.NQuad)
		multiNquad   = make(map[string][]*api.NQuad)
		singleFacets = make(map[string][]*api.Facet)
		multiFacets  = make(map[string]map[string]map[int]*api.Facet)
	)
	// 解析obj为struct值，如果最后无法解析为struct值，则报错
	structVal, err := d.getStructObj(obj)
	if err != nil {
		return nil, nil, err
	}
	// 获取dtype列表，获取失败则报错，操作的obj必须包含dtype
	dtype = d.getDtypeFromStruct(structVal)
	if len(dtype) == 0 {
		return nil, nil, errors.New("cannot get dtype")
	}
	// 如果为增加操作，则指定空uid；其它操作都要获取uid，获取不到则会报错
	if action == actionAdd {
		subjectId = fmt.Sprintf("_:%s", uuid.NewV1().String())
	} else {
		subjectId = d.getUidFromStructVal(structVal)
		if subjectId == "" {
			return nil, nil, errors.New(fmt.Sprintf("get uid failed,operation [%s] cannot operate without uid", action))
		}
	}
	// 遍历解析结构体
	for i := 0; i < structVal.NumField(); i++ {
		var (
			isFacet  bool
			idSet    bool
			mustSet  bool
			predName string
			facetKey string
		)
		field := structVal.Type().Field(i)
		fv := structVal.Field(i)
		// 如果是指针则解指针
		if fv.Kind() == reflect.Ptr {
			fv = fv.Elem()
		}
		tagList := strings.Split(field.Tag.Get(StructTag), ",")
		// dgraph tag为空则不解析该字段
		if len(tagList) == 0 {
			continue
		}
		// 解析结构体中tag字段
		for k, tag := range tagList {
			if k == 0 {
				predName = tag
				if strings.Contains(tag, "|") {
					tl := strings.Split(tag, "|")
					if len(tl) != 2 {
						return nil, nil, errors.New(fmt.Sprintf("wrong tag format of [%s]", tag))
					}
					predName = tl[0]
					facetKey = tl[1]
					isFacet = true
				}
				continue
			}
			if tag == tagId {
				if idSet == true {
					return nil, nil, errors.New("only 1 id directive can be set in type")
				}
				idSet = true
				continue
			}
			if tag == tagMust {
				mustSet = true
				continue
			}
		}
		// 越过uid分析,uid为subject，已在开头获取
		if predName == "uid" {
			continue
		}
		// 无法单独改/删 dgraph.type(除非删除节点)
		if predName == "dgraph.type" && action != actionAdd {
			continue
		}
		// id或must设置时，说明该字段必须存在
		if (idSet || mustSet) && fv.IsZero() && (action == actionAdd || action == actionUpdate) {
			return nil, nil, errors.New(fmt.Sprintf("id or must set while [%s] value is nil", predName))
		}
		if (idSet || mustSet) && !fv.IsZero() && action == actionDelete {
			return nil, nil, errors.New(fmt.Sprintf("id or must not set while [%s] value is nil on delete", predName))
		}
		if fv.IsZero() && len(omitempty) > 0 && omitempty[0] == true {
			continue
		}
		// 如果设置了id唯一键，则返回类型和id的值，在外部的方法中需要
		if idSet {
			switch fv.Interface().(type) {
			case string, int, int64:
				idArr = append(idArr, predName)
				idArr = append(idArr, fmt.Sprintf("%v", fv.Interface()))
				idArr = append(idArr, dtype...)
			default:
				return nil, nil, errors.New(fmt.Sprintf("id directive does not surpport set on [%s] datetype", fv.Type()))
			}
		}
		// 解析谓词对象
		if !isFacet {
			if fv.Kind() == reflect.Slice {
				for j := 0; j < fv.Len(); j++ {
					fvval := fv.Index(j)
					nquad, err := d.parseObjNquad(subjectId, predName, fvval.Interface())
					if err != nil {
						return nil, nil, err
					}
					multiNquad[predName] = append(multiNquad[predName], nquad)
				}
			} else {
				nquad, err := d.parseObjNquad(subjectId, predName, fv.Interface())
				if err != nil {
					return nil, nil, err
				}
				singleNquad[predName] = nquad
			}
			continue
		}
		// 解析面对象
		if isFacet && action != actionDelete {
			if fv.Kind() == reflect.Map {
				for _, mk := range fv.MapKeys() {
					mval := fv.MapIndex(mk)
					var mkey int
					switch mk.Interface().(type) {
					case int, int64:
						mkey = mk.Interface().(int)
					case string:
						ms := mk.Interface().(string)
						mi, err := strconv.Atoi(ms)
						if err != nil {
							return nil, nil, err
						}
						mkey = mi
					default:
						return nil, nil, errors.New(fmt.Sprintf("unsupport facets map key type [%s]", mk.Type()))
					}
					if mkey == 0 {
						return nil, nil, errors.New("parse facet map error, empty map key")
					}
					facet, err := d.parseObjFacet(facetKey, mval.Interface())
					if err != nil {
						return nil, nil, err
					}
					if multiFacets[predName] == nil {
						multiFacets[predName] = make(map[string]map[int]*api.Facet)
						if multiFacets[predName][facetKey] == nil {
							multiFacets[predName][facetKey] = make(map[int]*api.Facet)
						}
					}
					multiFacets[predName][facetKey][mkey] = facet
				}
			} else {
				facet, err := d.parseObjFacet(facetKey, fv.Interface())
				if err != nil {
					return nil, nil, err
				}
				singleFacets[predName] = append(singleFacets[predName], facet)
			}
		}
	}
	var r []*api.NQuad
	for singlekey, singleval := range singleNquad {
		if sfacet, ok := singleFacets[singlekey]; ok &&
			len(singleFacets[singlekey]) > 0 && action != actionDelete {
			singleval.Facets = sfacet
		}
		r = append(r, singleval)
	}
	for multikey, multival := range multiNquad {
		if mfacet, ok := multiFacets[multikey]; ok && len(mfacet) > 0 && action != actionDelete {
			for _, mfmap := range mfacet {
				for k, mf := range mfmap {
					if k < len(multival) {
						multival[k].Facets = append(multival[k].Facets, mf)
					} else {
						return nil, nil, errors.New(fmt.Sprintf("error facets key [%d], with pred val len [%d]", k, len(multival)))
					}
				}
			}
		}
		r = append(r, multival...)
	}
	return r, idArr, nil
}

// getStructObj 将传入的对象解析为结构体val,用于后续解析
func (d *Txn) getStructObj(obj interface{}) (reflect.Value, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return reflect.Value{}, errors.New("object is not a struct value")
	}
	return val, nil
}

// getDtypeFromStruct 从传入的结构体对象中获取dgraph.type字段的值
func (d *Txn) getDtypeFromStruct(value reflect.Value) []string {
	dtpVal := value.FieldByName("dgraph.type")
	dv := dtpVal.Interface()
	if v, ok := dv.([]string); ok {
		return v
	}
	return nil
}

// getUidFromStructVal 从传入的结构体对象中获取uid字段的值
func (d *Txn) getUidFromStructVal(value reflect.Value) string {
	uidVal := value.FieldByName("uid")
	uv := uidVal.Interface()
	switch uv.(type) {
	case string:
		return uv.(string)
	case int, int64:
		return fmt.Sprintf("%x", uv.(int64))
	}
	return ""
}

// parseObjNquad 解析一个单一的字面量(string,int64,float64,bool,struct)，或(geom.T,time.Time)
func (d *Txn) parseObjNquad(subject, pred string, obj interface{}) (*api.NQuad, error) {
	var (
		name string
		lang string
	)
	name = pred
	if strings.Contains(name, "@") {
		tnl := strings.Split(pred, "@")
		if len(tnl) != 2 {
			return nil, errors.New(fmt.Sprintf("error lang set in tag [%s]", pred))
		}
		name = tnl[0]
		lang = tnl[1]
	}
	r := &api.NQuad{
		Subject:   subject,
		Predicate: name,
		Lang:      lang,
	}
	switch obj.(type) {
	case string:
		r.ObjectValue = &api.Value{Val: &api.Value_StrVal{StrVal: obj.(string)}}
	case int, int64:
		r.ObjectValue = &api.Value{Val: &api.Value_IntVal{IntVal: obj.(int64)}}
	case float32, float64:
		r.ObjectValue = &api.Value{Val: &api.Value_DoubleVal{DoubleVal: obj.(float64)}}
	case bool:
		r.ObjectValue = &api.Value{Val: &api.Value_BoolVal{BoolVal: obj.(bool)}}
	case time.Time:
		tb, err := obj.(time.Time).MarshalBinary()
		if err != nil {
			return nil, err
		}
		r.ObjectValue = &api.Value{Val: &api.Value_DatetimeVal{DatetimeVal: tb}}
	case geom.T:
		gb, err := geojson.Marshal(obj.(geom.T))
		if err != nil {
			return nil, err
		}
		r.ObjectValue = &api.Value{Val: &api.Value_GeoVal{GeoVal: gb}}
	}
	if r.ObjectValue != nil {
		return r, nil
	}
	if reflect.TypeOf(obj).Kind() == reflect.Struct {
		val := reflect.ValueOf(obj)
		tp := reflect.TypeOf(obj)
		subUid := d.getUidFromStructVal(val)
		if subUid == "" {
			return nil, errors.New(fmt.Sprintf("pred [%s] must have a uid when exist", name))
		}
		r.ObjectId = subUid
		for i := 0; i < val.NumField(); i++ {
			subTag := tp.Field(i).Tag.Get(StructTag)
			if strings.HasPrefix(subTag, pred+"|") && len(subTag) > len(pred+"|") {
				key := subTag[len(pred):]
				facet, err := d.parseObjFacet(key, val.Field(i).Interface())
				if err != nil {
					return nil, err
				}
				r.Facets = append(r.Facets, facet)
			}
		}
		return r, nil
	}
	return nil, errors.New(fmt.Sprintf("unsupport datatype [%s] on [%s]", reflect.TypeOf(obj), name))
}

// parseObjFacet 解析单一字面量对应的面(string, bool, int, float and dateTime)
func (d *Txn) parseObjFacet(key string, obj interface{}) (*api.Facet, error) {
	var r = &api.Facet{
		Key: key,
	}
	// TODO:facet传入的数据类型有待测试
	switch obj.(type) {
	case string:
		r.Value = []byte(obj.(string))
		r.ValType = api.Facet_STRING
	case bool:
		r.Value = []byte(fmt.Sprintf("%t", obj.(bool)))
		r.ValType = api.Facet_BOOL
	case int, int64:
		r.Value = []byte(fmt.Sprintf("%d", obj.(int64)))
		r.ValType = api.Facet_INT
	case float64:
		r.Value = []byte(fmt.Sprintf("%f", obj.(float64)))
		r.ValType = api.Facet_FLOAT
	case time.Time:
		r.Value = []byte(fmt.Sprintf("%s", obj.(time.Time).String()))
		r.ValType = api.Facet_DATETIME
	default:
		return nil, errors.New(fmt.Sprintf("unsupport facet datatype [%s]", reflect.TypeOf(obj)))
	}
	return r, nil
}

func GetRandomString(l int) string {
	str := "abcdefghijklmnopqrstuvwxyz"
	bs := []byte(str)
	var rst []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		rst = append(rst, bs[r.Intn(len(bs))])
	}
	return string(rst)
}
