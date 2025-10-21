# LUMOS Development Guide

**Version**: 0.1.0
**Date**: 2025-10-21

---

## Development Setup

### Prerequisites

- Go 1.21+ (install via Homebrew: `brew install go`)
- macOS or Linux (WSL supported)
- Git
- A decent text editor (VSCode recommended)

### Initial Setup

```bash
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMOS

# Download dependencies
go mod download

# Build binary
go build -o lumos ./cmd/lumos

# Test basic functionality
./lumos --help
```

### Project Structure Quick Reference

```
LUMOS/
├── cmd/lumos/main.go          ← CLI entry point
├── pkg/
│   ├── pdf/
│   │   ├── document.go        ← PDF document model
│   │   ├── search.go          ← Search functionality
│   │   └── cache.go           ← LRU cache
│   ├── ui/
│   │   ├── model.go           ← Main Bubble Tea model
│   │   └── keybindings.go     ← Vim keybindings
│   ├── config/
│   │   └── theme.go           ← Color themes
│   └── terminal/              ← (Phase 3: image support)
├── test/
│   ├── fixtures/              ← Test PDF files
│   └── benchmarks/            ← Performance tests
└── docs/
    ├── ARCHITECTURE.md        ← This is here
    └── DEVELOPMENT.md         ← This guide
```

---

## Development Workflow

### 1. Understanding the Codebase

**Start with**:
1. `README.md` - Project overview
2. `docs/ARCHITECTURE.md` - System design
3. `cmd/lumos/main.go` - Entry point
4. `pkg/pdf/document.go` - Core model
5. `pkg/ui/model.go` - UI logic

**Read in order**: Entry point → Domain logic → UI

### 2. Making Changes

**General Process**:

```bash
# Create a feature branch
git checkout -b feature/your-feature-name

# Make changes to source files
# Edit files in pkg/pdf, pkg/ui, etc.

# Build to check for compilation errors
go build -o lumos ./cmd/lumos

# Test your changes
./lumos ~/Documents/test-pdf.pdf

# Run tests (when available)
go test ./...

# Commit changes
git add .
git commit -m "feat: clear description of what changed"

# Push to origin
git push origin feature/your-feature-name
```

### 3. Debugging

**Print Debugging**:
```go
import "log"

// In your code
log.Printf("Debug info: %v", variable)
```

**Structured Logging** (recommended):
```go
import "log/slog"

slog.Debug("msg", "key", value)
slog.Info("msg", "key", value)
slog.Error("msg", "error", err)
```

**Interactive Debugging**:
```bash
# Using Delve debugger
dlv debug ./cmd/lumos -- ~/Documents/test.pdf

# Then use 'help' to see commands
(dlv) help
```

### 4. Testing

**Run tests**:
```bash
go test ./...                  # All tests
go test ./pkg/pdf/...         # Specific package
go test -v ./...              # Verbose output
go test -race ./...           # Detect race conditions
```

**Write a test**:
```go
// In pkg/pdf/document_test.go
func TestPageExtraction(t *testing.T) {
    doc, err := NewDocument("test.pdf", 5)
    if err != nil {
        t.Fatalf("Failed to open PDF: %v", err)
    }

    page, err := doc.GetPage(1)
    if err != nil {
        t.Fatalf("Failed to get page: %v", err)
    }

    if page == nil {
        t.Error("Expected page, got nil")
    }
}
```

### 5. Benchmarking

**Run benchmarks**:
```bash
go test -bench=. ./pkg/pdf/
go test -bench=. -benchmem ./pkg/pdf/
```

**Write a benchmark**:
```go
// In pkg/pdf/document_test.go
func BenchmarkPageLoad(b *testing.B) {
    doc, _ := NewDocument("test.pdf", 5)
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        doc.GetPage(1)
    }
}
```

---

## Code Style Guide

### Go Conventions

1. **Naming**:
   - Exported: `CamelCase` (e.g., `Document`, `GetPage`)
   - Unexported: `camelCase` (e.g., `document`, `getPage`)
   - Constants: `UPPER_CASE` or `CamelCase`

2. **Package Layout**:
   - Group related functions in same file
   - Keep files <500 lines
   - One main type per file when possible

3. **Comments**:
   - Comment all exported types and functions
   - Keep comments concise and clear
   - Explain *why*, not just *what*

Example:
```go
// Document represents a loaded PDF file with page caching
type Document struct {
    filepath string
    reader   *pdf.Reader
    pages    int
}

// GetPage retrieves text content from a specific page.
// Results are cached in an LRU cache of size maxCache.
func (d *Document) GetPage(pageNum int) (*PageInfo, error) {
    // ...
}
```

### Error Handling

```go
// Good
if err != nil {
    return fmt.Errorf("failed to load page: %w", err)
}

// Bad
if err != nil {
    panic(err)  // Crashes app
}

// Bad
if err != nil {
    return err  // Loses context
}
```

### Concurrency

Use sync primitives correctly:
```go
type Cache struct {
    data map[int]string
    mu   sync.RWMutex  // Protects data
}

// Reads - use RLock
func (c *Cache) Get(key int) string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.data[key]
}

// Writes - use Lock
func (c *Cache) Set(key int, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}
```

