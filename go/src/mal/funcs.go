package mal

import "mal/types"

// FIXME(damnever): check arg number/type???????

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
