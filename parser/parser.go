package parser

import (
	"fmt"

	"github.com/ulricksennick/monkey/ast"
	"github.com/ulricksennick/monkey/lexer"
	"github.com/ulricksennick/monkey/token"
)

// Parser implementing recursive-descent parsing
type Parser struct {
	l         *lexer.Lexer // lexer containing tokenized source code
	curToken  token.Token  // current token under examination
	peekToken token.Token  // next token; checked when forming program statements
	errors    []string     // errors due to incorrect token types (syntax errors)
}

// Create a new parser which will use the given lexer
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// Return a program node which represents the top node of abstract syntax tree
// generated by the parser. The AST will contain nodes representing the source
// code provided to the parser's lexer.
func (p *Parser) ParseProgram() *ast.Program {
	// Create a new program node
	program := &ast.Program{}
	// Program statements (children nodes of <program> in the AST)
	program.Statements = []ast.Statement{}

	// Iterate over tokens until end of file, parsing and appending statements
	// to the program node's statement list
	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	// Return the program
	return program
}

func (p *Parser) Errors() []string {
	return p.errors
}

// Add an error to the parser error list due to incorrect token type
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// Advance the parser's current and next tokens
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Parse a program statement depending on its token type starting with the
// parser's curToken
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

// Construct a "let statement" node with the parser's current token. Advance the
// parser tokens while checking/asserting the next token's type for the next
// expected token type. (let <IDENT> = <expression>)
func (p *Parser) parseLetStatement() *ast.LetStatement {
	// Create "let" statement with current token (LET token)
	stmt := &ast.LetStatement{Token: p.curToken}

	// Check for identifier
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Check for ASSIGN node, "="
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: For now, we skip over expressions until a semicolon is encountered
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	// Create "return" statement with current token (RETURN token)
	stmt := &ast.ReturnStatement{Token: p.curToken}

	// Advance the parser to beginning of expression to be parsed
	p.nextToken()

	// TODO: For now, we skip over expressions until a semicolon is encountered
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Check whether the next token is of the expected token type, advance the parser
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	// Token type mismatch, add an error to the parser
	p.peekError(t)
	return false
}

// Check if current token is of the given type
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// Check if next token is of the given type
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}
