package token

import "fmt"

type Token uint

const (
	ILLEGAL Token = iota
	EOF
	COMMENT

	// Basic type literals
	NIL     // nil
	BOOL    // true/false
	INT     // 12345
	FLOAT   // 123.45
	STRING  // "abc"
	KEYWORD // :abc

	// Delimiters
	LPAREN // (
	RPAREN // )
	LBRACK // [
	RBRACK // ]
	LBRACE // {
	RBRACE // }

	// Math functions
	ADD // +
	SUB // -
	MUL // *
	QUO // /

	TILDEAT     // ~@
	SINGLEQUOTE // '
	BACKQUOTE   // `
	TILDE       // ~
	CIRCUMFLEX  // ^
	ATSIGN      // @
	ASNSCS      // [^\s\[\]{}()'"`@,;]+ a sequence of zero or more non special characters
)

type Pos struct {
	Offset int
	Line   int
	Column int
}

func (pos Pos) String() string {
	return fmt.Sprintf("line:%d, column:%d", pos.Line, pos.Column)
}
