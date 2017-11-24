package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const (
	// Default number of characters per line.
	defaultWidth = 74

	// Default padding to the left of flag lines.
	defaultFlagPad = 2

	// Default padding before flag description.
	defaultDescPad = 2

	// Default usage statement prefix.
	defaultUsagePrefix = "Usage: "

	// Default prefix to the flags block.
	defaultFlagsPrefix = "Flags:"

	// Separator for short and long settings ie. -o, --opt.
	commaSeparator = ", "
)

// Formatter is a utility for formatting a help string for a FlagSet.
type Formatter struct {
	Width       int
	FlagPad     int
	DescPad     int
	UsagePrefix string
	FlagsPrefix string
}

// NewFormatter constructs a new Formatter with the default values.
func NewFormatter() *Formatter {
	f := &Formatter{
		Width:       defaultWidth,
		FlagPad:     defaultFlagPad,
		DescPad:     defaultDescPad,
		UsagePrefix: defaultUsagePrefix,
		FlagsPrefix: defaultFlagsPrefix,
	}
	return f
}

// PrintHelp prints a generated help message for the FlagSet to the Writer.
// Returns an error if the usage string is empty.
func (f *Formatter) PrintHelp(w io.Writer, usage string, header string, flags FlagSet, footer string) error {
	if len(usage) == 0 {
		return errors.New("cli.Formatter.PrintHelp: usage string not provided")
	}

	f.PrintUsage(w, usage)

	if len(header) > 0 {
		f.printWrapped(w, header)
	}

	f.PrintFlags(w, flags)

	if len(footer) > 0 {
		fmt.Fprintln(w)
		f.printWrapped(w, footer)
	}

	return nil
}

// PrintUsage prints a generated usage message for the FlagSet to the Writer.
func (f *Formatter) PrintUsage(w io.Writer, usage string) {
	// when wrapping, indent from the second argument in the usage
	argPos := strings.IndexRune(usage, ' ') + 1

	// use a buffer to join strings
	buf := bytes.NewBufferString(f.UsagePrefix)
	buf.WriteString(usage)

	f.printWrappedIndent(w, buf.String(), len(f.UsagePrefix)+argPos)
}

// PrintFlags prints a generated message detailing the flags in the FlagSet to the Writer.
func (f *Formatter) PrintFlags(w io.Writer, flags FlagSet) {
	buf := new(bytes.Buffer)

	f.renderFlags(buf, flags)
	fmt.Fprint(w, buf.String())
}

// printWrappedIndent prints text to the Writer with line wrapping using Formatter width.
func (f *Formatter) printWrapped(w io.Writer, text string) {
	f.printWrappedIndent(w, text, 0)
}

// printWrappedIndent prints text to the Writer with line wrapping using Formatter width.
// 'newLineIndent' is the amount to indent each wrapped new line by
func (f *Formatter) printWrappedIndent(w io.Writer, text string, newLineIndent int) {
	buf := new(bytes.Buffer)

	f.renderWrappedTextBlock(buf, text, newLineIndent)
	fmt.Fprintln(w, buf.String())
}

