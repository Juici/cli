package cli

import (
	"reflect"
	"testing"
)

func Test_commandLine_addArg(t *testing.T) {
	want := &commandLine{
		args:   []string{"arg1", "arg2", "want"},
		flags:  []*Flag(nil),
		values: map[*Flag]string(nil),
	}

	c := &commandLine{
		args:   []string{"arg1", "arg2"},
		flags:  []*Flag(nil),
		values: map[*Flag]string(nil),
	}
	if c.addArg("arg3"); reflect.DeepEqual(c, want) {
		t.Errorf("commandLine.addArg() = %v, want %v", c, want)
	}
}

func Test_commandLine_Args(t *testing.T) {
	want := []string{"arg1", "arg2", "want"}
	c := &commandLine{
		args:   want,
		flags:  []*Flag(nil),
		values: map[*Flag]string(nil),
	}
	if got := c.Args(); !reflect.DeepEqual(got, want) {
		t.Errorf("commandLine.Args() = %v, want %v", got, want)
	}
}
