// Package swiftmessages implements a parser for SWIFT financial messages
package swiftmessages

import (
	"errors"
)

// ErrSwiftBlockInvalid is raised when an invalid SWIFT financial message received.
var ErrSwiftBlockInvalid = errors.New("swift block is invalid")

// SwiftMessage represents a collection of SWIFT financial blocks.
type SwiftMessage struct {
	Blocks []SwiftBlock
}

// SwiftBlock represents a SWIFT block.
type SwiftBlock struct {
	ID    string
	Value interface{}
}

// Parser represents a parser.
type Parser struct {
	*Lexer
}

// NewParser returns a new instance of Parser.
func NewParser(l *Lexer) *Parser {
	return &Parser{Lexer: l}
}

// Parse parses SWIFT message.
func (p *Parser) Parse() (message SwiftMessage, err error) {
	for {
		tok, _ := p.Lexer.Scan()
		switch tok {
		case TokenLBrace:
			swiftBlock, err := p.parseBlock()
			if err != nil {
				return message, err
			}

			message.Blocks = append(message.Blocks, swiftBlock)
		default:
			return message, nil
		}
	}
}

// Parses SWIFT block.
func (p *Parser) parseBlock() (SwiftBlock, error) {
	var (
		tok        Token
		lit        string
		swiftBlock SwiftBlock
	)

	if tok, lit = p.Lexer.Scan(); tok != TokenID {
		return SwiftBlock{}, ErrSwiftBlockInvalid
	}

	// Block's identifier
	swiftBlock.ID = lit

	if tok, _ = p.Lexer.Scan(); tok != TokenColon {
		return SwiftBlock{}, ErrSwiftBlockInvalid
	}

	tok, lit = p.Lexer.Scan()
	switch tok {
	case TokenLBrace:
		_ = p.Lexer.UnreadRune()

		for {
			if tok, _ := p.Lexer.Scan(); tok != TokenLBrace {
				_ = p.Lexer.UnreadRune()
				break
			}

			newBlock, err := p.parseBlock()
			if err != nil {
				return SwiftBlock{}, err
			}

			if _, ok := swiftBlock.Value.([]SwiftBlock); !ok {
				swiftBlock.Value = make([]SwiftBlock, 0)
			}

			swiftBlock.Value = append(swiftBlock.Value.([]SwiftBlock), newBlock)
		}
	case TokenLinebreak:
		for {
			tok, lit := p.Lexer.Scan()
			if tok == TokenRBrace {
				_ = p.Lexer.UnreadRune()
				break
			} else if tok == TokenColon || tok == TokenLinebreakColon || tok == TokenLinebreak {
				continue
			} else if len(lit) == 0 {
				continue
			}

			if _, ok := swiftBlock.Value.([]SwiftBlock); !ok {
				swiftBlock.Value = make([]SwiftBlock, 0)
			}

			if tok == TokenID {
				swiftBlock.Value = append(swiftBlock.Value.([]SwiftBlock), SwiftBlock{
					ID:    lit,
					Value: "",
				})
			} else if tok == TokenString {
				if swiftBlockValue, ok := swiftBlock.Value.([]SwiftBlock); ok {
					if len(swiftBlockValue) == 0 {
						return SwiftBlock{}, ErrSwiftBlockInvalid
					}

					idx := len(swiftBlockValue) - 1

					if len(swiftBlockValue[idx].Value.(string)) != 0 {
						swiftBlockValue[idx].Value = swiftBlockValue[idx].Value.(string) + "\n"
					}

					swiftBlockValue[idx].Value = swiftBlockValue[idx].Value.(string) + lit
				}
			}
		}
	case TokenRBrace:
		swiftBlock.Value = ""
	case TokenString:
		swiftBlock.Value = lit
	default:
		return SwiftBlock{}, ErrSwiftBlockInvalid
	}

	if tok != TokenRBrace {
		if tok, lit = p.Lexer.Scan(); tok != TokenRBrace {
			return SwiftBlock{}, ErrSwiftBlockInvalid
		}
	}

	return swiftBlock, nil
}
