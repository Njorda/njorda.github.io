---
layout: post
title: "SQL for gcloud bucket"
subtitle: "Building a SQL tool to query Google Cloud Bucket"
date: 2024-01-16
author: "Niklas Hansson"
URL: "/2024/01/16/"
---

TLDR: https://github.com/Njorda/cloudsql/tree/main

The goal with this blog post is to build a small tool to query Google Cloud buckets. We will do this using ChatGPT, I feel like I need to be even better at prompting for coding help. So to start by setting some constraints to limit the scope and make it a reasonable task to finish in a hour we will use `go` and we will limit our self to `SELECT` with `WHERE` we will support predicate push downs but thats it. We also aim for creating an experiance that is similar to [psql](https://www.postgresql.org/docs/current/app-psql.html#:~:text=Description,or%20from%20command%20line%20arguments.)

The first step is to be able to parse queries, we will not support [Common Table Expressions](https://www.postgresql.org/docs/current/queries-with.html) CTE or any other fancy feature for that part to make life easier. Initially I though to not include the exact prompts but after realising it might be tricky for people to play around with and reproduce I decided to keep them in an appendix. However in the end I changed my mind again to not do this since it was so much back and fourth with ChatGPT4. The first step in order to build my small CLI is to get the overview tasks, the steps I decided to follow for the project is the following: 

1) Lexer or Tokenizer: The lexer takes the SQL query as input and breaks it down into tokens. Tokens are the smallest units that have meaning in SQL, like keywords (SELECT, FROM, WHERE), identifiers (table names, column names), operators (=, >, <), and literals (numeric values, strings).
2) Intermediate Representation: The parser typically converts the query into an intermediate representation (IR). This IR is a data structure (like a parse tree or an abstract syntax tree) that represents the parsed query. The IR is used by other components of the database system to execute the query.
3) Generate API call: Extract the paramters for the API call to google cloud storage.


The code will be broken down to the following structure:


├── lexer
│   ├── lexer.go
│   ├── lexer_test.go
├── parser
│   ├── sql_parser.go
│   ├── sql_parser_test.go
├── main.go
├── go.mod
├── go.sum

The code for making the API calls will live in main. 

The first step is to build the parser, ChatGPT4 did most of the heavy lifting and the code only required some minor changes: 


```go
package lexer

import (
	"unicode"
)

type TokenType int

var keywords = []string{"SELECT", "FROM", "WHERE"}

const (
	TOKEN_EOF TokenType = iota
	TOKEN_ERROR
	TOKEN_KEYWORD
	TOKEN_IDENTIFIER
	TOKEN_SYMBOL
)

type Token struct {
	Type    TokenType
	Literal string
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case 0:
		tok = Token{Type: TOKEN_EOF, Literal: string("")}
	case '=', ';', '(', ')', ',', '\'', '/':
		tok = Token{Type: TOKEN_SYMBOL, Literal: string(l.ch)}
		l.readChar()
	default:
		if isLetter(l.ch) {

			literal := l.readIdentifier()
			if keyword(literal) {
				tok = Token{Type: TOKEN_KEYWORD, Literal: literal}
				return tok
			}
			tok = Token{Type: TOKEN_IDENTIFIER, Literal: literal}
			return tok
		}
		tok = Token{Type: TOKEN_ERROR, Literal: string(l.ch)}
		l.readChar()
	}
	return tok
}

func keyword(literal string) bool {
	for _, keyword := range keywords {
		if keyword == literal {
			return true
		}
	}
	return false
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_' || ch == '%'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
```

