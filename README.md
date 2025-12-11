# LOC Counter

A fast, cross-platform CLI tool written in Go to count lines of code (LOC) in your projects. It supports recursive directory traversal and accurate comment filtering for a wide variety of programming languages.

## Features

-   **Recursive Counting**: traverse directories to count lines for specific file extensions.
-   **Comment Filtering**: Optional `-skip-comments` flag to exclude comments from the count.
-   **Multi-Language Support**: Smart comment parsing for Go, Python, Java, C++, SQL, HTML, Shell, and many more.
-   **Cross-Platform**: Compile for macOS, Linux, and Windows.

## Installation

### From Source

Ensure you have Go installed (1.21+ recommended).

```bash
git clone https://github.com/onlyargon/loc-counter.git
cd loc-counter
go install
```

### Pre-built Binaries

You can build the binaries yourself using the included `Makefile`:

```bash
make build-all
```

This will generate binaries in the `bin/` directory for macOS, Linux, and Windows.

## Usage

The basic syntax is:

```bash
loc-counter -ext <extension> [options] [directory]
```

### Examples

**Count Go lines in the current directory:**

```bash
loc-counter -ext .go
```

**Count Python lines in a specific project:**

```bash
loc-counter -ext .py /path/to/my/project
```

**Count lines excluding comments:**

```bash
loc-counter -ext .js -skip-comments .
```

## Supported Languages (Comment Filtering)

When using `-skip-comments`, `loc-counter` automatically detects the comment syntax based on the file extension. Partial list of supported languages:

-   **C-Style** (`//`, `/* ... */`): C, C++, Java, JavaScript, TypeScript, Go, Rust, Swift, Kotlin, PHP
-   **Script-Style** (`#`): Python, Ruby, Perl, Bash/Shell, YAML, PowerShell
-   **Database/Config** (`--`, `<!--`): SQL, Lua, Haskell, HTML, XML
-   **Others**: VB (`'`), ASM/Lisp (`;`), Matlab/LaTeX (`%`)

## License

MIT
