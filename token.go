package swiftmessages

// Token represents a lexical token.
type Token int

const (
	ILLEGAL Token = iota
	EOF
	LBRACKET
	RBRACKET
	STRING
	ID
	COLON
	LINEBREAK
	LINEBREAK_COLON

	CHARACTER_COLON     = ':'
	CHARACTER_LBRACKET  = '{'
	CHARACTER_RBRACKET  = '}'
	CHARACTER_LINEBREAK = '\n'
)
