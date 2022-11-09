// Code generated by github.com/gnolang/gno. DO NOT EDIT.

//go:build gno
// +build gno

package exts // import "gno.land/p/demo/grc/exts"

type TokenMetadata interface {
	// Returns the name of the token.
	GetName() string

	// Returns the symbol of the token, usually a shorter version of the
	// name.
	GetSymbol() string

	// Returns the decimals places of the token.
	GetDecimals() uint
}
