package cli

import (
	"bytes"
	"strconv"
	"unicode"
)

const (
	defaultArgName    = "ARG"
	minLongFlagLength = 2
)

// Flag represents a command line flag, with a short and/or long variation.
type Flag struct {
	Short       rune   // the short flag (0 for no short flag)
	Long        string // the long flag (empty string for no long flag)
	Description string // the flag description

	Required bool // true if flag is required
	HasArg   bool // true if the flag has an argument

	ArgName string // the argument name for the help formatter
}

// NewFlag constructs a new flag.
func NewFlag(short rune, long string, desc string, hasArg bool) *Flag {
	return &Flag{
		Short:       short,
		Long:        long,
		Description: desc,
		Required:    false,
		HasArg:      hasArg,
		ArgName:     defaultArgName,
	}
}

// NewRequiredFlag constructs a new required flag.
func NewRequiredFlag(short rune, long string, desc string, hasArg bool) *Flag {
	return &Flag{
		Short:       short,
		Long:        long,
		Description: desc,
		Required:    true,
		HasArg:      hasArg,
		ArgName:     defaultArgName,
	}
}

// String returns a string representation of this flag.
func (f Flag) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString("cli.Flag{")
	if f.Short != 0 {
		buf.WriteString("Short='")
		buf.WriteRune(f.Short)
		buf.WriteString("', ")
	}
	if len(f.Long) != 0 {
		buf.WriteString("Long=\"")
		buf.WriteString(f.Long)
		buf.WriteString("\", ")
	}
	buf.WriteString("Description=\"")
	buf.WriteString(f.Description)
	buf.WriteString("\", ")
	buf.WriteString("Required=")
	buf.WriteString(strconv.FormatBool(f.Required))
	buf.WriteString(", HasArg=")
	buf.WriteString(strconv.FormatBool(f.HasArg))
	if f.HasArg {
		buf.WriteString(", ArgName=\"")
		buf.WriteString(f.ArgName)
		buf.WriteRune('"')
	}
	buf.WriteRune('}')

	return buf.String()
}

// FlagSlice represents a slice of Flag pointers, implements the sort interface.
type FlagSlice []*Flag

// Len is the number of elements in the collection.
func (p FlagSlice) Len() int {
	return len(p)
}

// Less reports whether the element with index i should sort before the element with index j.
func (p FlagSlice) Less(i, j int) bool {
	fi, fj := p[i], p[j]
	si, sj := fi.Short, fj.Short

	// Both have long flags, lexical sort
	if si == 0 && sj == 0 {
		return fi.Long < fj.Long
	}

	// Short flag sorts before long flag.
	if (si == 0) != (sj == 0) {
		return si != 0
	}
	il, jl := unicode.IsLower(si), unicode.IsLower(sj)

	// Same case, sort lexicographically.
	if il == jl {
		return si < sj
	}
	si2, sj2 := unicode.ToLower(si), unicode.ToLower(sj)

	// Same letter, lowercase sorts first.
	if si2 == sj2 {
		return il
	}

	// Different case, sort lexicographically.
	return si2 < sj2
}

// Swap swaps the elements with indexes i and j.
func (p FlagSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
