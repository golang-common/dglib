/**
 * @Author: daipengyuan
 * @Description: 谓词
 * @File:  pred
 * @Version: 1.0.0
 * @Date: 2021/8/11 12:18
 */

package dql

import (
	"errors"
	"fmt"
	"strings"
)

// Pred 谓词,也叫有向边
type Pred struct {
	Name  string
	Type  DType
	Index []DIndex
}

// Schema 转换为rdf格式
func (p Pred) Schema() string {
	var (
		model      = `$name: $type $indices .`
		name       = p.Name
		indexList  []string
		oIndexList []string

		index, count, upsert, reverse, lang bool
	)

	for _, idx := range p.Index {
		if idx == IList {
			name = fmt.Sprintf("[%s]", name)
			continue
		}
		if idx == IIndex {
			index = true
			continue
		}
		if idx == ICount {
			count = true
			continue
		}
		if idx == IUpsert {
			upsert = true
			continue
		}
		if idx == IReverse {
			reverse = true
			continue
		}
		if idx == ILang {
			lang = true
			continue
		}
		index = true
		indexList = append(indexList, string(idx))
	}
	if index && len(indexList) > 0 {
		oIndexList = append(oIndexList, fmt.Sprintf("@index(%s)", strings.Join(indexList, ",")))
	}
	if count {
		oIndexList = append(oIndexList, string(ICount))
	}
	if upsert {
		oIndexList = append(oIndexList, string(IUpsert))
	}
	if reverse {
		oIndexList = append(oIndexList, string(IReverse))
	}
	if lang {
		oIndexList = append(oIndexList, string(ILang))
	}
	replacer := strings.NewReplacer(
		"$name", p.Name,
		"$type", string(p.Type),
		"$indices", strings.Join(oIndexList, " "),
	)
	return replacer.Replace(model)
}

func (p Pred) Unmarshal(s string) error {
	var (
		name     string
		tp       DType
		indices  []DIndex
		errFmt   = errors.New("error predicate format")
		errIndex = errors.New("error predicate index")
	)
	s = strings.Replace(s, ", ", ",", -1)
	for k, v := range strings.Split(s, " ") {
		if k == 0 {
			name = strings.Trim(v, " <>:")
			continue
		}
		if k == 1 {
			if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
				indices = append(indices, IList)
				v = v[1 : len(v)-1]
			}
			_, ok := DTypeMap[DType(v)]
			if !ok {
				return errFmt
			}
			tp = DType(v)
			continue
		}
		if strings.HasPrefix(v, "@index") {
			istr := v[7 : len(v)-1]
			indices = append(indices, IIndex)
			for _, vi := range strings.Split(istr, ",") {
				_, ok := IIndexMap[DIndex(vi)]
				if !ok {
					return errIndex
				}
				indices = append(indices, DIndex(v))
			}
			continue
		}
		if strings.HasPrefix(v, "@") {
			_, ok := IGIndexMap[DIndex(v)]
			if !ok {
				return errIndex
			}
			indices = append(indices, DIndex(v))
			continue
		}
	}
	p.Name = name
	p.Type = tp
	p.Index = indices
	return nil
}
