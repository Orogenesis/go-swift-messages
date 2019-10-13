// Package go_swift_messages implements a parser for SWIFT financial messages
package go_swift_messages

import (
	"errors"
)

var ErrSwiftBlockInvalid = errors.New("swift block is invalid")

type SwiftMessage struct {
	Blocks []SwiftBlock
}

type SwiftBlock struct {
	ID    string
	Value interface{}
}

// Represents a parser.
type Parser struct {
	*Lexer
}

// Returns a new instance of Parser.
func NewParser(l *Lexer) *Parser {
	return &Parser{Lexer: l}
}

// Parses SWIFT message.
func (p *Parser) Parse() (message SwiftMessage, err error) {
	for {
		tok, _ := p.Lexer.Scan()
		switch tok {
		case LBRACKET:
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

	if tok, lit = p.Lexer.Scan(); tok != ID {
		return SwiftBlock{}, ErrSwiftBlockInvalid
	}

	// Block's identifier
	swiftBlock.ID = lit

	if tok, _ = p.Lexer.Scan(); tok != COLON {
		return SwiftBlock{}, ErrSwiftBlockInvalid
	}

	tok, lit = p.Lexer.Scan()
	switch tok {
	case LBRACKET:
		_ = p.Lexer.UnreadRune()

		for {
			if tok, _ := p.Lexer.Scan(); tok != LBRACKET {
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
	case LINEBREAK:
		for {
			tok, lit := p.Lexer.Scan()
			if tok == RBRACKET {
				_ = p.Lexer.UnreadRune()
				break
			} else if tok == COLON || tok == LINEBREAK_COLON || tok == LINEBREAK {
				continue
			} else if len(lit) == 0 {
				continue
			}

			if _, ok := swiftBlock.Value.([]SwiftBlock); !ok {
				swiftBlock.Value = make([]SwiftBlock, 0)
			}

			if tok == ID {
				swiftBlock.Value = append(swiftBlock.Value.([]SwiftBlock), SwiftBlock{
					ID:    lit,
					Value: "",
				})
			} else if tok == STRING {
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
	case RBRACKET:
		swiftBlock.Value = ""
	case STRING:
		swiftBlock.Value = lit
	default:
		return SwiftBlock{}, ErrSwiftBlockInvalid
	}

	if tok != RBRACKET {
		if tok, lit = p.Lexer.Scan(); tok != RBRACKET {
			return SwiftBlock{}, ErrSwiftBlockInvalid
		}
	}

	return swiftBlock, nil
}
