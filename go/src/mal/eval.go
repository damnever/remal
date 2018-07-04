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
	env.Set("+", types.NewFunc("+", funcAdd))
	env.Set("-", types.NewFunc("-", funcSub))
	env.Set("*", types.NewFunc("*", funcMul))
	env.Set("/", types.NewFunc("/", funcDiv))
	return &Evaler{env: env}
}

func (e *Evaler) EvalAST(a *ast.AST) (vs []types.Valuer, err error) {
	a.Walk(func(node ast.Node) bool {
		var v types.Valuer
		v, err = e.evalNode(node)
		if err == errIgnore {
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
		env, err := e.env.Get(x.Content)
		if err != nil {
			return nil, fmt.Errorf("[%s] %v", x.Pos(), err)
		}
		return env, nil
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

func (e *Evaler) evalList(l *ast.List) (types.Valuer, error) {
	symbol, ok := l.Elems[0].(*ast.Symbol)
	if !ok {
		return types.NewRaw(l), nil
	}

	// FIXME(damnever): fuck..
	switch symbol.Content {
	case "def!":
		v, err := e.evalNode(l.Elems[2])
		if err != nil {
			return nil, err
		}
		e.env.Set(l.Elems[1].(*ast.Symbol).Content, v)
		return v, nil
	case "let*":
		env := NewEnv(e.env)
		var elems []ast.Node
		switch x := l.Elems[1].(type) {
		case *ast.List:
			elems = x.Elems
		case *ast.AtomContainer:
			elems = x.Elems
		}
		for i := 0; i < len(elems); i = i + 2 {
			v, _ := NewEvaler(env).evalNode(elems[i+1])
			env.Set(elems[i].(*ast.Symbol).Content, v)
		}
		return NewEvaler(env).evalNode(l.Elems[2])
	default:
	}
	ev, err := e.env.Get(symbol.Content)
	if err != nil {
		return nil, fmt.Errorf("[%s] %v", symbol.Pos(), err)
	}
	if fn, ok := ev.(types.Func); ok {
		return e.evalFunc(fn, l.Elems[1:])
	}
	return ev, nil
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
