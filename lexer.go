package go_swift_messages

import (
	"bufio"
	"bytes"
	"io"
)

type Lexer struct {
	*bufio.Reader
	prev Token
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{Reader: bufio.NewReader(r)}
}

func (l *Lexer) Scan() (tok Token, lit string) {
	ch := l.read()
	switch ch {
	case rune(0):
		tok, lit = EOF, ""
	case CHARACTER_LBRACKET:
		tok, lit = LBRACKET, string(ch)
	case CHARACTER_RBRACKET:
		tok, lit = RBRACKET, string(ch)
	case CHARACTER_COLON:
		if l.prev == LINEBREAK {
			tok, lit = LINEBREAK_COLON, string(ch)
		} else {
			tok, lit = COLON, string(ch)
		}
	case CHARACTER_LINEBREAK:
		tok, lit = LINEBREAK, string(ch)
	default:
		l.unread()
		switch l.prev {
		case COLON, LINEBREAK:
			tok, lit = STRING, l.scanIdent(CHARACTER_RBRACKET)
		default:
			tok, lit = ID, l.scanIdent(CHARACTER_COLON)
		}
	}

	l.prev = tok
	return tok, lit
}

func (l *Lexer) scanIdent(stopCharacter rune) string {
	buf := &bytes.Buffer{}
	buf.WriteRune(l.read())

	for {
		if ch := l.read(); ch == rune(0) {
			break
		} else if ch == stopCharacter || ch == CHARACTER_LINEBREAK {
			l.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return buf.String()
}

func (l *Lexer) read() rune {
	ch, _, err := l.ReadRune()
	if err != nil {
		return rune(0)
	}

	return ch
}

func (l *Lexer) unread() { _ = l.UnreadRune() }
