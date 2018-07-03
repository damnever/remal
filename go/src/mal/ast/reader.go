package ast

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"mal/ast/token"
)

// XXX(damnever): using stack??????????

type TokenWraper struct {
	Token   token.Token
	Pos     token.Pos
	End     token.Pos
	Content string
}

type tokenReader struct {
	rd     *bufio.Reader
	pos    token.Pos
	end    token.Pos
	buf    *bytes.Buffer
	tokens []TokenWraper
}

func newTokenReader(code string) *tokenReader {
	return &tokenReader{
		rd:     bufio.NewReader(bytes.NewBufferString(code)),
		buf:    new(bytes.Buffer),
		tokens: []TokenWraper{},
		end:    token.Pos{Line: 1},
	}
}

func (r *tokenReader) Peek() (tw TokenWraper, err error) {
	if n := len(r.tokens); n > 0 {
		tw = r.tokens[n-1]
		return
	}
	tw, err = r.nexttoken()
	if err != nil {
		return
	}
	r.tokens = append(r.tokens, tw)
	return
}

func (r *tokenReader) Next() (tw TokenWraper, err error) {
	if n := len(r.tokens); n > 0 {
		tw = r.tokens[n-1]
		r.tokens = r.tokens[:n-1]
		return
	}

	return r.nexttoken()
}

func (r *tokenReader) nexttoken() (tw TokenWraper, err error) {
	r.buf.Reset()
	r.pos = r.end
	var t token.Token

	var b byte
	b, err = r.peekByte()
	if err != nil {
		if err == io.EOF {
			t, err = token.EOF, nil
			tw = r.makeTokenWrapper(t)
		}
		return
	}

	switch b {
	case '(':
		r.discardByte()
		t = token.LPAREN
	case ')':
		r.discardByte()
		t = token.RPAREN
	case '[':
		r.discardByte()
		t = token.LBRACK
	case ']':
		r.discardByte()
		t = token.RBRACK
	case '{':
		r.discardByte()
		t = token.LBRACE
	case '}':
		r.discardByte()
		t = token.RBRACE
	case '"':
		t, err = r.readString()
	case ';':
		t, err = r.readComment()
	case ':':
		t, err = r.readKeyword()
	case '\'', '`', '~', '^', '@':
		t, err = r.readSpecialSymbol()
	case '-': // '+' ????
		var bs []byte
		bs, err = r.peekBytes(2)
		if err == nil && (bs[1] >= '0' && bs[1] <= '9') {
			t, err = r.readNumber()
		} else {
			t, err = r.readSymbols()
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		t, err = r.readNumber()
	case ' ', '\n', '\r', '\t', '\f', '\b', ',':
		r.discardByte()
		return r.nexttoken()
	default:
		t, err = r.readSymbols()
	}

	if err != nil {
		return
	}
	tw = r.makeTokenWrapper(t)
	return
}

func (r *tokenReader) makeTokenWrapper(t token.Token) TokenWraper {
	return TokenWraper{
		Token:   t,
		Pos:     r.pos,
		End:     r.end,
		Content: r.buf.String(),
	}
}

func (r *tokenReader) readNumber() (t token.Token, err error) {
	isfloat := false
	b, _ := r.nextByte()
	r.buf.WriteByte(b)

	for {
		b, err = r.peekByte()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}
		if !(b >= '0' && b <= '9') && b != '.' {
			break
		}
		r.discardByte()
		if b == '.' {
			if isfloat {
				return
			}
			isfloat = true
		}
		r.buf.WriteByte(b)
	}

	if isfloat {
		t = token.FLOAT
	} else {
		t = token.INT
	}
	return
}

func (r *tokenReader) readString() (t token.Token, err error) {
	isfirst, slashs := true, 0
	before := func() {
		if slashs > 0 {
			r.buf.WriteString(strings.Repeat("\\", slashs))
		}
	}
	for {
		var b byte
		b, err = r.nextByte()
		if err != nil {
			if err == io.EOF {
				t, err = token.EOF, nil
			}
			return
		}
		switch b {
		case '\n':
			before()
			return
		case '\\':
			slashs++
		case '"':
			before()
			r.buf.WriteByte(b)
			if !isfirst && slashs%2 == 0 {
				t = token.STRING
				return
			}
			isfirst, slashs = false, 0
		default:
			before()
			slashs = 0
			r.buf.WriteByte(b)
		}
	}
}

func (r *tokenReader) readSpecialSymbol() (t token.Token, err error) {
	b, _ := r.nextByte()
	if b == '~' {
		r.buf.WriteByte(b)

		b1, err1 := r.peekByte()
		if err1 == nil && b1 == '@' {
			r.discardByte()
			r.buf.WriteByte(b1)
			t = token.TILDEAT
			return
		}

		if err != nil && err != io.EOF {
			return
		}
		t = token.TILDE
		return
	}

	switch b {
	case '\'':
		t = token.SINGLEQUOTE
	case '`':
		t = token.BACKQUOTE
	case '^':
		t = token.CIRCUMFLEX
	case '@':
		t = token.ATSIGN
	}
	r.buf.WriteByte(b)
	return
}

func (r *tokenReader) readSymbols() (t token.Token, err error) {
	if err = r.findSeqBytes(); err != nil {
		return
	}
	if r.buf.Len() == 0 {
		t = token.EOF
		return
	}
	switch r.buf.String() {
	case "nil":
		t = token.NIL
	case "true", "false":
		t = token.BOOL
	default:
		t = token.ASNSCS
	}
	return
}

func (r *tokenReader) readKeyword() (t token.Token, err error) {
	if err = r.findSeqBytes(); err == nil {
		t = token.KEYWORD
	}
	return
}

func (r *tokenReader) findSeqBytes() error {
	for {
		b, err := r.peekByte()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch b {
		case ' ', '\n', '\r', '\t', '\f', '\b', '(', ')', '[', ']', '{', '}':
			return nil
		default:
			r.buf.WriteByte(b)
			r.discardByte()
		}
	}
}

func (r *tokenReader) readComment() (t token.Token, err error) {
	for {
		var b byte
		b, err = r.nextByte()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}

		if b == '\n' {
			break
		}
		r.buf.WriteByte(b)
	}

	t = token.COMMENT
	return
}

func (r *tokenReader) peekByte() (b byte, err error) {
	var bs []byte
	bs, err = r.rd.Peek(1)
	if err == nil {
		b = bs[0]
	}
	return
}

func (r *tokenReader) peekBytes(n int) ([]byte, error) {
	return r.rd.Peek(n)
}

func (r *tokenReader) nextByte() (b byte, err error) {
	b, err = r.rd.ReadByte()
	if err != nil {
		return
	}
	r.advancePos([]byte{b})
	return
}

func (r *tokenReader) discardByte() {
	r.discardBytes(1)
}

func (r *tokenReader) discardBytes(n int) {
	p := make([]byte, n)
	r.rd.Read(p)
	r.advancePos(p)
}

func (r *tokenReader) advancePos(bs []byte) {
	for _, b := range bs {
		r.end.Offset++
		if b == '\n' {
			r.end.Line++
			r.end.Column = 0
		} else {
			r.end.Column++
		}
	}
}
