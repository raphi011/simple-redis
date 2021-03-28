package main

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct{input string; want interface{}}{
		{input: "+asd\r\n", want: "asd"},
		{input: "$-1\r\n", want: nil},
		{input: ":5\r\n", want: 5},
		{input: "$6\r\nfoobar\r\n", want: []byte("foobar")},
		{input: "-error\r\n", want: errors.New("error")},
		{input: "*2\r\n$3\r\nfoo\r\n$4\r\nbars\r\n", want: Array{[]byte("foo"), []byte("bars")}},
	}

	for _, tc := range tests {
		got, err := parse(bytes.NewBufferString(tc.input))

		if err != nil {
			t.Errorf("parse(%s) err: %v", tc.input, err)
		}

		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("parse(%s) expected: %v, got: %v", escape(tc.input), tc.want, got)
		}
	}
}

func TestToString(t *testing.T) {
	tests := []struct{input interface{}; want string}{
		{input: "asd", want: "+asd\r\n"},
		{input: nil, want: "$-1\r\n"},
		{input: 5, want: ":5\r\n"},
		{input: []byte("foobar"), want: "$6\r\nfoobar\r\n"},
		{input: errors.New("error"), want: "-error\r\n"},
		{input: Array{[]byte("foo"), []byte("bars")}, want: "*2\r\n$3\r\nfoo\r\n$4\r\nbars\r\n"},
	}

	for _, tc := range tests {
		buf := &strings.Builder{}
		err := toString(tc.input, buf)
		if err != nil {
			t.Errorf("toString(%s) err: %v", tc.input, err)
		}

		got := buf.String()

		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("toString(%s) expected: %v, got: %v", tc.input, escape(tc.want), escape(got))
		}
	}
}

func escape(input string) string {
	return strings.ReplaceAll(strings.ReplaceAll(input, "\r", "\\r"), "\n", "\\n")
}
