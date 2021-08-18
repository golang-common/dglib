/**
 * @Author: daipengyuan
 * @Description: S P O 三元组
 * @File:  nquad
 * @Version: 1.0.0
 * @Date: 2021/8/15 16:25
 */

package dql

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/dgo/v200/protos/api"
	uuid "github.com/satori/go.uuid"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"reflect"
	"strings"
	"time"
)

const (
	Uid      = "Uid"
	StarAll  = "_STAR_ALL"
	TagDb    = "db" // 解析结构时使用的tag
	TagIndex = "index"
	TagDtype = "dtype"
	tagId    = "id"
	tagMust  = "must"
)

var (
	starNqVal = &api.Value{Val: &api.Value_DefaultVal{DefaultVal: StarAll}}

	typeNqTypeMap = map[string]func(value reflect.Value) (*api.Value, string, error){
		TypeDefault: func(value reflect.Value) (*api.Value, string, error) {
			if v, ok := value.Interface().(string); ok {
				return &api.Value{Val: &api.Value_DefaultVal{DefaultVal: v}}, "", nil
			}
			return nil, "", errors.New("errors datatype " + value.Type().String())
		},
		TypeUid: func(value reflect.Value) (*api.Value, string, error) {
			if v, ok := value.Interface().(string); ok {
				return nil, v, nil
			}
			if v, ok := value.Interface().(uint64); ok {
				return nil, fmt.Sprintf("0x%x", v), nil
			}
			return nil, "", errors.New("errors datatype " + value.Type().String())
		},
		TypeInt: func(value reflect.Value) (*api.Value, string, error) {
			if v, ok := value.Interface().(int); ok {
				return &api.Value{Val: &api.Value_IntVal{IntVal: int64(v)}}, "", nil
			}
			if v, ok := value.Interface().(int64); ok {
				return &api.Value{Val: &api.Value_IntVal{IntVal: v}}, "", nil
			}
			return nil, "", errors.New("errors datatype " + value.Type().String())
		},
		TypeFloat: func(value reflect.Value) (*api.Value, string, error) {
			if v, ok := value.Interface().(float64); ok {
				return &api.Value{Val: &api.Value_DoubleVal{DoubleVal: v}}, "", nil
			}
			if v, ok := value.Interface().(float32); ok {
				return &api.Value{Val: &api.Value_DoubleVal{DoubleVal: float64(v)}}, "", nil
			}
			return nil, "", errors.New("errors datatype " + value.Type().String())
		},
		TypeString: func(value reflect.Value) (*api.Value, string, error) {
			if v, ok := value.Interface().(string); ok {
				return &api.Value{Val: &api.Value_StrVal{StrVal: v}}, "", nil
			}
			return nil, "", errors.New("errors datatype " + value.Type().String())
		},
		TypeBool: func(value reflect.Value) (*api.Value, string, error) {
			if v, ok := value.Interface().(bool); ok {
				return &api.Value{Val: &api.Value_BoolVal{BoolVal: v}}, "", nil
			}
			return nil, "", errors.New("errors datatype " + value.Type().String())
		},
		TypeDateTime: func(value reflect.Value) (*api.Value, string, error) {
			if v, ok := value.Interface().(time.Time); ok {
				bs, err := v.MarshalBinary()
				if err != nil {
					return nil, "", err
				}
				return &api.Value{Val: &api.Value_DatetimeVal{DatetimeVal: bs}}, "", nil
			}
			return nil, "", errors.New("errors datatype " + value.Type().String())
		},
		TypeGeo: func(value reflect.Value) (*api.Value, string, error) {
			if v, ok := value.Interface().(geom.T); ok {
				bs, err := geojson.Marshal(v)
				if err != nil {
					return nil, "", err
				}
				return &api.Value{Val: &api.Value_GeoVal{GeoVal: bs}}, "", nil
			}
			return nil, "", errors.New("errors datatype " + value.Type().String())
		},
		TypePassword: func(value reflect.Value) (*api.Value, string, error) {
			if v, ok := value.Interface().(string); ok {
				return &api.Value{Val: &api.Value_PasswordVal{PasswordVal: v}}, "", nil
			}
			return nil, "", errors.New("errors datatype " + value.Type().String())
		},
	}
)

