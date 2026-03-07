package main

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		err      bool
	}{
		{
			name:     "1.pos",
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
			err:      false,
		},
		{
			name:     "2.pos",
			input:    "abcd",
			expected: "abcd",
			err:      false,
		},
		{
			name:     "3.pos",
			input:    "",
			expected: "",
			err:      false,
		},
		{
			name:     "4.neg",
			input:    "45",
			expected: "",
			err:      true,
		},
		{
			name:     "5.pos",
			input:    "qwe\\4\\5",
			expected: "qwe45",
			err:      false,
		},
		{
			name:     "6.pos",
			input:    "qwe\\45",
			expected: "qwe44444",
			err:      false,
		},
		{
			name:     "7.neg",
			input:    "5abc",
			expected: "",
			err:      true,
		},
		{
			name:     "8.neg",
			input:    "abc\\",
			expected: "",
			err:      true,
		},
		{
			name:     "9.pos",
			input:    "a10b2",
			expected: "aaaaaaaaaabb",
			err:      false,
		},
		{
			name:     "10.pos",
			input:    "Го5ло4ва3нё2в",
			expected: "Гооооолоооовааанёёв",
			err:      false,
		},
		{
			name:     "11.pos",
			input:    "\\В\\а\\л\\е\\н\\т\\и\\н",
			expected: "Валентин",
			err:      false,
		},
		{
			name:     "12.pos",
			input:    "a\\4b\\5c",
			expected: "a4b5c",
			err:      false,
		},
		{
			name:     "13.pos",
			input:    "a0b",
			expected: "b",
			err:      false,
		},
		{
			name:     "14.pos",
			input:    "\\\\\\\\",
			expected: "\\\\",
			err:      false,
		},
		{
			name:     "15.pos",
			input:    "\\4\\5\\6",
			expected: "456",
			err:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := unpackString(tt.input)

			if (err != nil) != tt.err {
				t.Errorf("unpackString() error = %v, wantErr %v", err, tt.err)
				return
			}

			if !tt.err && result != tt.expected {
				t.Errorf("unpackString() = %s, want %s", result, tt.expected)
			}
		})
	}
}
