package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT // fields, table name

	// Misc characters
	ASTERISK // *
	COMMA    // ,

	// Keywords
	SELECT
	FROM
)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

var eof = rune(0)

// Scanner represents lexical scanner
type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// scan returns the next token and literal value
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace
	// if we see a letter then consume and ident our reserved word
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	}

	// Otherwise read the individual character
	switch ch {
	case eof:
		return EOF, ""
	case '*':
		return ASTERISK, string(ch)
	case ',':
		return COMMA, string(ch)
	}

	return ILLEGAL, string(ch)
}

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// create a buffer and read the current character in it
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// read every subsequent whitespace character into the buffer
	// Non-whitespace characters and EOF whill cause the loop to exit
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdent consumes current rune and all following ident runes
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// create a buffer and read the current character into it
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every following ident rune in the buffer
	// Non ident or eof will cause to exit
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	// if the string matches a keyword then return that keyword
	switch strings.ToUpper(buf.String()) {
	case "SELECT":
		return SELECT, buf.String()
	case "FROM":
		return FROM, buf.String()
	}

	// otherwise return as a regular identifier
	return IDENT, buf.String()
}

type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

func NewParse(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

func (p *Parser) Parse() (*SelectStatement, error) {
	stmt := &SelectStatement{}
	if tok, lit := p.scanIgnoreWhitespace(); tok != SELECT {
		return nil, fmt.Errorf("found %q, expected SELECT", lit)
	}

	for {
		// Read a field.
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmt.Fields = append(stmt.Fields, lit)

		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
	}

	// Next we should see the "FROM" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != FROM {
		return nil, fmt.Errorf("found %q, expected FROM", lit)
	}

	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmt.TableName = lit

	return stmt, nil
}

type SelectStatement struct {
	Fields    []string
	TableName string
}

func main() {
	p := NewParse(strings.NewReader("select name, age from things"))
	stmt, err := p.Parse()
	fmt.Printf("Got stmt %v with error: %v", stmt, err)
}