func (d *Txn) Add(obj interface{}, facets ...*Facet) (*api.Response, error) {
	muta, err := newMutation(obj, facets...)
	if err != nil {
		return nil, err
	}
	req, err := muta.MakeAdd()
	if err != nil {
		return nil, err
	}
	fmt.Println(IndentJson(req))
	defer d.Cancel()
	resp, err := d.Txn.Do(d.Ctx(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *Txn) Update(obj interface{}, facets ...*Facet) (*api.Response, error) {
	muta, err := newMutation(obj, facets...)
	if err != nil {
		return nil, err
	}
	req, err := muta.MakeUpd()
	if err != nil {
		return nil, err
	}
	defer d.Cancel()
	resp, err := d.Txn.Do(d.Ctx(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *Txn) Merge(obj interface{}, facets ...*Facet) (*api.Response, error) {
	muta, err := newMutation(obj, facets...)
	if err != nil {
		return nil, err
	}
	req, err := muta.MakeMerge()
	if err != nil {
		return nil, err
	}
	defer d.Cancel()
	resp, err := d.Txn.Do(d.Ctx(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *Txn) Delete(obj interface{}, facets ...*Facet) (*api.Response, error) {
	muta, err := newMutation(obj, facets...)
	if err != nil {
		return nil, err
	}
	req, err := muta.MakeDelVal()
	if err != nil {
		return nil, err
	}
	defer d.Cancel()
	resp, err := d.Txn.Do(d.Ctx(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *Txn) DelNode(obj interface{}, facets ...*Facet) (*api.Response, error) {
	muta, err := newMutation(obj, facets...)
	if err != nil {
		return nil, err
	}
	req, err := muta.MakeDelNode()
	if err != nil {
		return nil, err
	}
	defer d.Cancel()
	resp, err := d.Txn.Do(d.Ctx(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func newMutation(obj interface{}, facets ...*Facet) (*mutation, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, errors.New("obj must be struct type")
	}
	uid, dtype, err := parseUidField(val)
	if err != nil {
		return nil, err
	}
	if dtype == "" {
		return nil, errors.New("obj must have Uid field with dtype tag")
	}
	mu := &mutation{
		Subject: uid,
		Dtype:   dtype,
		Val:     val,
		Facets:  facets,
	}
	return mu, nil
}

type MuType int

// mutation 变更所需数据
type mutation struct {
	Subject    string
	Dtype      string
	Val        reflect.Value
	Facets     []*Facet
	curName    string
	curDt      string
	curPred    string
	curLang    string
	curReverse bool
	curMustSet bool
	idSet      bool
	idName     string
	idVal      string
}

func (m *mutation) MakeAdd() (*api.Request, error) {
	var (
		model     = `query{ a as var(func: type($type)) @filter(eq($name,$value)) }`
		q         string
		cond      string
		setNquads []*api.NQuad
	)
	m.Subject = fmt.Sprintf("_:%s", uuid.NewV1().String())
	setNquads = append(setNquads, &api.NQuad{
		Subject:     m.Subject,
		Predicate:   "dgraph.type",
		ObjectValue: &api.Value{Val: &api.Value_StrVal{StrVal: m.Dtype}},
	})
	for i := 0; i < m.Val.NumField(); i++ {
		f := m.Val.Type().Field(i)
		fv := m.Val.Field(i)
		if f.Name == Uid {
			continue
		}
		err := m.parseTag(f.Tag)
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(m.curName, "~") {
			continue
		}
		// 增加操作忽略空值
		if fv.IsZero() {
			if m.curMustSet {
				return nil, errors.New(fmt.Sprintf("%s must have a value", m.curName))
			}
			if m.idSet && m.idName == m.curName {
				return nil, errors.New("field with id set must not empty")
			}
			continue
		}
		if m.idSet && m.idName == m.curName {
			switch fv.Interface().(type) {
			case string:
				m.idVal = fv.Interface().(string)
			case int:
				m.idVal = fmt.Sprintf("%d", fv.Interface().(int))
			case float64:
				m.idVal = fmt.Sprintf("%f", fv.Interface().(float64))
			default:
				return nil, errors.New("unsupport id datatype")
			}
		}
		nql, err := m.setCurVal(fv)
		if err != nil {
			return nil, err
		}
		for _, fc := range m.Facets {
			if fc.PredWithLang == m.curName && fc.Seq < len(nql) {
				err = fc.Combine(nql[fc.Seq])
				if err != nil {
					return nil, err
				}
			}
		}
		setNquads = append(setNquads, nql...)
	}
	if len(setNquads) == 0 {
		return nil, errors.New("nothing to add")
	}
	if m.idSet {
		if m.idName == "" || reflect.ValueOf(m.idVal).IsZero() {
			return nil, errors.New("id is set but id value or id name is empty")
		}
		rplc := strings.NewReplacer(
			"$type", m.Dtype,
			"$name", m.idName,
			"$value", fmt.Sprintf("%v", m.idVal),
		)
		q = rplc.Replace(model)
		cond = `@if(eq(len(a),0))`
	}
	var req = &api.Request{
		Query:     q,
		Mutations: []*api.Mutation{{Cond: cond, Set: setNquads}},
	}
	return req, nil
}

func (m *mutation) MakeUpd() (*api.Request, error) {
	var (
		model    = `query{ a as var(func: type($type)) @filter(eq($name,$value) AND NOT(uid($uid)))}`
		q        string
		cond     string
		setNquad []*api.NQuad
		delNquad []*api.NQuad
	)
	if m.Subject == "" {
		return nil, errors.New("subject must not nil in update mutation")
	}
	for i := 0; i < m.Val.NumField(); i++ {
		f := m.Val.Type().Field(i)
		fv := m.Val.Field(i)
		// 忽略uid节点
		if f.Name == Uid {
			continue
		}
		if f.Name == "Friend" {
			fmt.Println(true)
		}
		err := m.parseTag(f.Tag)
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(m.curName, "~") {
			continue
		}
		if fv.IsZero() {
			if m.curMustSet {
				return nil, errors.New(fmt.Sprintf("%s must have a value", m.curName))
			}
			if m.idSet && m.idName == m.curName && m.Val.IsZero() {
				return nil, errors.New("field with id set must not empty")
			}
		}
		if m.idSet && m.idName == m.curName {
			switch fv.Interface().(type) {
			case string:
				m.idVal = fv.Interface().(string)
			case int:
				m.idVal = fmt.Sprintf("%d", fv.Interface().(int))
			case float64:
				m.idVal = fmt.Sprintf("%f", fv.Interface().(float64))
			default:
				return nil, errors.New("unsupport id datatype")
			}
		}
		delNql, err := m.delCurPred()
		if err != nil {
			return nil, err
		}
		setNql, err := m.setCurVal(fv)
		if err != nil {
			return nil, err
		}
		for _, fc := range m.Facets {
			if fc.PredWithLang == m.curName && fc.Seq < len(setNql) {
				err = fc.Combine(setNql[fc.Seq])
				if err != nil {
					return nil, err
				}
			}
		}
		delNquad = append(delNquad, delNql)
		setNquad = append(setNquad, setNql...)
	}
	if m.idSet {
		if m.idName == "" || reflect.ValueOf(m.idVal).IsZero() {
			return nil, errors.New("id is set but id value or id name is empty")
		}
		rplc := strings.NewReplacer(
			"$type", m.Dtype,
			"$name", m.idName,
			"$value", fmt.Sprintf("%v", m.idVal),
			"$uid", m.Subject,
		)
		q = rplc.Replace(model)
		cond = `@if(eq(len(a),0))`
	}
	var req = &api.Request{
		Query:     q,
		Mutations: []*api.Mutation{{Cond: cond, Set: setNquad, Del: delNquad}},
	}
	fmt.Println(IndentJson(req))
	return req, nil
}

func (m *mutation) MakeMerge() (*api.Request, error) {
	var (
		model    = `query{ a as var(func: type($type)) @filter(eq($name,$value) AND NOT(uid($uid)))}`
		q        string
		cond     string
		setNquad []*api.NQuad
	)
	if m.Subject == "" {
		return nil, errors.New("subject must not nil in update mutation")
	}
	for i := 0; i < m.Val.NumField(); i++ {
		f := m.Val.Type().Field(i)
		fv := m.Val.Field(i)
		// 忽略uid节点
		if f.Name == Uid {
			continue
		}
		if fv.IsZero() {
			continue
		}
		err := m.parseTag(f.Tag)
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(m.curName, "~") {
			continue
		}
		if m.idSet && m.idName == m.curName {
			switch fv.Interface().(type) {
			case string:
				m.idVal = fv.Interface().(string)
			case int:
				m.idVal = fmt.Sprintf("%d", fv.Interface().(int))
			case float64:
				m.idVal = fmt.Sprintf("%f", fv.Interface().(float64))
			default:
				return nil, errors.New("unsupport id datatype")
			}
		}
		setNql, err := m.setCurVal(fv)
		if err != nil {
			return nil, err
		}
		for _, fc := range m.Facets {
			if fc.PredWithLang == m.curName && fc.Seq < len(setNql) {
				err = fc.Combine(setNql[fc.Seq])
				if err != nil {
					return nil, err
				}
			}
		}
		setNquad = append(setNquad, setNql...)
	}
	if m.idSet {
		if m.idName == "" || reflect.ValueOf(m.idVal).IsZero() {
			return nil, errors.New("id is set but id value or id name is empty")
		}
		rplc := strings.NewReplacer(
			"$type", m.Dtype,
			"$name", m.idName,
			"$value", fmt.Sprintf("%v", m.idVal),
			"$uid", m.Subject,
		)
		q = rplc.Replace(model)
		cond = `@if(eq(len(a),0))`
	}
	var req = &api.Request{
		Query:     q,
		Mutations: []*api.Mutation{{Cond: cond, Set: setNquad}},
	}
	return req, nil
}

func (m *mutation) MakeDelVal() (*api.Request, error) {
	var (
		delNquad []*api.NQuad
	)
	if m.Subject == "" {
		return nil, errors.New("subject must not nil in update mutation")
	}
	for i := 0; i < m.Val.NumField(); i++ {
		f := m.Val.Type().Field(i)
		fv := m.Val.Field(i)
		// 忽略uid节点
		if f.Name == Uid {
			continue
		}
		err := m.parseTag(f.Tag)
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(m.curName, "~") {
			continue
		}
		if fv.IsZero() {
			continue
		}
		if m.curMustSet {
			return nil, errors.New(
				fmt.Sprintf("cannot delete %s with must tag", m.curName))
		}
		if m.idSet && m.idName == m.curName {
			if !m.Val.IsZero() {
				return nil, errors.New("field with id set must not delete")
			}
			switch fv.Interface().(type) {
			case string:
				m.idVal = fv.Interface().(string)
			case int:
				m.idVal = fmt.Sprintf("%d", fv.Interface().(int))
			case float64:
				m.idVal = fmt.Sprintf("%f", fv.Interface().(float64))
			default:
				return nil, errors.New("unsupport id datatype")
			}
		}
		setNql, err := m.setCurVal(fv)
		if err != nil {
			return nil, err
		}
		for _, fc := range m.Facets {
			if fc.PredWithLang == m.curName && fc.Seq < len(setNql) {
				err = fc.Combine(setNql[fc.Seq])
				if err != nil {
					return nil, err
				}
			}
		}
		delNquad = append(delNquad, setNql...)
	}
	var req = &api.Request{
		Mutations: []*api.Mutation{{Del: delNquad}},
	}
	return req, nil
}

func (m *mutation) MakeDelNode() (*api.Request, error) {
	var (
		revList  []string
		q        string
		mul      = []*api.Mutation{{Del: []*api.NQuad{{Subject: m.Subject, Predicate: StarAll, ObjectValue: starNqVal}}}}
		delModel = `query{ var(func: type($type)) {$reverse}}`
	)
	for i := 0; i < m.Val.NumField(); i++ {
		f := m.Val.Type().Field(i)
		// 忽略uid节点
		if f.Name == Uid {
			continue
		}
		err := m.parseTag(f.Tag)
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(m.curName, "~") {
			continue
		}
		revList = append(revList, m.curName)
	}
	if len(revList) > 0 {
		var rList []string
		var keyList []string
		var valList []string
		for k, v := range revList {
			key := fmt.Sprintf("a%d", k)
			valList = append(valList, strings.Trim(v, "~"))
			rList = append(rList, fmt.Sprintf("%s as %s", key, v))
			keyList = append(keyList, key)
		}
		rplc := strings.NewReplacer(
			"$type", m.Dtype,
			"$reverse", "\n"+strings.Join(rList, "\n"),
		)
		q = rplc.Replace(delModel)
		for k, v := range keyList {
			if k < len(valList) {
				mu := &api.Mutation{
					Cond: fmt.Sprintf("@if(gt(len(%s),0))", v),
					Del: []*api.NQuad{
						{
							Subject:     fmt.Sprintf("uid(%s)", v),
							Predicate:   valList[k],
							ObjectValue: starNqVal,
						},
					},
				}
				mul = append(mul, mu)
			}
		}
	}
	req := &api.Request{
		Query:     q,
		Mutations: mul,
	}
	fmt.Println(IndentJson(req))
	return req, nil
}

func (m *mutation) setCurVal(val reflect.Value) ([]*api.NQuad, error) {
	var r []*api.NQuad
	fc, ok := typeNqTypeMap[m.curDt]
	if !ok {
		return nil, errors.New("error datatype " + m.Dtype)
	}
	if val.Kind() != reflect.Slice {
		nqVal, nqobj, err := fc(val)
		if err != nil {
			return nil, err
		}
		nq := &api.NQuad{
			Subject:     m.Subject,
			Predicate:   m.curPred,
			ObjectId:    nqobj,
			ObjectValue: nqVal,
			Lang:        m.curLang,
		}
		r = append(r, nq)
	} else {
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)
			nqVal, nqobj, err := fc(v)
			if err != nil {
				return nil, err
			}
			nq := &api.NQuad{
				Subject:     m.Subject,
				Predicate:   m.curPred,
				ObjectId:    nqobj,
				ObjectValue: nqVal,
				Lang:        m.curLang,
			}
			r = append(r, nq)
		}
	}
	return r, nil
}

func (m *mutation) delCurPred() (*api.NQuad, error) {
	if m.Subject == "" || m.curPred == "" {
		return nil, errors.New("make del pred failed, subject or predicate is nil")
	}
	nq := &api.NQuad{
		Subject:     m.Subject,
		Predicate:   m.curPred,
		ObjectId:    StarAll,
		ObjectValue: starNqVal,
	}
	return nq, nil
}

func (m *mutation) delNode() (*api.NQuad, error) {
	if m.Subject == "" {
		return nil, errors.New("make del node failed, subject is nil")
	}
	nq := &api.NQuad{
		Subject:     m.Subject,
		Predicate:   StarAll,
		ObjectValue: starNqVal,
	}
	return nq, nil
}

// parseTag 解析tag相关变量到结构体参数中
func (m *mutation) parseTag(tag reflect.StructTag) error {
	m.curName = ""
	m.curDt = ""
	m.curPred = ""
	m.curLang = ""
	m.curMustSet = false
	dbtag := tag.Get(TagDb)
	dbtagList := strings.Split(dbtag, ",")
	if len(dbtagList) < 2 {
		return errors.New("tag must have predName and predType")
	}
	for k, tg := range dbtagList {
		if k == 0 {
			tgList := strings.Split(tg, "@")
			m.curName = tg
			m.curPred = tgList[0]
			if len(tgList) > 1 {
				m.curLang = tgList[2]
			}
			continue
		}
		if k == 1 {
			if _, ok := TypeAttrMap[tg]; !ok {
				return errors.New("unrecognized datatype set in tag")
			}
			m.curDt = tg
		}
		if tg == tagId {
			if m.idSet {
				return errors.New("id tag can only set once")
			}
			m.idSet = true
			m.idName = m.curName
		}
		if tg == tagMust {
			m.curMustSet = true
		}
	}
	if m.curName == "" || m.curPred == "" || m.curDt == "" {
		return errors.New("get predicate name or datatype failed in tag ")
	}
	return nil
}

func parseUidField(val reflect.Value) (string, string, error) {
	var (
		uid   string
		dtype string
	)
	if val.Kind() != reflect.Struct {
		return "", "", errors.New("need struct value")
	}
	tp := val.Type()
	fieldTp, ok := tp.FieldByName(Uid)
	if !ok {
		return "", "", errors.New("uid field not found")
	}
	fieldVal := val.FieldByName(Uid)
	if id, ok := fieldVal.Interface().(string); !fieldVal.IsZero() && ok {
		uid = id
	}
	dtype = fieldTp.Tag.Get(TagDtype)
	return uid, dtype, nil
}

func IndentJson(obj interface{}) string {
	ret, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(ret)
}
