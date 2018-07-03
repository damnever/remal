package mal

import (
	"mal/ast"
	"mal/types"
)

type Evaler struct{}

func (e *Evaler) EvalAST(a *ast.AST) []types.Valuer {
	vs := []types.Valuer{}
	a.Walk(func(node ast.Node) bool {
		if v := e.evalNode(node); v != nil {
			vs = append(vs, v)
		}
		return true
	})
	return vs
}

func (e *Evaler) evalNode(node ast.Node) types.Valuer {
	switch x := node.(type) {
	case *ast.Comment:
		return nil
	case *ast.Symbol:
		return types.NewRaw(x)
	case *ast.AtomSingle:
		return e.evalAtomSingle(x)
	case *ast.AtomContainer:
		return e.evalAtomContainer(x)
	case *ast.List:
		return e.evalList(x)
	default:
	}
	return types.NewRaw(node)
}

func (e *Evaler) evalList(l *ast.List) types.Valuer {
	nelems := len(l.Elems)
	if nelems != 3 {
		return types.NewRaw(l)
	}
	symbol, ok := l.Elems[0].(*ast.Symbol)
	if !ok {
		return types.NewRaw(l)
	}
	switch symbol.Content {
	case "+", "-", "*", "/":
	default:
		return types.NewRaw(l)
	}
	nums := make([]types.Number, nelems)
	for i, elem := range l.Elems[1:] {
		num, ok := e.evalNode(elem).(types.Number)
		if !ok {
			return types.NewRaw(l)
		}
		nums[i] = num
	}
	switch symbol.Content {
	case "+":
		return nums[0].Add(nums[1])
	case "-":
		return nums[0].Sub(nums[1])
	case "*":
		return nums[0].Mul(nums[1])
	case "/":
		return nums[0].Div(nums[1])
	}
	panic("how dare you")
}

func (e *Evaler) evalAtomSingle(as *ast.AtomSingle) types.Valuer {
	switch as.Kind {
	case ast.Nil:
		return types.Nil{}
	case ast.Bool:
		return types.NewBool(as.Content)
	case ast.Int:
		return types.NewInt(as.Content)
	case ast.Float:
		return types.NewFloat(as.Content)
	case ast.String:
		return types.String(as.Content[1 : len(as.Content)-1])
	case ast.Keyword:
		return types.Keyword(as.Content[1:])
	}
	panic("how dare you")
}

func (e *Evaler) evalAtomContainer(ac *ast.AtomContainer) types.Valuer {
	switch ac.Kind {
	case ast.Vector:
		vec := &types.Vector{}
		for _, elem := range ac.Elems {
			if v := e.evalNode(elem); v != nil {
				vec.Append(v)
			}
		}
		return vec
	case ast.Map:
		m := types.Map{}
		var k types.Valuer
		for _, elem := range ac.Elems {
			if k == nil {
				if vv := e.evalNode(elem); vv != nil {
					k = vv.(types.MapKey)
				}
			} else {
				if vv := e.evalNode(elem); vv != nil {
					m[k] = vv
					k = nil
				}
			}
		}
		return m
	}
	panic("how dare you")
}
