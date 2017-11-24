package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	ShortPrefix    = "-"
	LongPrefix     = "--"
	ValueSeparator = '='
)

// Parser represents a command line argument parser.
type Parser struct {
	cmd      *commandLine // the command-line instance
	flags    *FlagSet     // the flags being parsed against
	expected []*Flag      // the expected flags

	skipParsing bool   // true if no more flags should be parsed
	curFlag     *Flag  // the last flag parsed
	curToken    string // the token currently being parsed
}

// NewParser returns a new parser.
func NewParser() *Parser {
	p := &Parser{
		cmd:         nil,
		flags:       nil,
		expected:    nil,
		skipParsing: false,
		curFlag:     nil,
		curToken:    "",
	}
	return p
}

// Parse parses the command-line arguments passed when executing the program.
func (p *Parser) Parse(flags *FlagSet) (CommandLine, error) {
	cmd, err := p.ParseArgs(flags, os.Args[1:])
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// ParseArgs parses the specified slice of string arguments.
func (p *Parser) ParseArgs(flags *FlagSet, args []string) (CommandLine, error) {
	p.cmd = &commandLine{
		flags: make([]*Flag, 0),
		args:  make([]string, 0),
	}
	p.flags = flags

	p.expected = make([]*Flag, len(flags.required))
	copy(p.expected, flags.required)

	p.skipParsing = false
	p.curFlag = nil

	if args != nil {
		for _, token := range args {
			err := p.handleToken(token)
			if err != nil {
				return nil, err
			}
		}
	}

	if p.curFlag != nil && p.curFlag.HasArg {
		return nil, fmt.Errorf("missing argument for %v", p.curFlag)
	}
	if len(p.expected) > 0 {
		return nil, fmt.Errorf("missing required flags %v", p.expected)
	}

	return p.cmd, nil
}

func (p *Parser) handleToken(token string) error {
	p.curToken = token

	var err error

	switch {
	case p.skipParsing:
		p.cmd.addArg(token)
		break
	case token == LongPrefix || token == ShortPrefix:
		p.skipParsing = true
		break
	case p.curFlag != nil && p.curFlag.HasArg && p.isArg(token):
		err = p.cmd.processValue(p.curFlag, token)
		break
	case strings.HasPrefix(token, LongPrefix):
		err = p.handleLong(token)
		break
	case strings.HasPrefix(token, ShortPrefix) && token != ShortPrefix:
		err = p.handleShort(token)
		break
	default:
		err = p.handleUnknown(token)
	}

	if err != nil {
		return err
	}

	if p.curFlag != nil && !p.cmd.needsValue(p.curFlag) {
		p.curFlag = nil
	}

	return nil
}

func (p *Parser) handleLong(token string) error {
	// Long flags are lowercase.
	long := strings.ToLower(token)
	long = long[len(LongPrefix):]

	i := strings.IndexRune(long, ValueSeparator)
	if i == -1 {
		flag, ok := p.flags.longs[long]
		if !ok {
			return p.handleUnknown(token)
		}
		return p.handleFlag(flag)
	}

	if len(long) <= i+1 {
		return fmt.Errorf(`no value found for "%v" after '%v'`, token, ValueSeparator)
	}

	long, val := long[:i], long[i+1:]
	flag, ok := p.flags.longs[long]
	if !ok || !flag.HasArg {
		return p.handleUnknown(token)
	}

	err := p.handleFlag(flag)
	if err != nil {
		return err
	}

	// Strip leading and trailing ""
	if runes := []rune(long); runes[0] == '"' && runes[len(runes)-1] == '"' {
		long = string(runes[1 : len(runes)-2])
	}

	err = p.cmd.processValue(p.curFlag, val)
	if err != nil {
		return err
	}
	p.curFlag = nil

	return nil
}

func (p *Parser) handleShort(token string) error {
	// Flags are lowercase.
	short := token[len(ShortPrefix):]

	// Handles single short flag and concatenated short flags.
	runes := []rune(short)
	for i, short := range runes {
		flag, ok := p.flags.shorts[short]
		if !ok {
			if i > 0 {
				token = string(runes[i:])
			}
			return p.handleUnknown(token)
		}

		err := p.handleFlag(flag)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) handleFlag(flag *Flag) error {

}

func (p *Parser) handleUnknown(token string) error {
	if strings.HasPrefix(token, LongPrefix) && len(token) > len(LongPrefix) {
		return fmt.Errorf(`unrecognised flag "%v"`, token)
	}
	if strings.HasPrefix(token, ShortPrefix) && len(token) > len(ShortPrefix) {
		return fmt.Errorf(`unrecognised flag "%v"`, token)
	}

	p.cmd.addArg(token)
	p.skipParsing = true

	return nil
}

func (p *Parser) isArg(token string) bool {
	return !p.isFlag(token) || p.isNegativeNumber(token)
}

func (p *Parser) isFlag(token string) bool {
	return p.isShort(token) || p.isLong(token)
}

func (p *Parser) isShort(token string) bool {
	if len(token) <= 1 || !strings.HasPrefix(token, ShortPrefix) {
		return false
	}

	// Flags are lowercase.
	token = strings.ToLower(token)

	// Strip leading prefix.
	token = token[len(ShortPrefix):]

	// Strip trailing =... (if it exists).
	i := strings.IndexRune(token, ValueSeparator)
	if i != -1 {
		token = token[:i]
	}

	// Check if first rune in string is a short option.
	// Valid for a single short flag or concatenated short flags.
	// Does not check that other flags in the concatenated flag string are valid.
	_, ok := p.flags.shorts[[]rune(token)[0]]
	return ok
}

func (p *Parser) isLong(token string) bool {
	if len(token) == 0 || !strings.HasPrefix(token, LongPrefix) {
		return false
	}

	// Flags are lowercase.
	token = strings.ToLower(token)

	// Strip leading prefix.
	token = token[len(LongPrefix):]

	// Strip trailing =... (if it exists).
	i := strings.IndexRune(token, ValueSeparator)
	if i != -1 {
		token = token[:i]
	}

	// Check if token is a long option.
	_, ok := p.flags.longs[token]
	return ok
}

func (p *Parser) isNegativeNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}
