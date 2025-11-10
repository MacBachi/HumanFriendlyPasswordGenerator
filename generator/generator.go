// generator/generator.go
package generator

import (
	"bufio"
	"crypto/rand"
	_ "embed" // Required for go:embed, blank import to avoid unused error
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// 🚨 CLEANED EMBED BLOCK
//go:embed wordlist.txt
var defaultWordlist string // Für eine einzelne Datei als string korrekt

// Beispiel für mehrere Dateien:
// //go:embed *.txt
// var wordlists embed.FS

// --- Configuration ---
const invalidChars = "yzäöüß"
var defaultSpecialChars = []string{"+", "-", "_", "!"}
var randomSource = rand.Reader

// Config holds all settings for the generator
type Config struct {
	WordCount      int
	SeparatorCount int
	SeparatorPool  string
	DigitRange     string
	TypoRate       float64
	AltList        string
	MergeList      string
	CapsMode       string
	Verbose        bool // <--- NEU: für -v/verify/verbose
}

// Generator holds the prepared state (the filtered wordlist)
type Generator struct {
	filteredWords []string
	config        Config
}

// NewGenerator creates a new generator instance, loads, and filters words
func NewGenerator(config Config) (*Generator, error) {
	// 1. Determine which words to load
	var stringsToParse []string
	var filesToLoad []string

	if config.AltList != "" {
		// Use *only* the alternative list
		filesToLoad = append(filesToLoad, config.AltList)
	} else {
		// Use the default embedded list
		stringsToParse = append(stringsToParse, defaultWordlist)
		if config.MergeList != "" {
			// And merge the extra file
			filesToLoad = append(filesToLoad, config.MergeList)
		}
	}

	// 2. Load words from files
	fileWords, err := loadWordsFromFiles(filesToLoad)
	if err != nil {
		return nil, err
	}
	stringsToParse = append(stringsToParse, fileWords...)

	// 3. Filter and de-duplicate all words
	filteredWords := filterWords(stringsToParse)
	
	if len(filteredWords) < config.WordCount {
		return nil, fmt.Errorf("not enough valid words found (%d) to generate %d words", len(filteredWords), config.WordCount)
	}

	gen := &Generator{
		filteredWords: filteredWords,
		config:        config,
	}
	return gen, nil
}

// WordCount returns the number of valid words loaded
func (g *Generator) WordCount() int {
	return len(g.filteredWords)
}

// Generate creates a single new passphrase
func (g *Generator) Generate() (string, error) {
	typoApplied := false
	var wordChunks []string
	
	// 1. Parse config options needed for this run
	minDigits, maxDigits, err := parseDigitRange(g.config.DigitRange)
	if err != nil {
		return "", err
	}
	separatorPool := parseSeparatorPool(g.config.SeparatorPool)
	if len(separatorPool) == 0 {
		return "", fmt.Errorf("separator pool is empty")
	}

	// 2. Generate all word chunks
	for i := 0; i < g.config.WordCount; i++ {
		randomWord, err := g.getRandomWord()
		if err != nil {
			return "", fmt.Errorf("failed to get random word: %w", err)
		}

		// Apply Typo Logic
		if !typoApplied && g.config.TypoRate > 0 {
			if shouldApplyTypo(g.config.TypoRate) || i == g.config.WordCount-1 {
				randomWord = applyTypoTransposition(randomWord)
				typoApplied = true
			}
		}
		
		// Apply NEW Capitalization Logic
		wordChunks = append(wordChunks, applyCapitalization(randomWord, g.config.CapsMode))
	}
	
	// 3. Insert Separators
	for i := 0; i < g.config.SeparatorCount; i++ {
		sepChar, err := getRandomSeparatorChar(separatorPool)
		if err != nil {
			return "", fmt.Errorf("failed to get random char: %w", err)
		}
		
		numDigits := randomDigitCount(minDigits, maxDigits)
		digits, err := getRandomDigits(numDigits)
		if err != nil {
			return "", fmt.Errorf("failed to get random digits: %w", err)
		}
		
		separatorChunk := sepChar + digits + sepChar
		
		if len(wordChunks) < 2 {
			wordChunks = append(wordChunks, separatorChunk)
			continue
		}
		
		maxRange := big.NewInt(int64(len(wordChunks) - 1))
		insertionIndexBig, _ := rand.Int(randomSource, maxRange)
		insertionIndex := int(insertionIndexBig.Int64()) + 1
		
		wordChunks = append(wordChunks[:insertionIndex], append([]string{separatorChunk}, wordChunks[insertionIndex:]...)...)
	}

	if !typoApplied && g.config.TypoRate > 0 && g.config.Verbose {
		log.Println("Warning: Typo rate set but no typo was applied to any word.")
	}
	
	// 4. Final assembly
	return strings.Join(wordChunks, ""), nil
}

// --- ALL HELPER FUNCTIONS (now internal to the package) ---

// loadWordsFromFiles reads all specified files into a slice of strings
func loadWordsFromFiles(filenames []string) ([]string, error) {
	var fileWords []string
	for _, filename := range filenames {
		if filename == "" {
			continue
		}
		
		file, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("could not open file %s: %w", filename, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fileWords = append(fileWords, scanner.Text())
		}
		if err := scanner.Err(); err != nil && err != io.EOF {
			return nil, fmt.Errorf("error reading file %s: %w", filename, err)
		}
	}
	return fileWords, nil
}

// filterWords takes a slice of strings (from files/embed), filters, and de-duplicates them
func filterWords(wordsToParse []string) []string {
	uniqueWords := make(map[string]bool)
	
	for _, wordString := range wordsToParse {
		// This handles the embedded string, which might be multi-line
		scanner := bufio.NewScanner(strings.NewReader(wordString))
		for scanner.Scan() {
			word := scanner.Text()
			lowerWord := strings.ToLower(word)
			
			if !strings.ContainsAny(lowerWord, invalidChars) {
				uniqueWords[word] = true
			}
		}
	}
	
	wordsSlice := make([]string, 0, len(uniqueWords))
	for word := range uniqueWords {
		wordsSlice = append(wordsSlice, word)
	}
	return wordsSlice
}

// applyCapitalization handles all capitalization rules
func applyCapitalization(word string, mode string) string {
	if len(word) == 0 {
		return ""
	}
	
	switch mode {
	case "camel":
		// Uppercase the first letter and append the rest
		return strings.ToUpper(string(word[0])) + word[1:]
		
	case "random":
		runes := []rune(word)
		maxIndex := big.NewInt(int64(len(runes)))
		n, err := rand.Int(randomSource, maxIndex)
		if err != nil {
			return word // Fail silently, return original
		}
		idx := int(n.Int64())
		runes[idx] = unicode.ToUpper(runes[idx])
		return string(runes)

	case "none":
		return word // Return as-is (assuming list is mostly lowercase)
	
	default:
		return word
	}
}


func parseSeparatorPool(s string) []string {
	if len(s) > 0 {
		return strings.Split(s, "")
	}
	return defaultSpecialChars
}

func parseDigitRange(rangeStr string) (min int, max int, err error) {
	parts := strings.Split(rangeStr, "-")
	
	if len(parts) == 1 {
		val, e := strconv.Atoi(parts[0])
		if e != nil {
			return 0, 0, fmt.Errorf("invalid digit count '%s'", parts[0])
		}
		if val < 1 {
			return 0, 0, fmt.Errorf("digit count must be at least 1")
		}
		return val, val, nil
	} else if len(parts) == 2 {
		minVal, e1 := strconv.Atoi(parts[0])
		maxVal, e2 := strconv.Atoi(parts[1])
		
		if e1 != nil || e2 != nil {
			return 0, 0, fmt.Errorf("invalid digit range format")
		}
		if minVal < 1 || minVal > maxVal {
			return 0, 0, fmt.Errorf("invalid digit range (min must be >= 1 and <= max)")
		}
		return minVal, maxVal, nil
	}
	return 0, 0, fmt.Errorf("invalid digit range format. Use 'N' or 'N-M'")
}

func randomDigitCount(min int, max int) int {
	if min == max {
		return min
	}
	rangeSize := big.NewInt(int64(max - min + 1))
	n, err := rand.Int(randomSource, rangeSize)
	if err != nil {
		log.Fatalf("Error generating random digit count: %v", err) // Fatal error
	}
	return min + int(n.Int64())
}

func getRandomSeparatorChar(pool []string) (string, error) {
	max := big.NewInt(int64(len(pool)))
	n, err := rand.Int(randomSource, max)
	if err != nil {
		return "", err
	}
	return pool[n.Int64()], nil
}

func applyTypoTransposition(word string) string {
	if len(word) < 2 {
		return word
	}
	runes := []rune(word)
	maxIndex := big.NewInt(int64(len(runes) - 1))
	n, err := rand.Int(randomSource, maxIndex)
	if err != nil {
		// log.Printf("Warning: Failed to select typo index: %v", err)
		return word
	}
	i := int(n.Int64())
	runes[i], runes[i+1] = runes[i+1], runes[i]
	return string(runes)
}

func shouldApplyTypo(rate float64) bool {
	if rate <= 0 {
		return false
	}
	max := big.NewInt(10000)
	n, err := rand.Int(randomSource, max)
	if err != nil {
		// log.Printf("Warning: Failed to generate random number for typo: %v", err)
		return false
	}
	return n.Int64() < int64(rate*10000)
}

func getRandomDigits(n int) (string, error) {
	if n <= 0 {
		return "", nil
	}
	max := big.NewInt(int64(math.Pow10(n)))
	num, err := rand.Int(randomSource, max)
	if err != nil {
		return "", err
	}
	format := fmt.Sprintf("%%0%dd", n)
	return fmt.Sprintf(format, num), nil
}

func (g *Generator) getRandomWord() (string, error) {
	max := big.NewInt(int64(len(g.filteredWords)))
	n, err := rand.Int(randomSource, max)
	if err != nil {
		return "", err
	}
	return g.filteredWords[n.Int64()], nil
}
