package cli

import "fmt"

type commandLine struct {
	args   []string // arguments
	flags  []*Flag  // parsed flags
	values map[*Flag]string
}

func (c *commandLine) addArg(arg string) {
	c.args = append(c.args, arg)
}

func (c *commandLine) processValue(flag *Flag, value string) error {
	if val, ok := c.values[flag]; ok {
		return fmt.Errorf(`%v already has a argument "%v"`, flag, val)
	}
	if !flag.HasArg {
		return fmt.Errorf("%v does not accept an argument", flag)
	}
	c.values[flag] = value
	return nil
}

func (c *commandLine) needsValue(flag *Flag) bool {
	if !flag.HasArg {
		return false
	}

	_, ok := c.values[flag]
	return !ok
}

func (c *commandLine) Value(flag *Flag) (string, bool) {
	if val, ok := c.values[flag]; ok {
		return val, ok
	}
	for _, f := range c.flags {
		if flag == f {
			return "", true
		}
	}
	return "", false
}

func (c *commandLine) Args() []string {
	return c.args
}

type CommandLine interface {
	// Value returns the value parsed for the specified flag.
	// If the flag was not parsed
	Value(flag *Flag) (string, bool)
	Args() []string
}
