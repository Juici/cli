package cli

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// FlagSet represents a collection of Flags to be parsed by a Parser.
type FlagSet struct {
	shorts   map[rune]*Flag
	longs    map[string]*Flag
	required []*Flag
}

// NewFlagSet constructs and returns a new empty FlagSet.
func NewFlagSet() *FlagSet {
	return &FlagSet{
		shorts:   make(map[rune]*Flag),
		longs:    make(map[string]*Flag),
		required: make([]*Flag, 0),
	}
}

// AddFlag adds the specified Flag to the FlagSet.
// Returns an error if the short/long flag is invalid or already exists.
func (f *FlagSet) AddFlag(flag *Flag) error {
	s, l := false, false

	// Make sure neither flag is invalid before adding either.
	if flag.Short != 0 {
		if !unicode.IsLetter(flag.Short) {
			return fmt.Errorf("cli.FlagSet: short flag '%v' is not a letter", flag.Short)
		}
		if _, ok := f.shorts[flag.Short]; ok {
			return fmt.Errorf("cli.FlagSet: short flag '%v' already exists", flag.Short)
		}

		s = true
	}
	if len(flag.Long) != 0 {
		flag.Long = strings.ToLower(flag.Long)

		if i := strings.IndexFunc(flag.Long, func(r rune) bool {
			return !unicode.IsLetter(r)
		}); i > -1 {
			return fmt.Errorf(`cli.FlagSet: long flag "%v" contains a non letter '%v'`,
				flag.Long, []rune(flag.Long)[i])
		}
		if _, ok := f.longs[flag.Long]; ok {
			return fmt.Errorf(`cli.FlagSet: long flag "%v" already exists`, flag.Long)
		}
		if len([]rune(flag.Long)) < MinimumLongFlagLength {
			return fmt.Errorf(`cli.FlagSet: long flag "%v" must be %d or more letters`, flag.Long, MinimumLongFlagLength)
		}

		l = true
	}
	// Flag must have a short or long variation, or both.
	if !(s || l) {
		return errors.New("cli.FlagSet: no short or long flag specified")
	}

	if s {
		f.shorts[flag.Short] = flag
	}
	if l {
		f.longs[flag.Long] = flag
	}
	if flag.Required {
		f.required = append(f.required, flag)
	}

	return nil
}

// AddNewFlag creates a new Flag and adds it to the FlagSet.
// Returns the created Flag, or an error if the short/long flag is invalid or already exists.
func (f *FlagSet) AddNewFlag(short rune, long string, desc string, hasArg bool) (*Flag, error) {
	flag := NewFlag(short, long, desc, hasArg)

	err := f.AddFlag(flag)
	if err != nil {
		return nil, err
	}

	return flag, nil
}

// AddNewFlag creates a new required Flag and adds it to the FlagSet.
// Returns the created Flag, or an error if the short/long flag is invalid or already exists.
func (f *FlagSet) AddNewRequiredFlag(short rune, long string, desc string, hasArg bool) (*Flag, error) {
	flag := NewRequiredFlag(short, long, desc, hasArg)

	err := f.AddFlag(flag)
	if err != nil {
		return nil, err
	}

	return flag, nil
}

// Flags returns a slice with all the Flags in this FlagSet.
// Slice is sorted by FlagSlice.Sort()
func (f *FlagSet) Flags() []*Flag {
	minCap := len(f.shorts)
	if l := len(f.longs); l > minCap {
		minCap = l
	}
	flags := make([]*Flag, 0, minCap)

	for _, flag := range f.shorts {
		flags = append(flags, flag)
	}
	for _, flag := range f.longs {
		if flag.Short != 0 {
			continue
		}
		flags = append(flags, flag)
	}

	sort.Sort(FlagSlice(flags))
	return flags
}

// ShortFlags returns a slice with all the short Flags in this FlagSet, sorted by FlagSlice.Sort().
func (f *FlagSet) ShortFlags() []*Flag {
	flags := make([]*Flag, 0, len(f.shorts))
	for _, flag := range f.shorts {
		flags = append(flags, flag)
	}
	sort.Sort(FlagSlice(flags))
	return flags
}

// ShortFlags returns a slice with all the long Flags in this FlagSet, sorted by FlagSlice.Sort().
func (f *FlagSet) LongFlags() []*Flag {
	flags := make([]*Flag, 0, len(f.longs))
	for _, flag := range f.longs {
		flags = append(flags, flag)
	}
	sort.Sort(FlagSlice(flags))
	return flags
}

// ShortFlags returns a slice with all the required Flags in this FlagSet, sorted by FlagSlice.Sort().
func (f *FlagSet) RequiredFlags() []*Flag {
	flags := make([]*Flag, len(f.required))
	copy(flags, f.required)
	sort.Sort(FlagSlice(flags))
	return flags
}

// Lookup returns the Flag with the long or short name specified.
// Returns (Flag, true) if a matching Flag was found.
// Or (nil, false) if no matching Flag was found.
func (f *FlagSet) Lookup(name string) (*Flag, bool) {
	if runes := []rune(name); len(runes) == 1 {
		s := runes[0]
		if flag, ok := f.shorts[s]; ok {
			return flag, ok
		}
		return nil, false
	}

	name = strings.ToLower(name)
	if flag, ok := f.longs[name]; ok {
		return flag, ok
	}

	return nil, false
}

// Matches returns long flags starting with the name specified.
// Casing is ignored for long flags and the slice will be sorted lexicographically.
// If a perfect match is found only that match will be returned.
func (f *FlagSet) Matches(name string) []string {
	var ret []string

	if _, ok := f.Lookup(name); ok {
		if len([]rune(name)) > 1 {
			// Lowercase if long flag.
			name = strings.ToLower(name)
		}

		ret = append(ret, name)
		return ret
	}

	name = strings.ToLower(name)
	for _, flag := range f.longs {
		if strings.HasPrefix(flag.Long, name) {
			ret = append(ret, flag.Long)
		}
	}
	return ret
}
