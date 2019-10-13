# go-swift-messages

go-swift-messages parse SWIFT messages with different message types (MT101, MT103, MT104, MT202, MT509, MT900, MT910, MT940, MT942, MT950) into abstract syntax tree (AST) and convert ASTs back to SWIFT message.  

### What's an MT103?

MT103 is a SWIFT payment message type/format used for cash transfer specifically for cross border/international wire transfer.

### Quick Usage

```go
package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/orogenesis/go-swift-messages"
)

func main()  {
	p := go_swift_messages.NewParser(go_swift_messages.NewLexer(&bytes.Buffer{}))
	message, err := p.Parse()
	if err != nil {
		log.Fatal(err)
	}

	for _, swiftBlock := range message.Blocks {
		fmt.Println("ID:", swiftBlock.ID)
		
		switch v := swiftBlock.Value.(type) {
		case []go_swift_messages.SwiftBlock:
			for _, newBlock := range v {
				fmt.Printf("ID: %v, value: %v\n", newBlock.ID, newBlock.Value)
			}
		case string:
			fmt.Println(v)
		}
	}
}
```