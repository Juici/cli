package cli

import (
	"reflect"
	"testing"
)

func TestNewFlagSet(t *testing.T) {
	want := &FlagSet{
		shorts:   make(map[rune]*Flag),
		longs:    make(map[string]*Flag),
		required: make([]*Flag, 0),
	}
	if got := NewFlagSet(); !reflect.DeepEqual(got, want) {
		t.Errorf("NewFlagSet() = %v, want %v", got, want)
	}
}

func TestFlagSet_AddFlag(t *testing.T) {
	fs := NewFlagSet()
	fl := NewFlagSet()

	type fields struct {
		shorts   map[rune]*Flag
		longs    map[string]*Flag
		required []*Flag
	}
	type args struct {
		flag *Flag
	}
	tests := []struct {
		name    string
		f       *FlagSet
		args    args
		wantErr bool
	}{
		{
			"short non-letter",
			NewFlagSet(),
			args{NewFlag('5', "", "", false)},
			true,
		},
		{
			"valid short",
			fs,
			args{NewFlag('s', "", "", false)},
			false,
		},
		{
			"short already exists",
			fs,
			args{NewFlag('s', "", "", false)},
			true,
		},
		{
			"long non-letter",
			NewFlagSet(),
			args{NewFlag(0, "l0ng", "", false)},
			true,
		},
		{
			"valid long required",
			fl,
			args{NewRequiredFlag(0, "long", "", false)},
			false,
		},
		{
			"long already exists",
			fl,
			args{NewFlag(0, "long", "", false)},
			true,
		},
		{
			"long too short",
			NewFlagSet(),
			args{NewFlag(0, "long"[:MinimumLongFlagLength-1], "", false)},
			true,
		},
		{
			"no flag",
			NewFlagSet(),
			args{NewFlag(0, "", "", false)},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.AddFlag(tt.args.flag); (err != nil) != tt.wantErr {
				t.Errorf("FlagSet.AddFlag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFlagSet_AddNewFlag(t *testing.T) {
	fs := NewFlagSet()
	fl := NewFlagSet()

	type args struct {
		short  rune
		long   string
		desc   string
		hasArg bool
	}
	tests := []struct {
		name    string
		f       *FlagSet
		args    args
		want    *Flag
		wantErr bool
	}{
		{
			"short non-letter",
			NewFlagSet(),
			args{'5', "", "", false},
			nil,
			true,
		},
		{
			"valid short",
			fs,
			args{'s', "", "", false},
			NewFlag('s', "", "", false),
			false,
		},
		{
			"short already exists",
			fs,
			args{'s', "", "", false},
			nil,
			true,
		},
		{
			"long non-letter",
			NewFlagSet(),
			args{0, "l0ng", "", false},
			nil,
			true,
		},
		{
			"valid long",
			fl,
			args{0, "long", "", false},
			NewFlag(0, "long", "", false),
			false,
		},
		{
			"long already exists",
			fl,
			args{0, "long", "", false},
			nil,
			true,
		},
		{
			"long too short",
			NewFlagSet(),
			args{0, "long"[:MinimumLongFlagLength-1], "", false},
			nil,
			true,
		},
		{
			"no flag",
			NewFlagSet(),
			args{0, "", "", false},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.AddNewFlag(tt.args.short, tt.args.long, tt.args.desc, tt.args.hasArg)
			if (err != nil) != tt.wantErr {
				t.Errorf("FlagSet.AddNewFlag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlagSet.AddNewFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlagSet_AddNewRequiredFlag(t *testing.T) {
	fs := NewFlagSet()
	fl := NewFlagSet()

	type args struct {
		short  rune
		long   string
		desc   string
		hasArg bool
	}
	tests := []struct {
		name    string
		f       *FlagSet
		args    args
		want    *Flag
		wantErr bool
	}{
		{
			"short non-letter",
			NewFlagSet(),
			args{'5', "", "", false},
			nil,
			true,
		},
		{
			"valid short",
			fs,
			args{'s', "", "", false},
			NewRequiredFlag('s', "", "", false),
			false,
		},
		{
			"short already exists",
			fs,
			args{'s', "", "", false},
			nil,
			true,
		},
		{
			"long non-letter",
			NewFlagSet(),
			args{0, "l0ng", "", false},
			nil,
			true,
		},
		{
			"valid long",
			fl,
			args{0, "long", "", false},
			NewRequiredFlag(0, "long", "", false),
			false,
		},
		{
			"long already exists",
			fl,
			args{0, "long", "", false},
			nil,
			true,
		},
		{
			"long too short",
			NewFlagSet(),
			args{0, "long"[:MinimumLongFlagLength-1], "", false},
			nil,
			true,
		},
		{
			"no flag",
			NewFlagSet(),
			args{0, "", "", false},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.AddNewRequiredFlag(tt.args.short, tt.args.long, tt.args.desc, tt.args.hasArg)
			if (err != nil) != tt.wantErr {
				t.Errorf("FlagSet.AddNewRequiredFlag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlagSet.AddNewRequiredFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlagSet_Flags(t *testing.T) {
	f := NewFlagSet()
	fa, _ := f.AddNewFlag('a', "aaa", "", false)
	fb, _ := f.AddNewFlag(0, "bbb", "", false)
	fc, _ := f.AddNewFlag('C', "ccc", "", false)
	fd, _ := f.AddNewFlag(0, "ddd", "", false)

	want := []*Flag{fa, fc, fb, fd}
	if got := f.Flags(); !reflect.DeepEqual(got, want) {
		t.Errorf("FlagSet.Flags() = %v, want %v", got, want)
	}
}

func TestFlagSet_ShortFlags(t *testing.T) {
	f := NewFlagSet()
	fa, _ := f.AddNewFlag('a', "aaa", "", false)
	f.AddNewFlag(0, "bbb", "", false)
	fc, _ := f.AddNewFlag('C', "ccc", "", false)
	f.AddNewFlag(0, "ddd", "", false)

	want := []*Flag{fa, fc}
	if got := f.ShortFlags(); !reflect.DeepEqual(got, want) {
		t.Errorf("FlagSet.ShortFlags() = %v, want %v", got, want)
	}
}

func TestFlagSet_LongFlags(t *testing.T) {
	f := NewFlagSet()
	fa, _ := f.AddNewFlag('a', "aaa", "", false)
	fb, _ := f.AddNewFlag(0, "bbb", "", false)
	fc, _ := f.AddNewFlag('C', "ccc", "", false)
	fd, _ := f.AddNewFlag(0, "ddd", "", false)

	want := []*Flag{fa, fc, fb, fd}
	if got := f.LongFlags(); !reflect.DeepEqual(got, want) {
		t.Errorf("FlagSet.LongFlags() = %v, want %v", got, want)
	}
}

func TestFlagSet_RequiredFlags(t *testing.T) {
	f := NewFlagSet()
	f.AddNewFlag('a', "aaa", "", false)
	fb, _ := f.AddNewRequiredFlag(0, "bbb", "", false)
	fc, _ := f.AddNewRequiredFlag('C', "ccc", "", false)
	f.AddNewFlag(0, "ddd", "", false)

	want := []*Flag{fc, fb}
	if got := f.RequiredFlags(); !reflect.DeepEqual(got, want) {
		t.Errorf("FlagSet.RequiredFlags() = %v, want %v", got, want)
	}
}

func TestFlagSet_Lookup(t *testing.T) {
	fa := NewFlag('a', "aaa", "", false)

	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  *Flag
		want1 bool
	}{
		{"find short", args{"a"}, fa, true},
		{"no short", args{"b"}, nil, false},
		{"find long", args{"aaa"}, fa, true},
		{"no long", args{"bbb"}, nil, false},
	}
	for _, tt := range tests {
		f := NewFlagSet()
		f.AddFlag(fa)

		t.Run(tt.name, func(t *testing.T) {
			got, got1 := f.Lookup(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlagSet.Lookup() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FlagSet.Lookup() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestFlagSet_Matches(t *testing.T) {
	f1 := NewFlag('a', "aaa", "", false)
	f2 := NewFlag(0, "bbb", "", false)
	f3 := NewFlag(0, "aab", "", false)

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"perfect match short", args{"a"}, []string{"a"}},
		{"perfect match long case insensitive", args{"AAA"}, []string{"aaa"}},
		{"partial matches case insensitive", args{"AA"}, []string{"aaa", "aab"}},
		{"no matches", args{"ccc"}, []string(nil)},
	}
	for _, tt := range tests {
		f := NewFlagSet()
		f.AddFlag(f1)
		f.AddFlag(f2)
		f.AddFlag(f3)

		t.Run(tt.name, func(t *testing.T) {
			if got := f.Matches(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlagSet.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
