package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// BlockComment defines the start and end tokens for a block comment
type BlockComment struct {
	Start string
	End   string
}

// CommentStyle defines the comment syntax for a language
type CommentStyle struct {
	LineComments  []string
	BlockComments []BlockComment
}

// LanguageStyles maps file extensions to their comment styles
var languageStyles = map[string]CommentStyle{
	// Single-line: //, Block: /* */
	".c":    cStyle,
	".cpp":  cStyle,
	".java": cStyle,
	".js":   cStyle,
	".ts":   cStyle,
	".cs":   cStyle,
	".swift": cStyle,
	".rs":   cStyle, // Rust
	".go":   cStyle,
	".kt":   cStyle, // Kotlin
	".php":  cStyle,
	".css":  {BlockComments: []BlockComment{{"/*", "*/"}}}, // CSS usually only has block comments

	// Single-line: #
	".py": {
		LineComments: []string{"#"},
		BlockComments: []BlockComment{
			{`"""`, `"""`},
			{`'''`, `'''`},
		},
	},
	".rb": {
		LineComments: []string{"#"},
		BlockComments: []BlockComment{
			{"=begin", "=end"},
		},
	},
	".pl":  shellStyle, // Perl
	".sh":  shellStyle, // Bash
	".yaml": shellStyle,
	".yml":  shellStyle,
	".r":   shellStyle,
	".ps1": shellStyle, // PowerShell

	// Single-line: --
	".sql":  {LineComments: []string{"--"}, BlockComments: []BlockComment{{"/*", "*/"}}},
	".lua":  {LineComments: []string{"--"}, BlockComments: []BlockComment{{"--[[", "]]"}}},
	".hs":   {LineComments: []string{"--"}, BlockComments: []BlockComment{{"{-", "-}"}}},
	".ada":  {LineComments: []string{"--"}},

	// Single-line: '
	".vb": {LineComments: []string{"'"}},

	// Single-line: ;
	".asm": {LineComments: []string{";"}},
	".lisp": {LineComments: []string{";"}},
	".clj":  {LineComments: []string{";"}},
	".ini":  {LineComments: []string{";", "#"}}, // INI often supports both

	// Single-line: %
	".m":   {LineComments: []string{"%"}, BlockComments: []BlockComment{{"%{", "%}"}}}, // Matlab/Octave
	".tex": {LineComments: []string{"%"}},

	// Mixed / Others
	".bat": {LineComments: []string{"REM", "::"}},
	".html": {BlockComments: []BlockComment{{"<!--", "-->"}}},
	".xml":  {BlockComments: []BlockComment{{"<!--", "-->"}}},
	".pas":  {BlockComments: []BlockComment{{"(*", "*)"}, {"{", "}"}}},
}

var (
	cStyle     = CommentStyle{LineComments: []string{"//"}, BlockComments: []BlockComment{{"/*", "*/"}}}
	shellStyle = CommentStyle{LineComments: []string{"#"}}
	luaStyle   = CommentStyle{LineComments: []string{"--"}}
)

func main() {
	// 1. Parse Arguments
	ext := flag.String("ext", "", "File extension to count (e.g., .go, .py)")
	skipComments := flag.Bool("skip-comments", false, "Exclude comments from count")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <directory>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *ext == "" {
		fmt.Println("Please provide a file extension using -ext")
		flag.Usage()
		os.Exit(1)
	}

	// Ensure extension has a dot
	if !strings.HasPrefix(*ext, ".") {
		*ext = "." + *ext
	}

	root := "."
	if len(flag.Args()) > 0 {
		root = flag.Args()[0]
	}

	var totalLines int64
	var fileCount int

	fmt.Printf("Counting lines for files with extension: %s in %s\n", *ext, root)
	
	style, hasStyle := languageStyles[*ext]
	if *skipComments {
		if hasStyle {
			fmt.Println("Excluding comments...")
		} else {
			fmt.Printf("Warning: No comment style defined for %s. Treating all lines as code.\n", *ext)
			// fallback: disable comment skipping if we don't know the style
			*skipComments = false
		}
	}

	// 2. Walk Directory
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		// 3. Filter Files
		if strings.HasSuffix(path, *ext) {
			lines, err := countLines(path, *skipComments, style)
			if err != nil {
				fmt.Printf("Error reading %s: %v\n", path, err)
				return nil
			}
			totalLines += lines
			fileCount++
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	// 5. Output
	fmt.Printf("Total files: %d\n", fileCount)
	fmt.Printf("Total lines: %d\n", totalLines)
}

// 4. Count Lines
func countLines(path string, skipComments bool, style CommentStyle) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines int64
	
	// State for block comments
	inBlockComment := false
	var currentBlockEnd string

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if !skipComments {
			lines++
			continue
		}

		if text == "" {
			continue
		}

		// If inside a block comment, look for the end token
		if inBlockComment {
			if idx := strings.Index(text, currentBlockEnd); idx != -1 {
				inBlockComment = false
				text = text[idx+len(currentBlockEnd):]
				// clear end token
				currentBlockEnd = ""
			} else {
				continue // Still in block comment
			}
		}

		// Process the line for comments
		for {
			text = strings.TrimSpace(text)
			if text == "" {
				break
			}

			// Find nearest comment start
			bestIdx := -1
			matchType := 0 // 1: line, 2: block
			var matchedBlock BlockComment

			// check line comments
			for _, start := range style.LineComments {
				idx := strings.Index(text, start)
				if idx != -1 {
					if bestIdx == -1 || idx < bestIdx {
						bestIdx = idx
						matchType = 1
					}
				}
			}

			// check block comments
			for _, bc := range style.BlockComments {
				idx := strings.Index(text, bc.Start)
				if idx != -1 {
					if bestIdx == -1 || idx < bestIdx {
						bestIdx = idx
						matchType = 2
						matchedBlock = bc
					}
				}
			}

			if bestIdx == -1 {
				// No comments found in remaining text
				break
			}

			if matchType == 1 {
				// Line comment found, ignore rest of line
				text = text[:bestIdx]
				break
			} else {
				// Block comment found
				// Check if it closes on the same line
				remaining := text[bestIdx+len(matchedBlock.Start):]
				endIdx := strings.Index(remaining, matchedBlock.End)
				
				if endIdx != -1 {
					// Closed on same line
					// remove the comment part
					text = text[:bestIdx] + remaining[endIdx+len(matchedBlock.End):]
					continue // loop again to find more comments
				} else {
					// Opens here, continues
					inBlockComment = true
					currentBlockEnd = matchedBlock.End
					text = text[:bestIdx]
					break
				}
			}
		}

		if strings.TrimSpace(text) != "" {
			lines++
		}
	}
	return lines, scanner.Err()
}