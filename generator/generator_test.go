// generator/generator_test.go
package generator

import (
		"reflect" // Used for comparing maps/slices
		"strings"
		"testing"
)

// TestParseDigitRange uses a table-driven test to check various inputs.
func TestParseDigitRange(t *testing.T) {
	// Define test cases
	tests := []struct {
		name    string
		input   string
		wantMin int
		wantMax int
		wantErr bool
	}{
		{"SingleValid", "4", 4, 4, false},
		{"RangeValid", "3-6", 3, 6, false},
		{"InvalidString", "abc", 0, 0, true},
		{"InvalidRangeOrder", "6-3", 0, 0, true},
		{"InvalidRangeZero", "0-5", 0, 0, true},
		{"InvalidFormat", "3-6-9", 0, 0, true},
	}

		for _, tt := range tests {
		// t.Run allows running sub-tests, which gives clearer output
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax, err := parseDigitRange(tt.input)

			// Check if error presence (or absence) matches what we want
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDigitRange() error = %v, wantErr %v", err, tt.wantErr)
				return // Don't check values if error state is wrong
			}

			// If we expected no error, check if the values are correct
			if !tt.wantErr {
				if gotMin != tt.wantMin {
					t.Errorf("parseDigitRange() gotMin = %v, want %v", gotMin, tt.wantMin)
				}
				if gotMax != tt.wantMax {
					t.Errorf("parseDigitRange() gotMax = %v, want %v", gotMax, tt.wantMax)
				}
			}
		})
	}
}

// TestFilterWords checks if filtering (invalidChars) and de-duplication works.
func TestFilterWords(t *testing.T) {
	// We simulate loading one embedded string and one file
	input := []string{
		"Haus\nBoot\nYacht\nZebra", // Simulates defaultWordlist
		"Auto\nÄpfel\nHaus",       // Simulates a loaded file
	}

		// We expect duplicates ("Haus") to be removed and invalid chars ("Yacht", "Zebra", "Äpfel")
		// to be filtered out.
		want := map[string]bool{
		"Haus": true,
		"Boot": true,
		"Auto": true,
	}

		gotSlice := filterWords(input)

		// Convert the resulting slice to a map for easy comparison,
		// as the order of the slice is not guaranteed.
		gotMap := make(map[string]bool)
	for _, w := range gotSlice {
		gotMap[w] = true
	}

		if !reflect.DeepEqual(gotMap, want) {
		t.Errorf("filterWords() = %v, want %v", gotMap, want)
	}
}

// TestApplyCapitalization checks the different caps modes.
func TestApplyCapitalization(t *testing.T) {
	word := "test"

		t.Run("camel", func(t *testing.T) {
		want := "Test"
		got := applyCapitalization(word, "camel")
		if got != want {
			t.Errorf("applyCapitalization('camel') = %v, want %v", got, want)
		}
	})

		t.Run("none", func(t *testing.T) {
		want := "test"
		got := applyCapitalization(word, "none")
		if got != want {
			t.Errorf("applyCapitalization('none') = %v, want %v", got, want)
		}
	})

		t.Run("random", func(t *testing.T) {
		got := applyCapitalization(word, "random")
		if got == word {
			// This could technically fail by chance (if 't' is chosen), but it's unlikely
			t.Logf("Warning: 'random' capitalization might have returned the original word by chance.")
		}
		if len(got) != len(word) {
			t.Errorf("'random' capitalization changed word length")
		}
		if strings.ToLower(got) != word {
			t.Errorf("'random' capitalization changed word content: got %v", got)
		}
	})
}
