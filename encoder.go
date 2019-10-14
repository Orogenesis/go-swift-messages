package swiftmessages

import (
	"container/list"
	"errors"
)

// ErrUnexpectedType is raised when an invalid value in a SWIFT block received.
var ErrUnexpectedType = errors.New("unexpected type")

// SwiftBlockRule represents a SWIFT block settings.
type SwiftBlockRule struct {
	SwiftBlock
	ShortMode bool
}

type swiftBlockElement struct {
	SwiftBlockRule
	closure   bool
	shortMode bool
}

// MustEncode returns message if err is nil and panics otherwise.
func MustEncode(swiftBlock ...SwiftBlockRule) string {
	m, err := Encode(swiftBlock...)
	if err != nil {
		panic(err)
	}

	return m
}

// Encode encodes swiftBlock into a SWIFT message or returns an error.
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
			message += string(CharacterLinebreak) + string(CharacterColon) + current.ID + string(CharacterColon)
		} else {
			// Add open brace character, block identifier and colon character
			message += string(CharacterLBrace) + current.ID + string(CharacterColon)
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
				message += string(CharacterRBrace)
			}

			if current.closure {
				for braceCounter > 0 {
					// Add close brace character
					message += string(CharacterRBrace)
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
