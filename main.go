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

func main() {
	// 1. Parse Arguments
	ext := flag.String("ext", "", "File extension to count (e.g., .go, .txt)")
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
	if *skipComments {
		fmt.Println("Excluding comments...")
	}

	// 2. Walk Directory
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Skip .git and other common hidden directories if needed,
			// but for now we just walk everything.
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		// 3. Filter Files
		if strings.HasSuffix(path, *ext) {
			lines, err := countLines(path, *skipComments)
			if err != nil {
				fmt.Printf("Error reading %s: %v\n", path, err)
				return nil // Continue even if one file fails
			}
			totalLines += lines
			fileCount++
			// Optional: print per-file count
			// fmt.Printf("%s: %d\n", path, lines)
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
func countLines(path string, skipComments bool) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines int64
	inBlockComment := false

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if !skipComments {
			lines++
			continue
		}

		if text == "" {
			continue
		}

		// Handle block comments (C-style /* ... */)
		if inBlockComment {
			if idx := strings.Index(text, "*/"); idx != -1 {
				inBlockComment = false
				text = text[idx+2:]
			} else {
				continue // Still in block comment
			}
		}

		// Clean up the line iteratively
		for {
			startIdx := strings.Index(text, "/*")
			lineCommentIdx := strings.Index(text, "//")

			// Check for single line comment "//"
			// It takes precedence if it appears before "/*"
			if lineCommentIdx != -1 && (startIdx == -1 || lineCommentIdx < startIdx) {
				text = text[:lineCommentIdx]
				break
			}

			// Check for block comment start "/*"
			if startIdx != -1 {
				// Look for end of block comment "*/"
				endIdx := strings.Index(text[startIdx+2:], "*/")
				if endIdx != -1 {
					// Complete block comment on same line: remove it and continue
					text = text[:startIdx] + text[startIdx+2+endIdx+2:]
					continue
				} else {
					// Block comment continues to next line
					inBlockComment = true
					text = text[:startIdx]
					break
				}
			}

			break // No more comments found
		}

		if strings.TrimSpace(text) != "" {
			lines++
		}
	}
	return lines, scanner.Err()
}