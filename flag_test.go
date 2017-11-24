package cli

import (
	"reflect"
	"testing"
)

func TestNewFlag(t *testing.T) {
	want := &Flag{'a', "", "", "", false, false, defaultArgName}
	if got := NewFlag('a', "", "", false); !reflect.DeepEqual(got, want) {
		t.Errorf("NewFlag() = %v, want %v", got, want)
	}
}

func TestNewRequiredFlag(t *testing.T) {
	want := &Flag{'a', "", "", "", true, false, defaultArgName}
	if got := NewRequiredFlag('a', "", "", false); !reflect.DeepEqual(got, want) {
		t.Errorf("NewFlag() = %v, want %v", got, want)
	}
}

func TestFlag_String(t *testing.T) {
	type fields struct {
		Short       rune
		Long        string
		Description string
		Required    bool
		HasArg      bool
		ArgName     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"short, req=false, no arg",
			fields{'a', "", "desc", false, false, defaultArgName},
			`cli.Flag{Short='a', Description="desc", DefValue="aaa", Required=false, HasArg=false}`,
		},
		{
			"long, req=true, arg=BBB",
			fields{0, "bbb", "desc2", true, true, "BBB"},
			`cli.Flag{Long="bbb", Description="desc2", DefValue="b", Required=true, HasArg=true, ArgName="BBB"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Flag{
				Short:       tt.fields.Short,
				Long:        tt.fields.Long,
				Description: tt.fields.Description,
				Required:    tt.fields.Required,
				HasArg:      tt.fields.HasArg,
				ArgName:     tt.fields.ArgName,
			}
			if got := f.String(); got != tt.want {
				t.Errorf("Flag.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlagSlice_Len(t *testing.T) {
	want := 3
	if got := (FlagSlice{nil, nil, nil}).Len(); got != want {
		t.Errorf("FlagSlice.Len() = %v, want %v", got, want)
	}
}

func TestFlagSlice_Less(t *testing.T) {
	slice := FlagSlice{
		NewFlag(0, "aaa", "", false),
		NewFlag(0, "bbb", "", false),
		NewFlag('c', "", "", false),
		NewFlag('d', "", "", false),
		NewFlag('C', "", "", false),
	}

	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		p    FlagSlice
		args args
		want bool
	}{
		{"long and long: lexical", slice, args{0, 1}, true},
		{"short xor short: short first", slice, args{2, 0}, true},
		{"same case: lexical", slice, args{2, 3}, true},
		{"same letter: lower first", slice, args{2, 4}, true},
		{"diff case: lexical", slice, args{3, 4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("FlagSlice.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlagSlice_Swap(t *testing.T) {
	fi := NewFlag('a', "", "", false)
	fj := NewFlag('b', "", "", false)

	slice, want := FlagSlice{fi, fj}, FlagSlice{fj, fi}

	i, j := 0, 1
	if slice.Swap(i, j); !reflect.DeepEqual(slice, want) {
		t.Errorf("FlagSlice.Swap() = %v, want %v", slice, want)
	}
}