---

## Common Tasks

### Adding a New Command

1. Add handler to `KeyHandler` in `pkg/ui/keybindings.go`
2. Add message type (e.g., `MyCommandMsg`)
3. Handle in `Model.Update()`
4. Implement logic and rendering

Example (adding `:w` to write):
```go
// 1. In keybindings.go
case ":w":
    return WriteCmd

// 2. Define message
type WriteMsg struct{}

// 3. In model.go Update()
case WriteMsg:
    return m.handleWrite()

// 4. Implement
func (m *Model) handleWrite() tea.Cmd {
    // Logic here
    return nil
}
```

### Adding a New Theme

1. Add theme to `pkg/config/theme.go`:
```go
var MyTheme = Theme{
    Name:       "mytheme",
    Background: "#...",
    Text:       "#...",
    // ... other colors
}
```

2. Update `GetTheme()` function:
```go
func GetTheme(name string) Theme {
    switch name {
    case "mytheme":
        return MyTheme
    // ...
    }
}
```

3. Users can now use: `lumos --theme mytheme file.pdf` (Phase 2+)

### Adding a New Keybinding

1. Add case to `handleNormalKey()` in `keybindings.go`:
```go
case "your-key":
    return YourCommand
```

2. Define the command function:
```go
var YourCommand = func() tea.Msg {
    return YourMsg{}
}
```

3. Handle in `Model.Update()`:
```go
case YourMsg:
    // Logic
```

4. Add to `VimKeybindingReference` for help:
```go
"your-key": "Description of what it does",
```

---

## Performance Optimization

### Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof ./pkg/pdf/
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof ./pkg/pdf/
go tool pprof mem.prof
```

### Common Optimizations

1. **Cache Misses**: Monitor `cache.HitRate()` - aim for >80%
2. **Page Loading**: Profile `GetPage()` - should be <50ms
3. **Search**: Consider lazy evaluation instead of processing all pages
4. **Memory**: Use `strings.Builder` for string concatenation

### Benchmarking Results Target

```
BenchmarkPageLoad:       <60ms per page
BenchmarkSearch:         <100ms for typical doc
BenchmarkCacheLookup:    <1ms
BenchmarkMemory:         <50MB for 100-page doc
```

---

## Dependency Management

### Adding a Dependency

```bash
# Add new dependency
go get github.com/user/package@latest

# Update go.mod and go.sum
go mod tidy

# Verify no issues
go mod verify
```

### Updating Dependencies

```bash
# Update all dependencies to latest minor versions
go get -u ./...

# Check for security vulnerabilities
go list -json -m all | nancy sleuth
```

### Current Dependencies

See `go.mod` for current versions. Key dependencies:
- `bubbletea` - TUI framework
- `bubbles` - UI components
- `lipgloss` - Terminal styling
- `ledongthuc/pdf` - PDF parsing
- `pdfcpu` - PDF manipulation (Phase 2+)

---

## Debugging Common Issues

### Issue: Build Fails

```bash
# Check for syntax errors
go build ./...

# Check imports
go mod tidy

# Verify dependencies are downloaded
go mod download
```

### Issue: PDF Not Loading

1. Check file exists and is readable
2. Check if PDF is corrupted
3. Check if PDF is encrypted (not supported yet)
4. Check logs for specific error

### Issue: Slow Page Navigation

1. Profile with pprof to find bottleneck
2. Check cache hit rate with `cache.HitRate()`
3. Check if memory is limited
4. Profile `ledongthuc/pdf` extraction

### Issue: Memory Leak

1. Run with race detector: `go test -race`
2. Check for goroutine leaks
3. Verify all mutexes unlocked properly
4. Profile memory over time

---

## Release Process

### Preparing a Release

```bash
# Update version in cmd/lumos/main.go
# Update CHANGELOG.md
# Create release notes

# Tag the release
git tag v0.1.0
git push origin v0.1.0

# GoReleaser will build cross-platform binaries
# (configured in .goreleaser.yml)
goreleaser release
```

---

## Useful Resources

### Go Resources
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Testing](https://golang.org/pkg/testing/)

### Library Documentation
- [ledongthuc/pdf](https://github.com/ledongthuc/pdf)
- [Bubble Tea Guide](https://github.com/charmbracelet/bubbletea)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)

### Similar Projects
- [LUMINA CCN](../../../LUMINA/ccn) - Markdown viewer reference
- [Glow](https://github.com/charmbracelet/glow) - Markdown renderer
- [Pdfless](https://github.com/Aaaaatyle/pdfless) - Rust PDF reader

---

## Getting Help

### Documentation
- `README.md` - Project overview
- `docs/ARCHITECTURE.md` - System design
- `docs/KEYBINDINGS.md` - Keyboard shortcuts (planned)
- `docs/PERFORMANCE.md` - Performance notes (planned)

### Code Comments
All modules have inline documentation. Read the source!

### External Help
- Go forums: https://forum.golangbridge.org/
- Stack Overflow: tag `go` and `pdf`
- Charm ecosystem: https://charm.sh

---

**Last Updated**: 2025-10-21
