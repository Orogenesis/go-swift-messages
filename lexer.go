package swiftmessages

import (
	"bufio"
	"bytes"
	"io"
)

// Lexer represents a lexical scanner.
type Lexer struct {
	*bufio.Reader
	prev Token
}

// NewLexer returns a new instance of Lexer.
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{Reader: bufio.NewReader(r)}
}

// Scan returns the next token and literal value.
func (l *Lexer) Scan() (tok Token, lit string) {
	ch := l.read()
	switch ch {
	case rune(0):
		tok, lit = TokenEOF, ""
	case CharacterLBrace:
		tok, lit = TokenLBrace, string(ch)
	case CharacterRBrace:
		tok, lit = TokenRBrace, string(ch)
	case CharacterColon:
		if l.prev == TokenLinebreak {
			tok, lit = TokenLinebreakColon, string(ch)
		} else {
			tok, lit = TokenColon, string(ch)
		}
	case CharacterLinebreak:
		tok, lit = TokenLinebreak, string(ch)
	default:
		l.unread()
		switch l.prev {
		case TokenColon, TokenLinebreak:
			tok, lit = TokenString, l.scanIdent(CharacterRBrace)
		default:
			tok, lit = TokenID, l.scanIdent(CharacterColon)
		}
	}

	l.prev = tok
	return tok, lit
}

// Consumes the current rune and all contiguous ident runes.
func (l *Lexer) scanIdent(stopCharacter rune) string {
	buf := &bytes.Buffer{}
	buf.WriteRune(l.read())
	for {
		if ch := l.read(); ch == rune(0) {
			break
		} else if ch == stopCharacter || ch == CharacterLinebreak {
			l.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return buf.String()
}

// Reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (l *Lexer) read() rune {
	ch, _, err := l.ReadRune()
	if err != nil {
		return rune(0)
	}

	return ch
}

// Places the previously read rune back on the reader.
func (l *Lexer) unread() { _ = l.UnreadRune() }
