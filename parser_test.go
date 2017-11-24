package cli

import (
	"reflect"
	"testing"
)

func TestNewParser(t *testing.T) {
	want := &Parser{
		cmd:         nil,
		flags:       nil,
		expected:    nil,
		skipParsing: false,
		curFlag:     nil,
		curToken:    "",
	}
	if got := NewParser(); !reflect.DeepEqual(got, want) {
		t.Errorf("NewParser() = %v, want %v", got, want)
	}
}

func TestParser_isNegativeNumber(t *testing.T) {
	type fields struct {
		cmd         *commandLine
		flags       *FlagSet
		expected    []*Flag
		skipParsing bool
		curFlag     *Flag
		curToken    string
	}
	type args struct {
		token string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"negative number",
			fields{nil, nil, nil, false, nil, ""},
			args{"-5"},
			true,
		},
		{
			"positive number",
			fields{nil, nil, nil, false, nil, ""},
			args{"5"},
			false,
		},
		{
			"short flag",
			fields{nil, nil, nil, false, nil, ""},
			args{string(append([]byte(ShortPrefix), "s"...))},
			false,
		},
		{
			"long flag",
			fields{nil, nil, nil, false, nil, ""},
			args{string(append([]byte(LongPrefix), "long"...))},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				cmd:         tt.fields.cmd,
				flags:       tt.fields.flags,
				expected:    tt.fields.expected,
				skipParsing: tt.fields.skipParsing,
				curFlag:     tt.fields.curFlag,
				curToken:    tt.fields.curToken,
			}
			if got := p.isNegativeNumber(tt.args.token); got != tt.want {
				t.Errorf("Parser.isNegativeNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
