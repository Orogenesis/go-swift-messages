package go_swift_messages

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	t.Run("when parse header block", func(t *testing.T) {
		r := strings.NewReader(`{1:F01AAAAGRA0AXXX0057000289}`)
		p := NewParser(NewLexer(r))
		message, err := p.Parse()
		assert.NoError(t, err)
		assert.Equal(t, "1", message.Blocks[0].ID)
		assert.Equal(t, "F01AAAAGRA0AXXX0057000289", message.Blocks[0].Value)
	})

	t.Run("when parse empty block", func(t *testing.T) {
		r := strings.NewReader(`{10:}`)
		p := NewParser(NewLexer(r))
		message, err := p.Parse()
		assert.NoError(t, err)
		assert.Equal(t, "10", message.Blocks[0].ID)
		assert.Empty(t, message.Blocks[0].Value)
	})

	t.Run("when parse multiple blocks", func(t *testing.T) {
		r := strings.NewReader(`{1:F01AAAAGRA0AXXX0057000289}{2:O1030919010321BBBBGRA0AXXX00570001710103210920N}`)
		p := NewParser(NewLexer(r))
		message, err := p.Parse()
		assert.NoError(t, err)
		assert.Equal(t, "F01AAAAGRA0AXXX0057000289", message.Blocks[0].Value)
		assert.Equal(t, "O1030919010321BBBBGRA0AXXX00570001710103210920N", message.Blocks[1].Value)
	})

	t.Run("when parse block with nested blocks", func(t *testing.T) {
		r := strings.NewReader(`{3:{108:MT103 003 OF 045}{121:c8b66b47-2bd9-48fe-be90-93c2096f27d2}}`)
		p := NewParser(NewLexer(r))
		message, err := p.Parse()
		assert.NoError(t, err)
		assert.Len(t, message.Blocks[0].Value, 2)
	})

	t.Run("when nested blocks are named with letters", func(t *testing.T) {
		r := strings.NewReader(`{5:{MAC:00000000}{CHK:24857F4599E7}{TNG:}}`)
		p := NewParser(NewLexer(r))
		message, err := p.Parse()
		assert.NoError(t, err)
		assert.Len(t, message.Blocks[0].Value, 3)
		assert.Equal(t, "MAC", message.Blocks[0].Value.([]SwiftBlock)[0].ID)
		assert.Equal(t, "CHK", message.Blocks[0].Value.([]SwiftBlock)[1].ID)
		assert.Equal(t, "TNG", message.Blocks[0].Value.([]SwiftBlock)[2].ID)
	})

	t.Run("when parse multiline block", func(t *testing.T) {
		data := []string{
			"{4:",
			":20:5387354",
			":23B:CRED",
			":23E:PHOB/20.527.19.60",
			":32A:000526USD1101,50",
			":33B:USD1121,50",
			":50K:FRANZ HOLZAPFEL GMBH",
			"VIENNA",
			":52A:BKAUATWW",
			":59:723491524",
			"C. KLEIN",
			"BLOEMENGRACHT 15",
			"AMSTERDAM",
			":71A:SHA",
			":71F:USD10,",
			"-}",
		}

		r := strings.NewReader(strings.Join(data, "\n"))
		p := NewParser(NewLexer(r))
		message, err := p.Parse()
		assert.NoError(t, err)
		assert.Equal(t, "20", message.Blocks[0].Value.([]SwiftBlock)[0].ID)
		assert.Equal(t, "5387354", message.Blocks[0].Value.([]SwiftBlock)[0].Value)
	})
}
