package mal

import (
	"fmt"
	"mal/types"
)

type Env struct {
	outer *Env
	data  map[string]types.Valuer
}

func NewEnv(outer *Env) *Env {
	return &Env{
		outer: outer,
		data:  map[string]types.Valuer{},
	}
}

func (e *Env) Set(symbol string, value types.Valuer) {
	e.data[symbol] = value
}

func (e *Env) Get(symbol string) (types.Valuer, error) {
	v, ok := e.Find(symbol)
	if !ok {
		return nil, fmt.Errorf("symbol(%s) not found", symbol)
	}
	return v, nil
}

func (e *Env) Find(symbol string) (types.Valuer, bool) {
	for env := e; env != nil; env = env.outer {
		if v, ok := env.data[symbol]; ok {
			return v, ok
		}
	}
	return nil, false
}
