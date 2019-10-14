package swiftmessages

// Token represents a lexical token.
type Token int

const (
	// Tokens, literals and keywords
	TokenIllegal Token = iota
	TokenEOF
	TokenLBrace
	TokenRBrace
	TokenString
	TokenID
	TokenColon
	TokenLinebreak
	TokenLinebreakColon

	// Characters
	CharacterColon     = ':'
	CharacterLBrace    = '{'
	CharacterRBrace    = '}'
	CharacterLinebreak = '\n'
)
