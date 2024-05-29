package parser

import (
	"fmt"
	"strings"

	"github.com/njorda.github.io/go/CloudSQL/lexer"
)

type KeyValue struct {
	Key   string
	Value string
}

type SQLQuery struct {
	Select []string
	From   string
	Where  KeyValue
	Equals KeyValue
}

// Lets do it like the parser instead!
type Parser struct {
	query lexer.Lexer
}

func NewParser(input string) *Parser {
	return &Parser{query: *lexer.NewLexer(input)}
}

func (p *Parser) nextIdentifier() string {
	for {
		tok := p.query.NextToken()
		if tok.Type == lexer.TOKEN_SYMBOL {
			continue
		}
		return tok.Literal
	}
}

// ParseSQLQuery parses a simple SQL query
func (p *Parser) ParseQuery() (*SQLQuery, error) {
	query := &SQLQuery{}
	for tok := p.query.NextToken(); tok.Type != lexer.TOKEN_EOF; tok = p.query.NextToken() {
		switch tok.Type {
		case lexer.TOKEN_KEYWORD:
			switch strings.ToUpper(tok.Literal) {
			case "SELECT":
				for tok = p.query.NextToken(); tok.Type == lexer.TOKEN_IDENTIFIER || tok.Type == lexer.TOKEN_SYMBOL; tok = p.query.NextToken() {
					switch tok.Type {
					case lexer.TOKEN_SYMBOL:
						continue
					case lexer.TOKEN_IDENTIFIER:
						query.Select = append(query.Select, tok.Literal)
					}
				}
				fallthrough
			case "FROM":
				// no inner query support
				query.From = p.query.NextToken().Literal
				fmt.Printf("FROM: %v\n", query.From)
			// Currently only supports one where clause, either with = or %.
			case "WHERE":

				kV := KeyValue{}
				kV.Key = p.nextIdentifier()
				kV.Value = p.nextIdentifier()
				fmt.Printf("WHERE: %v, %v\n", kV.Key, kV.Value)

				if strings.HasSuffix(kV.Value, "%") {
					kV.Value = strings.TrimSuffix(kV.Value, "%")
					query.Where = kV
					continue
				}
				query.Equals = kV

			}
		}
	}
	return query, nil
}
