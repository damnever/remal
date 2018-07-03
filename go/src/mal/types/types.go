package types

import (
	"container/list"
	"fmt"
	"strings"
)

type Valuer interface {
	String() string
	IsEqaulTo(oth Valuer) bool
}

type (
	Nil     struct{}
	Bool    bool
	Int     int64
	Float   float64
	String  string
	Keyword string
	List    struct{ *list.List }
	Vector  struct{ elems []Valuer }
	Map     map[Valuer]Valuer
)

func (n Nil) IsEqaulTo(oth Valuer) bool {
	_, ok := oth.(Nil)
	return ok
}

func (n Nil) String() string {
	return "nil"
}

func (i Int) IsEqaulTo(oth Valuer) bool {
	o, ok := oth.(Int)
	return ok && i == o
}

func (i Int) String() string {
	return fmt.Sprintf("%d", i)
}

func (f Float) IsEqaulTo(oth Valuer) bool {
	o, ok := oth.(Float)
	return ok && f == o
}

func (f Float) String() string {
	return fmt.Sprintf("%d", f)
}

func (b Bool) IsEqaulTo(oth Valuer) bool {
	o, ok := oth.(Bool)
	return ok && b == o
}

func (b Bool) String() string {
	if b {
		return "true"
	}
	return "false"
}

func (s String) IsEqaulTo(oth Valuer) bool {
	o, ok := oth.(String)
	return ok && s == o
}

func (s String) String() string {
	return fmt.Sprintf("\"%s\"", string(s))
}

func (k Keyword) IsEqaulTo(oth Valuer) bool {
	o, ok := oth.(Keyword)
	return ok && k == o
}

func (k Keyword) String() string {
	return fmt.Sprintf(":\"%s\"", string(k))
}

func (l List) IsEqaulTo(oth Valuer) bool {
	o, ok := oth.(List)
	if !ok {
		return false
	}
	if l.Len() != o.Len() {
		return false
	}
	for e1, e2 := l.Front(), o.Front(); ; e1, e2 = e1.Next(), e2.Next() {
		if e1 == nil && e2 == nil {
			break
		}
		if e1 == nil && e2 != nil {
			return false
		}
		if e1 != nil && e2 == nil {
			return false
		}
		if !e1.Value.(Valuer).IsEqaulTo(e2.Value.(Valuer)) {
			return false
		}
	}
	return true
}

func (l List) String() string {
	elems := []string{"("}
	for i, elem := 0, l.Front(); elem != nil; i, elem = i+1, elem.Next() {
		if i == 0 {
			elems = append(elems, fmt.Sprintf("%s", elem.Value))
		} else {
			elems = append(elems, fmt.Sprintf(" %s", elem.Value))
		}
	}
	elems = append(elems, ")")
	return strings.Join(elems, "")
}

func (v Vector) Append(elems ...Valuer) {
	v.elems = append(v.elems, elems...)
}

func (v Vector) Remove(elem Valuer) {
	for i, x := range v.elems {
		if x.IsEqaulTo(elem) {
			copy(v.elems[i:], v.elems[i+1:])
			n := len(v.elems) - 1
			v.elems[n] = nil
			v.elems = v.elems[:n]
			break
		}
	}
}

func (v Vector) IsEqaulTo(oth Valuer) bool {
	o, ok := oth.(Vector)
	if !ok {
		return false
	}

	n := len(v.elems)
	if n != len(o.elems) {
		return false
	}
	for i := 0; i < n; i++ {
		if !v.elems[i].IsEqaulTo(o.elems[i]) {
			return false
		}
	}
	return true
}

func (v Vector) String() string {
	elems := []string{"["}
	for i, elem := range v.elems {
		if i == 0 {
			elems = append(elems, fmt.Sprintf("%s", elem))
		} else {
			elems = append(elems, fmt.Sprintf(" %s", elem))
		}
	}
	elems = append(elems, "]")
	return strings.Join(elems, "")
}

func (m Map) IsEqaulTo(oth Valuer) bool {
	o, ok := oth.(Map)
	if !ok {
		return false
	}

	if len(m) != len(o) {
		return false
	}
	for k, v := range m {
		if !v.IsEqaulTo(o[k]) {
			return false
		}
	}
	return true
}

func (m Map) String() string {
	elems := []string{"{"}
	i := 0
	for k, v := range m {
		if i == 0 {
			elems = append(elems, fmt.Sprintf("%s %s", k, v))
		} else {
			elems = append(elems, fmt.Sprintf(" %s %s", k, v))
		}
		i++
	}
	elems = append(elems, "}")
	return strings.Join(elems, "")
}