The code also contains a test which shows the output, check [here](https://github.com/Njorda/cloudsql/blob/main/lexer/lexer_test.go). Next step is to parse the tokenized the input: 


```go
package parser

import (
	"fmt"
	"strings"

	"github.com/Njorda.cloudsql/lexer"
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
			// Currently only supports one where clause, either with = or %.
			case "WHERE":
				kV := KeyValue{}
				kV.Key = p.nextIdentifier()
				// Need to get all the stuff until we get the end of it.
				kV.Value = p.nextIdentifier()
			Exit:
				for tok = p.query.NextToken(); tok.Type == lexer.TOKEN_IDENTIFIER || tok.Type == lexer.TOKEN_SYMBOL; tok = p.query.NextToken() {
					switch {
					case tok.Type == lexer.TOKEN_SYMBOL && tok.Literal == `/`:
						kV.Value = fmt.Sprintf("%v%v", kV.Value, tok.Literal)
					case tok.Type == lexer.TOKEN_SYMBOL && tok.Literal == `=`:
						continue
					case tok.Type == lexer.TOKEN_IDENTIFIER:
						kV.Value = fmt.Sprintf("%v%v", kV.Value, tok.Literal)
					default:
						continue Exit
					}
				}
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
```


Most of the code so far as been written baser upon the go standard libs. However for the CLI we will use a [library](github.com/chzyer/readline) in order to give a lot of the basic functionality such as searching previous commands and history. We will also use a [library](github.com/jedib0t/go-pretty/v6/table) for pretty printing the results.


```go
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/Njorda.cloudsql/parser"
	"google.golang.org/api/iterator"

	"github.com/chzyer/readline"

	"github.com/jedib0t/go-pretty/v6/table"
)

//var columns = []string{"name", "size", "timeCreated", "timeUpdated", "storageClass", "owner", "contentType", "contentEncoding", "contentDisposition", "retentionTime", "updated"}

func handleInput(ctx context.Context, client *storage.Client, input string) error {
	query, err := parser.NewParser(input).ParseQuery()
	if err != nil {
		return err
	}

	// Name is the only value we have ...
	rows, err := ListObjects(ctx, client, query.From, query.Where.Value, query.Select)
	if err != nil {
		return err
	}
	format(query.Select, rows)
	return nil
}

// CreateClient initializes a new Google Cloud Storage client
func CreateClient(ctx context.Context) (*storage.Client, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// GetObjects lists objects in a given bucket, optionally filtered by a prefix
func GetObjects(ctx context.Context, client *storage.Client, bucketName, objet string, selected []string) ([]string, error) {
	attrs, err := client.Bucket(bucketName).Object(objet).Attrs(ctx)
	if err != nil {
		return nil, err
	}
	var objects []string
	out := map[string]string{}
	out["name"] = attrs.Name
	out["size"] = fmt.Sprintf("%d", attrs.Size)
	out["timeCreated"] = attrs.Created.String()
	out["timeUpdated"] = attrs.Updated.String()
	out["storageClass"] = attrs.StorageClass
	out["owner"] = attrs.Owner
	out["contentType"] = attrs.ContentType
	out["contentEncoding"] = attrs.ContentEncoding
	out["contentDisposition"] = attrs.ContentDisposition
	out["retentionTime"] = attrs.RetentionExpirationTime.GoString()
	out["updated"] = attrs.Updated.String()
	for _, column := range selected {
		objects = append(objects, out[column])
	}
	return objects, nil
}

// ListObjects lists objects in a given bucket, optionally filtered by a prefix
func ListObjects(ctx context.Context, client *storage.Client, bucketName, prefix string, selected []string) ([][]string, error) {
	fmt.Println("The prefix is: ", prefix)
	it := client.Bucket(bucketName).Objects(ctx, &storage.Query{Prefix: prefix})
	var rows [][]string
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		// Here we could do something with reflect to get all the stuff out!
		// This would be the way.
		out := map[string]string{}
		out["name"] = attrs.Name
		out["size"] = fmt.Sprintf("%d", attrs.Size)
		out["timeCreated"] = attrs.Created.String()
		out["timeUpdated"] = attrs.Updated.String()
		out["storageClass"] = attrs.StorageClass
		out["owner"] = attrs.Owner
		out["contentType"] = attrs.ContentType
		out["contentEncoding"] = attrs.ContentEncoding
		out["contentDisposition"] = attrs.ContentDisposition
		out["retentionTime"] = attrs.RetentionExpirationTime.GoString()
		out["updated"] = attrs.Updated.String()
		row := []string{}
		for _, column := range selected {
			row = append(row, out[column])
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func format(columns []string, tuples [][]string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	row := table.Row{}
	for _, col := range columns {
		row = append(row, col)
	}
	t.AppendHeader(row)
	rows := []table.Row{}
	for _, tuple := range tuples {
		row := table.Row{}
		for _, val := range tuple {
			row = append(row, val)
		}
		rows = append(rows, row)
	}
	t.AppendRows(rows)
	t.AppendSeparator()
	t.Render()
}

func main() {
	ctx := context.Background()
	client, err := CreateClient(ctx) // Assuming CreateClient is a function you've defined
	if err != nil {
		panic(err)
	}

	rl, err := readline.New("GCSQL> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	fmt.Println("Welcome to GCSQL, the Google Cloud Storage SQL interface.")
	fmt.Println("Type 'EXIT' to quit.")

	for {
		input, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}

		if strings.ToUpper(input) == "EXIT" {
			fmt.Println("Goodbye!")
			break
		}

		// Add the input to history
		rl.SaveHistory(input)

		// Handle the input
		if err := handleInput(ctx, client, input); err != nil {
			fmt.Printf("Error: %v", err)
		}
	}
}

```

In order to build the binary and try the CLI out `go build -o cloudsql` and then `./cloudsql`:

```bash
GCSQL> SELECT name, size FROM ceedai WHERE path = 'ceedai/files/%';
+---------------------------------------------------------------------------------+----------+
| NAME                                                                            | SIZE     |
+---------------------------------------------------------------------------------+----------+
| ceedai/files/tmp1.pdf                                                           | 3372362  |
| ceedai/files/smp.pdf                                                            | 989059   |
| ceedai/files/12.pdf                                                             | 72202    |
| ceedai/files/p29-neumann-cidr20.pdf                                             | 335856   |
| ceedai/files/simon.json                                                         | 30       |
+---------------------------------------------------------------------------------+----------+
GCSQL>
```

The current implementation is very limited but demonstrates how it can be implemented. Feel free to fork the code and add features. I might come back in later posts and add more features.
