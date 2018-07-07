package mal

import (
	"errors"
	"fmt"

	"mal/ast"
	"mal/types"
)

var (
	errIgnore = errors.New("ignore")
)

type Evaler struct {
	env *Env
}

func NewEvaler(env *Env) *Evaler {
	for k, v := range funcmap {
		env.Set(k, types.NewFunc(k, v))
	}
	return &Evaler{env: env}
}

func (e *Evaler) EvalAST(a *ast.AST) (vs []types.Valuer, err error) {
	a.Walk(func(node ast.Node) bool {
		var v types.Valuer
		v, err = e.evalNode(node)
		if err == errIgnore {
			err = nil
			return true
		}
		if err != nil {
			return false
		}
		vs = append(vs, v)
		return true
	})
	return
}

func (e *Evaler) evalNode(node ast.Node) (types.Valuer, error) {
	switch x := node.(type) {
	case *ast.Comment:
		return nil, errIgnore
	case *ast.Symbol:
		return e.evalSymbol(x)
	case *ast.AtomSingle:
		return e.evalAtomSingle(x), nil
	case *ast.AtomContainer:
		return e.evalAtomContainer(x)
	case *ast.List:
		return e.evalList(x)
	default:
	}
	return types.NewRaw(node), nil
}

func (e *Evaler) evalSymbol(symbol *ast.Symbol) (types.Valuer, error) {
	env, err := e.env.Get(symbol.Content)
	if err != nil {
		return nil, fmt.Errorf("[%s] %v", symbol.Pos(), err)
	}
	return env, nil
}

func (e *Evaler) evalList(list *ast.List) (types.Valuer, error) {
	// FIXME(damnever): fuck..
	evaler := e
	var n ast.Node = list

	for {
		l, ok := n.(*ast.List)
		if !ok {
			return evaler.evalNode(n)
		}
		if len(l.Elems) == 0 {
			return types.NewRaw(l), nil
		}

		symbol, ok := l.Elems[0].(*ast.Symbol)
		if !ok { // In place lambda call
			v, err := e.evalNode(l.Elems[0])
			if err != nil {
				return types.NewRaw(l), nil
			}
			if fn, ok := v.(types.LambdaFunc); ok {
				evaler, n, err = evaler.evalLambaFunc(fn, l.Elems[1:])
				if err != nil {
					return nil, err
				}
				continue
			}
			return types.NewRaw(l), nil
		}

		switch symbol.Content {
		case "def!":
			v, err := evaler.evalNode(l.Elems[2])
			if err != nil {
				return nil, err
			}
			evaler.env.Set(l.Elems[1].(*ast.Symbol).Content, v)
			return v, nil

		case "let*":
			var elems []ast.Node
			switch x := l.Elems[1].(type) {
			case *ast.List:
				elems = x.Elems
			case *ast.AtomContainer:
				if x.Kind != ast.Vector {
					return nil, fmt.Errorf("[%s] expect list or vector, got map", x.Pos())
				}
				elems = x.Elems
			}
			letenv := NewEnv(evaler.env, nil, nil)
			for i := 0; i < len(elems); i = i + 2 {
				v, _ := NewEvaler(letenv).evalNode(elems[i+1])
				letenv.Set(elems[i].(*ast.Symbol).Content, v)
			}

			evaler = NewEvaler(letenv)
			n = l.Elems[2]

		case "do":
			last := len(l.Elems) - 1
			for _, elem := range l.Elems[1:last] {
				if _, err := evaler.evalNode(elem); err != nil {
					if err == errIgnore {
						continue
					}
					return nil, err
				}
			}
			n = l.Elems[last]

		case "if":
			v1, err := evaler.evalNode(l.Elems[1])
			if err != nil && err != errIgnore {
				return nil, err
			}
			ok := true
			switch x := v1.(type) {
			case types.Nil:
				ok = false
			case types.Bool:
				ok = bool(x)
			default:
			}

			if !ok {
				if len(l.Elems) < 4 {
					return types.Nil{}, nil
				}
				n = l.Elems[3]
			} else {
				n = l.Elems[2]
			}

		case "fn*":
			var elems []ast.Node
			switch x := l.Elems[1].(type) {
			case *ast.List:
				elems = x.Elems
			case *ast.AtomContainer:
				if x.Kind != ast.Vector {
					return nil, fmt.Errorf("[%s] expect list or vector, got map", x.Pos())
				}
				elems = x.Elems
			}
			binds := []string{}
			for _, elem := range elems {
				binds = append(binds, elem.(*ast.Symbol).Content)
			}
			return types.NewLambdaFunc(evaler.env, l.Elems[2], binds), nil

		default:
			ev, err := evaler.evalSymbol(symbol)
			if err != nil {
				return nil, err
			}
			switch fn := ev.(type) {
			case types.Func:
				return evaler.evalFunc(fn, l.Elems[1:])
			case types.LambdaFunc:
				var err error
				evaler, n, err = evaler.evalLambaFunc(fn, l.Elems[1:])
				if err != nil {
					return nil, err
				}
			default:
				return ev, nil
			}
		}
	}
}

func (e *Evaler) evalLambaFunc(fn types.LambdaFunc, nodes []ast.Node) (evaler *Evaler, n ast.Node, err error) {
	exprs := []types.Valuer{}
	var ev types.Valuer
	for _, nn := range nodes {
		ev, err = e.evalNode(nn)
		if err != nil {
			if err == errIgnore {
				err = nil
				continue
			}
			return
		}
		exprs = append(exprs, ev)
	}

	evaler = NewEvaler(NewEnv(fn.Env.(*Env), fn.Binds, exprs))
	n = fn.Expr.(ast.Node)
	return
}

func (e *Evaler) evalFunc(fn types.Func, nodes []ast.Node) (types.Valuer, error) {
	args := make([]types.Valuer, len(nodes))
	for i, node := range nodes {
		v, err := e.evalNode(node)
		if err != nil {
			if err == errIgnore {
				continue
			}
			return nil, err
		}
		args[i] = v
	}
	return fn.Exec(args...)
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

func (e *Evaler) evalAtomContainer(ac *ast.AtomContainer) (types.Valuer, error) {
	switch ac.Kind {
	case ast.Vector:
		vec := &types.Vector{}
		for _, elem := range ac.Elems {
			v, err := e.evalNode(elem)
			if err != nil {
				if err == errIgnore {
					continue
				}
				return nil, err
			}
			vec.Append(v)
		}
		return vec, nil
	case ast.Map:
		m := types.Map{}
		var k types.Valuer
		for _, elem := range ac.Elems {
			v, err := e.evalNode(elem)
			if err != nil {
				if err == errIgnore {
					continue
				}
				return nil, err
			}
			if k == nil {
				k = v.(types.MapKey)
			} else {
				m[k] = v
				k = nil
			}
		}
		if k != nil {
			return nil, fmt.Errorf("[%s] key/value pair required", ac.End())
		}
		return m, nil
	}
	panic("how dare you")
}
