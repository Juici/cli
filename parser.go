package cli

const (
	defaultShortPrefix = "-"
	defaultLongPrefix  = "--"
)

// ParserOptions represents the options used by a Parser when parsing command line arguments.
type ParserOptions struct {
	ShortOptPrefix string // the prefix for a short option
	LongOptPrefix  string // the prefix for a long option
}

// DefaultParserOptions returns the default options for use by a Parser.
func DefaultParserOptions() *ParserOptions {
	o := &ParserOptions{
		ShortOptPrefix: defaultShortPrefix,
		LongOptPrefix:  defaultLongPrefix,
	}
	return o
}

// Parser represents a command line argument parser.
type Parser struct {
	Options *ParserOptions // the parser options

	curFlag *Flag
}

// NewParser returns a new parser with the specified ParserOptions.
// If opts is nil then the Parser will use DefaultParserOptions.
func NewParser(opts *ParserOptions) *Parser {
	if opts == nil {
		opts = DefaultParserOptions()
	}
	p := &Parser{Options: opts}
	return p
}

func (p *Parser) Parse(args []string) {

}
