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
			lines, err := countLines(path)
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
func countLines(path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines int64
	for scanner.Scan() {
		lines++
	}
	return lines, scanner.Err()
}