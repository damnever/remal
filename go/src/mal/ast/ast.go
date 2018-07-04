package ast

import (
	"fmt"

	"mal/ast/token"
)

type AST struct {
	tr    *tokenReader
	nodes []Node
}

func (ast *AST) Parse(code string) error {
	ast.tr = newTokenReader(code)
	ast.nodes = []Node{}

	for {
		node, err := ast.processForm()
		if err != nil {
			return err
		}
		if node == nil {
			break
		}
		ast.nodes = append(ast.nodes, node)
	}
	return nil
}

func (ast *AST) Walk(visitor func(node Node) bool) {
	for _, node := range ast.nodes {
		if !visitor(node) {
			break
		}
	}
}

func (ast *AST) processForm() (Node, error) {
	t, err := ast.tr.Peek()
	if err != nil {
		return nil, err
	}
	var node Node

	switch t.Token {
	case token.EOF:
		return nil, nil
	case token.ILLEGAL:
		err = fmt.Errorf("[%s] illegal syntax: %s", t.Pos, t.Content)
	case token.COMMENT: // ;
		node, err = ast.processComment()
	case token.TILDEAT: // ~@
		node, err = ast.processMacro("splice-unquote", t.Pos)
	case token.SINGLEQUOTE: // '
		node, err = ast.processMacro("quote", t.Pos)
	case token.BACKQUOTE: // `
		node, err = ast.processMacro("quasiquote", t.Pos)
	case token.TILDE: // ~
		node, err = ast.processMacro("unquote", t.Pos)
	case token.CIRCUMFLEX: // ^
		node, err = ast.processMacro("with-meta", t.Pos)
	case token.ATSIGN: // @
		node, err = ast.processMacro("deref", t.Pos)
	case token.LPAREN: // (
		node, err = ast.processList()
	default:
		node, err = ast.processAtom()
	}
	return node, err
}

func (ast *AST) processComment() (Node, error) {
	t, err := ast.tr.Next()
	if err != nil {
		return nil, err
	}
	return &Comment{
		pos:     t.Pos,
		Content: t.Content,
	}, nil
}

func (ast *AST) processMacro(macro string, pos token.Pos) (Node, error) {
	ast.tr.Next()
	node := &List{
		Symbol: &Symbol{
			pos:     pos,
			Content: macro,
		},
		Elems: make([]Node, 1),
	}
	if macro == "with-meta" {
		meta, err := ast.processForm()
		if err != nil {
			return nil, err
		}
		node.Elems = append(node.Elems, meta)
	}
	n, err := ast.processForm()
	if err != nil {
		return nil, err
	}
	node.Elems[0] = n
	return node, nil
}

func (ast *AST) processList() (Node, error) {
	t, err := ast.tr.Next()
	if err != nil {
		return nil, err
	}
	if t.Token != token.LPAREN {
		return nil, fmt.Errorf("[%s] expect '(', got '%s'", t.Pos, t.Content)
	}

	node := &List{
		Symbol: &Symbol{pos: t.Pos},
		Elems:  []Node{},
	}
	t, err = ast.processContainer(&(node.Elems), token.RPAREN)
	if err != nil {
		return nil, err
	}
	node.end = t.End
	return node, err
}

func (ast *AST) processAtom() (Node, error) {
	t, err := ast.tr.Next()
	if err != nil {
		return nil, err
	}

	var node Node
	switch t.Token {
	case token.LBRACK: // [
		node, err = ast.processAtomContainer(Vector, t)
	case token.LBRACE: // {
		node, err = ast.processAtomContainer(Map, t)
	case token.NIL:
		node = ast.processKindAtomSingle(Nil, t)
	case token.BOOL:
		node = ast.processKindAtomSingle(Bool, t)
	case token.INT:
		node = ast.processKindAtomSingle(Int, t)
	case token.FLOAT:
		node = ast.processKindAtomSingle(Float, t)
	case token.STRING:
		node = ast.processKindAtomSingle(String, t)
	case token.KEYWORD:
		node = ast.processKindAtomSingle(Keyword, t)
	default:
		node, err = ast.processAtomSingle(t)
	}
	return node, err
}

func (ast *AST) processKindAtomSingle(kind AtomKind, t TokenWraper) Node {
	return &AtomSingle{
		pos:     t.Pos,
		Kind:    kind,
		Content: t.Content,
	}
}

func (ast *AST) processAtomSingle(t TokenWraper) (node Node, err error) {
	// FIXME(damnever)
	switch t.Content {
	default:
		node = &Symbol{
			pos:     t.Pos,
			Content: t.Content,
		}
		// err = fmt.Errorf("[%s] illegal syntax: %s", t.Pos, t.Content)
	}
	return
}

func (ast *AST) processAtomContainer(kind AtomKind, t TokenWraper) (Node, error) {
	node := &AtomContainer{
		pos:   t.Pos,
		Kind:  kind,
		Elems: []Node{},
	}
	var endToken token.Token
	if kind == Vector {
		endToken = token.RBRACK
	} else if kind == Map {
		endToken = token.RBRACE
	}
	t, err := ast.processContainer(&(node.Elems), endToken)
	if err != nil {
		return nil, err
	}
	node.end = t.End
	return node, err
}

func (ast *AST) processContainer(elems *[]Node, endToken token.Token) (t TokenWraper, err error) {
	var n Node
	for {
		t, err = ast.tr.Peek()
		if err != nil {
			return
		}
		if t.Token == token.EOF {
			err = fmt.Errorf("[%s] unexpected EOF", t.Pos)
			return
		}
		if t.Token == endToken {
			break
		}
		n, err = ast.processForm()
		if err != nil {
			return
		}
		*elems = append(*elems, n)
	}
	t, _ = ast.tr.Next()
	return
}
