package ast

import (
	"fmt"
	"strings"

	"mal/ast/token"
)

type AtomKind uint8

const (
	Nil AtomKind = iota
	Bool
	Int
	Float
	String
	Keyword
	Vector
	Map
)

type (
	Node interface {
		fmt.Stringer
		Pos() token.Pos
		End() token.Pos
	}

	Comment struct {
		pos     token.Pos
		Content string // No newline included
	}
	Symbol struct {
		pos     token.Pos
		Content string
	}
	AtomSingle struct { // nil true false number string keyword
		pos     token.Pos
		Kind    AtomKind
		Content string
	}
	AtomContainer struct { // vector map
		pos   token.Pos
		end   token.Pos
		Kind  AtomKind
		Elems []Node
	}
	List struct {
		end    token.Pos
		Symbol *Symbol
		Elems  []Node
	}
)

func (c *Comment) Pos() token.Pos {
	return c.pos
}

func (c *Comment) End() token.Pos {
	n := len(c.Content)
	return token.Pos{
		Offset: c.pos.Offset + n,
		Line:   c.pos.Line,
		Column: c.pos.Column + n,
	}
}

func (c *Comment) String() string {
	return c.Content
}

func (s *Symbol) Pos() token.Pos {
	return s.pos
}

func (s *Symbol) End() token.Pos {
	n := len(s.Content)
	return token.Pos{
		Offset: s.pos.Offset + n,
		Line:   s.pos.Line,
		Column: s.pos.Column + n,
	}
}

func (s *Symbol) String() string {
	return s.Content
}

func (a *AtomSingle) Pos() token.Pos {
	return a.pos
}

func (a *AtomSingle) End() token.Pos {
	n := len(a.Content)
	return token.Pos{
		Offset: a.pos.Offset + n,
		Line:   a.pos.Line,
		Column: a.pos.Column + n,
	}
}

func (a *AtomSingle) String() string {
	if a.Kind == String {
		return fmt.Sprintf("%q", a.Content[1:len(a.Content)-1])
	}
	return a.Content
}

func (a *AtomContainer) Pos() token.Pos {
	return a.pos
}

func (a *AtomContainer) End() token.Pos {
	return a.end
}

func (a *AtomContainer) String() string {
	n := len(a.Elems)
	elems := make([]string, n+2)
	switch a.Kind {
	case Vector:
		elems[0] = "["
		elems[n+1] = "]"
	case Map:
		elems[0] = "{"
		elems[n+1] = "}"
	}
	for i, elem := range a.Elems {
		elems[i+1] = elem.String()
	}
	return fmt.Sprintf("%s%s%s", elems[0], strings.Join(elems[1:n+1], " "), elems[n+1])
}

func (l *List) Pos() token.Pos {
	return l.Symbol.pos
}

func (l *List) End() token.Pos {
	return l.end
}

func (l *List) String() string {
	n := len(l.Elems)
	elems := make([]string, n)
	for i, elem := range l.Elems {
		elems[i] = elem.String()
	}
	macro := l.Symbol.Content
	if macro != "" {
		macro += " "
	}
	return fmt.Sprintf("(%s%s)", macro, strings.Join(elems, " "))
}
