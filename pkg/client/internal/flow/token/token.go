package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// identifiers, literals
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	DOT = "."

	// delimiters
	LBRACKET = "["
	RBRACKET = "]"
)

type Token struct {
	Type    TokenType
	Literal string
}

func LookupIdent(ident string) TokenType {
	return IDENT
}
