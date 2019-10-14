[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b023b6b1b3ac465591c0b519eac15d5a)](https://www.codacy.com/manual/Orogenesis/go-swift-messages)
[![Build Status](https://travis-ci.org/Orogenesis/go-swift-messages.svg?branch=master)](https://travis-ci.org/Orogenesis/go-swift-messages)
[![Coverage Status](https://coveralls.io/repos/github/Orogenesis/go-swift-messages/badge.svg?branch=master)](https://coveralls.io/github/Orogenesis/go-swift-messages?branch=master)
[![GoDoc](http://godoc.org/github.com/orogenesis/go-swift-messages?status.svg)](http://godoc.org/github.com/orogenesis/go-swift-messages)

### go-swift-messages

Parses SWIFT financial messages with different message types (MT101, MT103, MT104, MT202, MT509, MT900, MT910, MT940, MT942, MT950) into abstract syntax tree (AST) and convert ASTs back to SWIFT financial message.  

### About MT103

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
	p := swiftmessages.NewParser(swiftmessages.NewLexer(&bytes.Buffer{}))
	message, err := p.Parse()
	if err != nil {
		log.Fatal(err)
	}

	for _, swiftBlock := range message.Blocks {
		fmt.Println("ID:", swiftBlock.ID)
		
		switch v := swiftBlock.Value.(type) {
		case []swiftmessages.SwiftBlock:
			for _, newBlock := range v {
				fmt.Printf("ID: %v, value: %v\n", newBlock.ID, newBlock.Value)
			}
		case string:
			fmt.Println(v)
		}
	}
}
```