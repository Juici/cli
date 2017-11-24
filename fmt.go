package cli

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
