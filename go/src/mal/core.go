package mal

import (
	"fmt"
	"mal/types"
	"strings"
)

// FIXME(damnever): check arg number/type???????

var funcmap = map[string]types.FuncType{
	"+":       funcAdd,
	"-":       funcSub,
	"*":       funcMul,
	"/":       funcDiv,
	"<":       funcLess,
	"<=":      funcLessOrEqual,
	">":       funcGreater,
	">=":      funcGreaterOrEqual,
	"=":       funcIsEqual,
	"list":    funcToList,
	"list?":   funcIsList,
	"empty?":  funcIsEmpty,
	"count":   funcCount,
	"prn":     funcPrint,
	"pr-str":  funcPrintStr,
	"str":     funcStr,
	"println": funcPrintln,
}

func funcAdd(vs ...types.Valuer) (types.Valuer, error) {
	return vs[0].(types.Number).Add(vs[1].(types.Number)), nil
}

func funcSub(vs ...types.Valuer) (types.Valuer, error) {
	return vs[0].(types.Number).Sub(vs[1].(types.Number)), nil
}

func funcMul(vs ...types.Valuer) (types.Valuer, error) {
	return vs[0].(types.Number).Mul(vs[1].(types.Number)), nil
}

func funcDiv(vs ...types.Valuer) (types.Valuer, error) {
	return vs[0].(types.Number).Div(vs[1].(types.Number)), nil
}

func funcLess(vs ...types.Valuer) (types.Valuer, error) {
	r := vs[0].(types.Number).Compare(vs[1].(types.Number)) < 0
	return types.Bool(r), nil
}

func funcLessOrEqual(vs ...types.Valuer) (types.Valuer, error) {
	r := vs[0].(types.Number).Compare(vs[1].(types.Number)) <= 0
	return types.Bool(r), nil
}

func funcGreater(vs ...types.Valuer) (types.Valuer, error) {
	r := vs[0].(types.Number).Compare(vs[1].(types.Number)) > 0
	return types.Bool(r), nil
}

func funcGreaterOrEqual(vs ...types.Valuer) (types.Valuer, error) {
	r := vs[0].(types.Number).Compare(vs[1].(types.Number)) >= 0
	return types.Bool(r), nil
}

func funcPrint(vs ...types.Valuer) (types.Valuer, error) {
	s, err := funcprint(" ", true, vs...)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", s)
	return types.Nil{}, nil
}

func funcPrintStr(vs ...types.Valuer) (types.Valuer, error) {
	return funcprint(" ", true, vs...)
}

func funcStr(vs ...types.Valuer) (types.Valuer, error) {
	return funcprint("", false, vs...)
}

func funcPrintln(vs ...types.Valuer) (types.Valuer, error) {
	s, err := funcprint(" ", false, vs...)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", s)
	return types.Nil{}, nil
}

func funcprint(sep string, readable bool, vs ...types.Valuer) (types.Valuer, error) {
	ss := []string{}
	for _, v := range vs {
		ss = append(ss, v.SPrint(readable))
	}
	return types.String(strings.Join(ss, sep)), nil
}

func funcToList(vs ...types.Valuer) (types.Valuer, error) {
	l := types.NewList()
	for _, v := range vs {
		l.PushBack(v)
	}
	return l, nil
}

func funcIsList(vs ...types.Valuer) (types.Valuer, error) {
	_, ok := vs[0].(types.List)
	return types.Bool(ok), nil
}

func funcIsEmpty(vs ...types.Valuer) (types.Valuer, error) {
	switch x := vs[0].(type) {
	case types.List:
		return types.Bool(x.Len() == 0), nil
	case *types.Vector:
		return types.Bool(x.Len() == 0), nil
	default:
	}
	return types.Bool(false), nil
}

func funcCount(vs ...types.Valuer) (types.Valuer, error) {
	switch x := vs[0].(type) {
	case types.List:
		return types.Int(x.Len()), nil
	case *types.Vector:
		return types.Int(x.Len()), nil
	default:
	}
	return types.Int(0), nil
}

func funcIsEqual(vs ...types.Valuer) (types.Valuer, error) {
	return types.Bool(vs[0].IsEqaulTo(vs[1])), nil
}