func (f *Formatter) renderFlags(buf *bytes.Buffer, fs FlagSet) *bytes.Buffer {
	flagPad := createPad(f.FlagPad)              // padding before short flag
	optPad := createPad(2 + len(commaSeparator)) // padding to fill if no short flag is present
	descPad := createPad(f.DescPad)              // padding before description

	// The length of the longest flag ie. len("  -o, --opt=value")=17.
	// Used in description alignment.
	maxLen := 0

	flags := fs.Flags()
	info := make([]string, 0, len(flags))

	if len(flags) == 0 {
		return buf
	}

	buf.WriteByte('\n')
	buf.WriteString(f.FlagsPrefix)
	buf.WriteByte('\n')

	for _, flag := range flags {
		fBuf := new(bytes.Buffer)

		fBuf.WriteString(flagPad) // add initial padding

		if flag.Short == 0 && len(flag.Long) == 0 {
			panic(fmt.Sprintf("cli.renderFlags: %s has no short or long option", flag))
		}

		if flag.Short != 0 {
			fBuf.WriteString(ShortPrefix)
			fBuf.WriteRune(flag.Short)

			if len(flag.Long) > 1 {
				// Add separator if long option exists.
				fBuf.WriteString(commaSeparator)
			}
		} else {
			// No short option add padding to align long option.
			fBuf.WriteString(optPad)
		}

		if len(flag.Long) > 1 {
			fBuf.WriteString(LongPrefix)
			fBuf.WriteString(flag.Long)
		}

		if flag.HasArg {
			sep := ValueSeparator
			if len(flag.Long) == 0 {
				sep = ' '
			}
			fBuf.WriteRune(sep)
			fBuf.WriteString(flag.ArgName)
		}

		s := fBuf.String()
		if len(s) > maxLen && len(s) < f.Width*2/5 {
			maxLen = len(s)
		}

		info = append(info, s)
	}

	for i, flag := range flags {
		if i > 0 {
			buf.WriteByte('\n')
		}

		s := info[i]
		fBuf := bytes.NewBufferString(s)

		if len(s) < maxLen {
			fBuf.WriteString(createPad(maxLen - len(s)))
		}

		// Special conditions if the description is on a new line from the flag.
		if len(s) > maxLen {
			fBuf.WriteByte('\n')
			fBuf.WriteString(createPad(maxLen + f.DescPad))
		} else {
			fBuf.WriteString(descPad)
		}

		newLineIndent := maxLen + f.DescPad*2

		if len(flag.Description) > 0 {
			fBuf.WriteString(flag.Description)
		}

		f.renderWrappedText(buf, fBuf.String(), newLineIndent)
	}
	buf.WriteByte('\n')

	return buf
}

func (f *Formatter) renderWrappedText(buf *bytes.Buffer, text string, newLineIndent int) *bytes.Buffer {
	pos := findWrapPos(text, f.Width)

	// No need to wrap.
	if pos == -1 {
		buf.WriteString(strings.TrimRightFunc(text, unicode.IsSpace))
		return buf
	}

	// Write first line up to wrap position.
	buf.WriteString(strings.TrimRightFunc(text[:pos], unicode.IsSpace))
	buf.WriteByte('\n')

	// Stop infinite looping.
	if f.Width <= newLineIndent {
		newLineIndent = 1
	}

	pad := createPad(newLineIndent)
	buf2 := new(bytes.Buffer)

	for {
		buf2.Reset()
		buf2.WriteString(pad)
		buf2.WriteString(strings.TrimSpace(text[pos:]))

		text = buf2.String()
		pos = findWrapPos(text, f.Width)

		if pos == -1 {
			buf.WriteString(text)
			return buf
		}

		if f.Width < len(text) && pos == newLineIndent-1 {
			pos = f.Width
		}

		buf.WriteString(strings.TrimRightFunc(text[:pos], unicode.IsSpace))
		buf.WriteByte('\n')
	}
}

func (f *Formatter) renderWrappedTextBlock(buf *bytes.Buffer, text string, newLineIndent int) *bytes.Buffer {
	for i, line := range strings.Split(text, "\n") {
		if i > 0 {
			buf.WriteByte('\n')
		}
		f.renderWrappedText(buf, line, newLineIndent)
	}

	return buf
}

// findWrapPos returns the wrap position for the text in 'text' with column width 'width'.
// Wraps at the last whitespace character before 'width', if there is no whitespace returns 'width'.
func findWrapPos(text string, width int) int {
	// Start from each new line if it ends before width.
	end := width
	nl := strings.IndexByte(text, '\n')
	if nl != -1 && nl < width {
		end += nl
	}

	// Doesn't need wrapping.
	if len(text) <= end {
		return -1
	}

	// Search for last whitespace character.
	for pos, runes := end, []rune(text); pos >= 0; pos-- {
		c := runes[pos]
		if unicode.IsSpace(c) {
			return pos
		}
	}

	// Couldn't find any whitespace to wrap at, just gonna wrap at max width.
	return width
}

// createPad returns a string length 'length' of spaces.
func createPad(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = ' '
	}
	return string(b)
}
