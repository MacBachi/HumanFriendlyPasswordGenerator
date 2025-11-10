// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/atotto/clipboard"
	// Import your new local generator package
	"github.com/MacBachi/HumanFriendlyPasswordGenerator/generator" 
)

func main() {
	// --- Flag parsing ---
	wordCount := flag.Int("w", 4, "Number of words")
	separatorCount := flag.Int("s", 1, "Number of separator blocks")
	separatorPoolStr := flag.String("sep", "!", "Allowed separator characters (e.g., '+-_!')")
	digitRangeStr := flag.String("d-range", "4", "Digit count per block (e.g., '4' or '3-6')")
	typoRate := flag.Float64("typo-rate", 0.33, "Probability (0.0-1.0) of a typo per word")
	altList := flag.String("altlist", "", "Alternative wordlist file (ignores default)")
	mergeList := flag.String("mergelist", "", "Additional wordlist file (merges with default)")
	
	// NEW: QoL Flags
	numToGenerate := flag.Int("n", 1, "Number of passphrases to generate")
	copyToClipboard := flag.Bool("c", false, "Copy the first generated passphrase to the clipboard")
	
	// NEW: Entropie Flag
	capsMode := flag.String("caps", "camel", "Capitalization mode: camel, random, none")
	
	// NEW: Verbose Flag
	verbose := flag.Bool("v", false, "Verbose output")
	
	flag.Parse()

	// --- Input Validation ---
	if *numToGenerate < 1 {
		*numToGenerate = 1
	}
	*capsMode = strings.ToLower(*capsMode)
	if *capsMode != "camel" && *capsMode != "random" && *capsMode != "none" {
		log.Fatalf("Error: Invalid --caps mode. Use 'camel', 'random', or 'none'.")
	}

	// --- Create Config ---
	config := generator.Config{
		WordCount:      *wordCount,
		SeparatorCount: *separatorCount,
		SeparatorPool:  *separatorPoolStr,
		DigitRange:     *digitRangeStr,
		TypoRate:       *typoRate,
		AltList:        *altList,
		MergeList:      *mergeList,
		CapsMode:       *capsMode,
		Verbose:        *verbose,
	}

	// --- Initialize Generator ---
	// NewGenerator now handles loading and filtering the wordlists (incl. embed)
	gen, err := generator.NewGenerator(config)
	if err != nil {
		log.Fatalf("Error initializing generator: %v", err)
	}

	if config.Verbose {
		fmt.Println("Starting HumanFriendlyPasswordGenerator...")
		fmt.Printf("Loaded %d valid words.\n", gen.WordCount())
		fmt.Println("-------------------------------------------------")
	}

	// --- Generation Loop ---
	var firstPassword string
	for i := 0; i < *numToGenerate; i++ {
		passphrase, err := gen.Generate()
		if err != nil {
			log.Fatalf("Error generating passphrase: %v", err)
		}
		
		fmt.Println(passphrase)
		
		if i == 0 {
			firstPassword = passphrase
		}
	}

	if config.Verbose {
		fmt.Println("-------------------------------------------------")
	}

	// --- QoL: Copy to Clipboard ---
	if *copyToClipboard {
		if err := clipboard.WriteAll(firstPassword); err != nil {
			log.Printf("Warning: Failed to copy to clipboard: %v", err)
		} else {
			fmt.Println("Copied first passphrase to clipboard!")
		}
	}
}
