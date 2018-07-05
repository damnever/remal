package types

import (
	"container/list"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type (
	Valuer interface {
		SPrint(readable bool) string
		IsEqaulTo(oth Valuer) bool
	}
	Number interface {
		Valuer
		Add(Number) Number
		Sub(Number) Number
		Mul(Number) Number
		Div(Number) Number
		Compare(Number) int
	}
	MapKey interface {
		Valuer
		Key()
	}

	FuncType func(...Valuer) (Valuer, error)
)

type (
	Raw     struct{ x fmt.Stringer }
	Nil     struct{}
	Bool    bool
	Int     int64
	Float   float64
	String  string
	Keyword string
	List    struct{ *list.List }
	Vector  []Valuer
	Map     map[Valuer]Valuer
	Func    struct {
		name string
		Exec FuncType
	}
)

func NewRaw(x fmt.Stringer) Raw {
	return Raw{x: x}
}

func (r Raw) IsEqaulTo(oth Valuer) bool {
	o, ok := oth.(Raw)
	return ok && r == o
}

func (r Raw) SPrint(readable bool) string {
	return r.x.String()
}

func (n Nil) IsEqaulTo(oth Valuer) bool {
	_, ok := oth.(Nil)
	return ok
}

func (n Nil) SPrint(readable bool) string {
	return "nil"
}

func NewInt(v string) Int {
	x, _ := strconv.ParseInt(v, 10, 64)
	return Int(x)
}

func (i Int) Add(oth Number) Number {
	switch x := oth.(type) {
	case Int:
		return i + x
	case Float:
		return Float(i) + x
	}
	panic("try again?")
}

func (i Int) Sub(oth Number) Number {
	switch x := oth.(type) {
	case Int:
		return i - x
	case Float:
		return Float(i) - x
	}
	panic("try again?")
}

func (i Int) Mul(oth Number) Number {
	switch x := oth.(type) {
	case Int:
		return i * x
	case Float:
		return Float(i) * x
	}
	panic("try again?")
}

func (i Int) Div(oth Number) Number {
	// FIXME: divide by zero
	switch x := oth.(type) {
	case Int:
		return i / x
	case Float:
		return Float(i) / x
	}
	panic("try again?")
}

func (i Int) Compare(oth Number) int {
	switch x := oth.(type) {
	case Int:
		return int(i - x)
	case Float:
		return int(math.Ceil(float64(Float(i) - x)))
	}
	panic("try again?")
}

func (i Int) IsEqaulTo(oth Valuer) bool {
	o, ok := oth.(Int)
	return ok && i == o
}

func (i Int) SPrint(readable bool) string {
	return fmt.Sprintf("%d", i)
}

func NewFloat(v string) Float {
	x, _ := strconv.ParseFloat(v, 64)
	return Float(x)
}

func (f Float) Add(oth Number) Number {
	switch x := oth.(type) {
	case Int:
		return f + Float(x)
	case Float:
		return f + x
	}
	panic("try again?")
}

func (f Float) Sub(oth Number) Number {
	switch x := oth.(type) {
	case Int:
		return f - Float(x)
	case Float:
		return f - x
	}
	panic("try again?")
}

func (f Float) Mul(oth Number) Number {
	switch x := oth.(type) {
	case Int:
		return f * Float(x)
	case Float:
		return f * x
	}
	panic("try again?")
}

func (f Float) Div(oth Number) Number {
	// FIXME: divide by zero
	switch x := oth.(type) {
	case Int:
		return f / Float(x)
	case Float:
		return f / x
	}
	panic("try again?")
}

func (f Float) Compare(oth Number) int {
	switch x := oth.(type) {
	case Int:
		return int(math.Ceil(float64(f - Float(x))))
	case Float:
		return int(math.Ceil(float64(f - x)))
	}
	panic("try again?")
}

func (f Float) IsEqaulTo(oth Valuer) bool {
	o, ok := oth.(Float)
	return ok && f == o
}

func (f Float) SPrint(readable bool) string {
	return fmt.Sprintf("%f", f)
}

func NewBool(v string) Bool {
	if v == "true" {
		return true
	}
	return false
}

func (b Bool) IsEqaulTo(oth Valuer) bool {
	o, ok := oth.(Bool)
	return ok && b == o
}

func (b Bool) SPrint(readable bool) string {
	if b {
		return "true"
	}
	return "false"
}

func (s String) IsEqaulTo(oth Valuer) bool {
	o, ok := oth.(String)
	return ok && s == o
}

func (s String) SPrint(readable bool) string {
	if readable {
		return fmt.Sprintf("\"%s\"", string(s))
	}
	return string(s)
}

func (String) Key() {}

func (k Keyword) IsEqaulTo(oth Valuer) bool {
	o, ok := oth.(Keyword)
	return ok && k == o
}

func (k Keyword) SPrint(readable bool) string {
	return fmt.Sprintf(":%s", string(k))
}

func (Keyword) Key() {}

func NewList() List {
	return List{List: list.New()}
}

func (l List) Append(vs ...Valuer) {
	for _, v := range vs {
		l.PushBack(v)
	}
}

func (l List) ToVector() *Vector {
	vec := &Vector{}
	for e := l.Front(); e != nil; e = e.Next() {
		vec.Append(e.Value.(Valuer))
	}
	return vec
}

func (l List) IsEqaulTo(oth Valuer) bool {
	var o List
	switch x := oth.(type) {
	case *Vector:
		return l.ToVector().IsEqaulTo(x)
	case List:
		o = x
	default:
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

func (l List) SPrint(readable bool) string {
	elems := []string{"("}
	for i, elem := 0, l.Front(); elem != nil; i, elem = i+1, elem.Next() {
		s := elem.Value.(Valuer).SPrint(readable)
		if i == 0 {
			elems = append(elems, s)
		} else {
			elems = append(elems, " "+s)
		}
	}
	elems = append(elems, ")")
	return strings.Join(elems, "")
}

func (v *Vector) Append(elems ...Valuer) {
	*v = append(*v, elems...)
}

func (v *Vector) Remove(elem Valuer) {
	for i, x := range *v {
		if x.IsEqaulTo(elem) {
			copy((*v)[i:], (*v)[i+1:])
			n := len(*v) - 1
			(*v)[n] = nil
			*v = (*v)[:n]
			break
		}
	}
}

func (v *Vector) Len() int {
	return len(*v)
}

func (v *Vector) ToList() List {
	l := NewList()
	l.Append((*v)...)
	return l
}

func (v *Vector) IsEqaulTo(oth Valuer) bool {
	var o *Vector
	switch x := oth.(type) {
	case List:
		return v.ToList().IsEqaulTo(x)
	case *Vector:
		o = x
	default:
		return false
	}
	o, ok := oth.(*Vector)
	if !ok {
		return false
	}

	n := len(*v)
	if n != len(*o) {
		return false
	}
	for i := 0; i < n; i++ {
		if !(*v)[i].IsEqaulTo((*o)[i]) {
			return false
		}
	}
	return true
}

func (v *Vector) SPrint(readable bool) string {
	elems := []string{"["}
	for i, elem := range *v {
		s := elem.SPrint(readable)
		if i == 0 {
			elems = append(elems, s)
		} else {
			elems = append(elems, " "+s)
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

func (m Map) SPrint(readable bool) string {
	elems := []string{"{"}
	i := 0
	for k, v := range m {
		sk := k.SPrint(readable)
		sv := v.SPrint(readable)
		if i == 0 {
			elems = append(elems, fmt.Sprintf("%s %s", sk, sv))
		} else {
			elems = append(elems, fmt.Sprintf(" %s %s", sk, sv))
		}
		i++
	}
	elems = append(elems, "}")
	return strings.Join(elems, "")
}

func NewFunc(name string, fn FuncType) Func {
	return Func{name: name, Exec: fn}
}

func (f Func) IsLambda() bool {
	return f.name == ""
}

func (f Func) IsEqaulTo(Valuer) bool {
	return false
}

func (f Func) SPrint(readable bool) string {
	return "#<function>"
}

func escape(s string) string {
	s = strings.Replace(s, "\\", "\\\\", -1)
	return s
}
