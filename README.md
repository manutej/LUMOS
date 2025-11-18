# LUMOS

A dark mode PDF reader for developers, built with Go and Bubble Tea.

**Status**: Phase 1 MVP - 50% Complete | [Specifications](.specify/SPECIFICATION_INDEX.md) | [Handoff Guide](.specify/HANDOFF.md)

---

## Features

- 🌙 **Dark mode by default** - Optimized for long reading sessions
- ⌨️ **Vim-style keybindings** - Navigate naturally with j/k/h/l
- 🎯 **Distraction-free** - Minimal UI, maximum focus
- ⚡ **Lightning fast** - Instant startup (<70ms), smooth scrolling
- 🔍 **Full-text search** - Find content quickly with highlighting (<50μs)
- 📱 **Terminal-native** - Works in any terminal emulator
- 🎨 **Themeable** - Dark and light modes
- ✅ **Production tested** - 94.4% test coverage

---

## Quick Start

```bash
# Clone the repository
git clone https://github.com/manutej/LUMOS.git
cd LUMOS

# Build the binary
make build

# Open a PDF
./build/lumos ~/Documents/paper.pdf
```

---

## Current Status

### ✅ Completed (50%)
- **Core PDF Engine** - Text extraction, caching, search (94.4% tested)
- **Build System** - Clean compilation, dependency management
- **Test Infrastructure** - 42 tests, 9 benchmarks, fixtures
- **Performance** - Exceeding all targets (<70ms startup, <20ms cache)

### 🚧 In Progress
- **TUI Framework** - Bubble Tea integration (Milestone 1.4)
- **3-Pane Layout** - Metadata | Viewer | Search

### ⏳ Upcoming
- **Vim Keybindings** - Full navigation suite (Milestone 1.5)
- **Dark Mode Polish** - Theme refinement (Milestone 1.6)

See [Roadmap](.specify/specs/phase-1-mvp.md) for detailed timeline.

---

## Architecture

```
┌─────────────────┐
│  CLI Entry      │  cmd/lumos/main.go
│  (tea.Program)  │
└────────┬────────┘
         │
┌─────────────────┐
│  PDF Package    │  pkg/pdf/ (✅ Complete)
│  - Document     │  - 94.4% test coverage
│  - Search       │  - <50μs search performance
│  - LRU Cache    │  - <100ns cache operations
└────────┬────────┘
         │
┌─────────────────┐
│  UI Package     │  pkg/ui/ (🚧 In Progress)
│  - MVU Model    │  - Bubble Tea framework
│  - Keybindings  │  - Vim navigation
│  - 3-Pane Layout│  - Responsive design
└─────────────────┘
```

---

## Keybindings (Coming in v0.1.0)

| Key | Action | Status |
|-----|--------|--------|
| `j/k` | Scroll down/up | ⏳ |
| `d/u` | Half page down/up | ⏳ |
| `gg/G` | First/last page | ⏳ |
| `Ctrl+N/P` | Next/previous page | ⏳ |
| `/` | Search | ⏳ |
| `n/N` | Next/previous match | ⏳ |
| `q` | Quit | 🚧 |
| `?` | Help | ⏳ |

---

## Development

### Building

```bash
# Build binary
make build

# Run tests (42 passing)
make test

# Check coverage (94.4%)
make coverage

# Run benchmarks
make bench

# Full CI checks
make ci-check
```

### Project Structure

```
LUMOS/
├── .specify/              # Specification-driven development
│   ├── constitution.md    # Architectural principles
│   ├── HANDOFF.md        # Quick start for developers
│   └── specs/            # Detailed specifications
├── cmd/lumos/            # CLI entry point
├── pkg/                  # Core packages
│   ├── pdf/             # PDF operations (complete)
│   ├── ui/              # TUI components (in progress)
│   └── config/          # Configuration
├── test/                 # Test fixtures and guides
└── docs/                 # Additional documentation
```

### Contributing

1. Read the [Handoff Guide](.specify/HANDOFF.md)
2. Check [Priorities](.specify/PRIORITIES.md)
3. Follow [Constitution](.specify/constitution.md)
4. Write tests first (TDD)
5. Maintain >80% coverage

---

## Performance

| Metric | Target | Current | Status |
|--------|--------|---------|--------|
| Cold start | <100ms | ~70ms | ✅ |
| Page switch (cached) | <50ms | <20ms | ✅ |
| Search (100 pages) | <100ms | <5ms | ✅ |
| Memory (10MB PDF) | <50MB | ~15MB | ✅ |
| Test coverage | >80% | 94.4% | ✅ |

---

## Roadmap

### Phase 1: MVP (Nov 2025) - Current
- ✅ Core PDF engine
- ✅ Test infrastructure
- 🚧 TUI implementation
- ⏳ Vim keybindings
- ⏳ Dark mode polish

### Phase 2: Enhanced (Q1 2026)
- Table of contents
- Bookmarks
- Annotations
- Configuration file

### Phase 3: Images (Q2 2026)
- Image rendering
- Table detection
- Complex layouts

### Phase 4: AI (Q3 2026)
- Claude Agent SDK
- Summarization
- Q&A features

---

## Related Projects

- [LUMINA](https://github.com/manutej/LUMINA) - Markdown viewer with similar TUI

---

## License

MIT License - See [LICENSE](LICENSE) file for details.

---

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [ledongthuc/pdf](https://github.com/ledongthuc/pdf) - PDF parsing
- [Charm](https://charm.sh) - Terminal UI tools

---

**Repository**: https://github.com/manutej/LUMOS
**Documentation**: [Specifications](.specify/SPECIFICATION_INDEX.md)
**Quick Start**: [Handoff Guide](.specify/HANDOFF.md)