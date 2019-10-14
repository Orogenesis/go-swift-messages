package swiftmessages

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	t.Run("when nested blocks", func(t *testing.T) {
		message := MustEncode(
			SwiftBlockRule{
				SwiftBlock: SwiftBlock{
					ID:    "1",
					Value: "F01COPZBEB0AXXX0377002089",
				},
			},
			SwiftBlockRule{
				SwiftBlock: SwiftBlock{
					ID: "3",
					Value: []SwiftBlock{
						{
							ID:    "108",
							Value: "MT103 003 OF 045",
						},
						{
							ID:    "121",
							Value: "c8b66b47-2bd9-48fe-be90-93c2096f27d2",
						},
					},
				},
			},
			SwiftBlockRule{
				SwiftBlock: SwiftBlock{
					ID: "5",
					Value: []SwiftBlock{
						{
							ID:    "MAC",
							Value: "00000000",
						},
						{
							ID: "CHK",
							Value: []SwiftBlock{
								{
									ID:    "TNG",
									Value: "7E0FAA1CFBE1",
								},
							},
						},
					},
				},
			},
			SwiftBlockRule{
				SwiftBlock: SwiftBlock{
					ID: "4",
					Value: []SwiftBlock{
						{
							ID:    "20",
							Value: "00345",
						},
						{
							ID:    "13C",
							Value: "/A234567Z/1359+0100",
						},
					},
				},
				ShortMode: true,
			},
		)

		expected := []string{
			"{1:F01COPZBEB0AXXX0377002089}{3:{108:MT103 003 OF 045}{121:c8b66b47-2bd9-48fe-be90-93c2096f27d2}}{5:{MAC:00000000}{CHK:{TNG:7E0FAA1CFBE1}}}{4:",
			":20:00345",
			":13C:/A234567Z/1359+0100}",
		}

		assert.Equal(t, strings.Join(expected, "\n"), message)
	})
}
