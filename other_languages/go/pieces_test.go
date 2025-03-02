package main

import "testing"

func TestNewPiece(t *testing.T) {
	for _, tc := range []struct {
		name           string
		matrix         [][]string
		expectedBitmap uint64
	}{
		{
			name: "1",
			matrix: [][]string{
				{"{1}", ""},
				{"{1}", "{1}"},
				{"{1}", "{1}"},
			},
			expectedBitmap: 0b111101,
		},
		{
			name: "2",
			matrix: [][]string{
				{"{2}", "{2}", "{2}", "{2}"},
				{"{2}", "", "", ""},
			},
			expectedBitmap: 0b00011111,
		},
		{
			name: "3",
			matrix: [][]string{
				{"{3}", "{3}", "{3}", ""},
				{"", "", "{3}", "{3}"},
			},
			expectedBitmap: 0b11000111,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := newPiece(tc.matrix)

			if p.bitmap != tc.expectedBitmap {
				t.Fatalf("not equal expected: %b, got: %b", tc.expectedBitmap, p.bitmap)
			}
			t.Log(p)
		})
	}

}
