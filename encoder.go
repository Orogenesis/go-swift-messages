package go_swift_messages

import (
	"container/list"
	"errors"
)

var ErrUnexpectedType = errors.New("unexpected type")

type SwiftBlockRule struct {
	SwiftBlock
	ShortMode bool
}

type swiftBlockElement struct {
	SwiftBlockRule
	closure   bool
	shortMode bool
}

func MustEncode(swiftBlock ...SwiftBlockRule) string {
	m, err := Encode(swiftBlock...)
	if err != nil {
		panic(err)
	}

	return m
}

func Encode(swiftBlock ...SwiftBlockRule) (message string, err error) {
	queue := list.New()
	braceCounter := 0

	for _, v := range swiftBlock {
		queue.PushBack(swiftBlockElement{
			SwiftBlockRule: v,
		})
	}

	for queue.Len() > 0 {
		element := queue.Front()
		current := element.Value.(swiftBlockElement)

		if current.shortMode {
			// Add colon character, block identifier, colon character and new line
			message += string(CHARACTER_LINEBREAK) + string(CHARACTER_COLON) + current.ID + string(CHARACTER_COLON)
		} else {
			// Add open brace character, block identifier and colon character
			message += string(CHARACTER_LBRACKET) + current.ID + string(CHARACTER_COLON)
		}

		switch v := current.Value.(type) {
		case []SwiftBlock:
			for i := len(v) - 1; i >= 0; i-- {
				queue.PushFront(swiftBlockElement{
					SwiftBlockRule: SwiftBlockRule{SwiftBlock: v[i]},
					closure:        i == (len(v) - 1),
					shortMode:      current.ShortMode,
				})
			}

			braceCounter++
		case string:
			message += v

			if !current.shortMode {
				// Add close brace character
				message += string(CHARACTER_RBRACKET)
			}

			if current.closure {
				for braceCounter > 0 {
					// Add close brace character
					message += string(CHARACTER_RBRACKET)
					braceCounter--
				}
			}
		default:
			return "", ErrUnexpectedType
		}

		// Dequeue element
		queue.Remove(element)
	}

	return message, nil
}
