package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luxor/lumos/pkg/pdf"
	"github.com/luxor/lumos/pkg/ui"
)

const Version = "0.1.0"

func main() {
	// Parse command line flags
	version := flag.Bool("version", false, "Show version")
	help := flag.Bool("help", false, "Show help")
	keys := flag.Bool("keys", false, "Show keyboard shortcuts")
	flag.BoolVar(help, "h", false, "Show help (short)")
	flag.BoolVar(version, "v", false, "Show version (short)")
	flag.BoolVar(keys, "k", false, "Show keyboard shortcuts (short)")

	flag.Parse()

	// Handle flags
	if *version {
		fmt.Printf("LUMOS v%s - PDF Dark Mode Reader\n", Version)
		os.Exit(0)
	}

	if *help {
		printHelp()
		os.Exit(0)
	}

	if *keys {
		printKeys()
		os.Exit(0)
	}

	// Get PDF file path from arguments
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: lumos [flags] <pdf-file>\n")
		fmt.Fprintf(os.Stderr, "Try 'lumos --help' for more information.\n")
		os.Exit(1)
	}

	pdfPath := args[0]

	// Expand home directory if needed
	if pdfPath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not expand home directory: %v\n", err)
			os.Exit(1)
		}
		pdfPath = filepath.Join(homeDir, pdfPath[1:])
	}

	// Check if file exists
	if _, err := os.Stat(pdfPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: File not found: %s\n", pdfPath)
		os.Exit(1)
	}

	// Load PDF
	doc, err := pdf.NewDocument(pdfPath, 5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading PDF: %v\n", err)
		os.Exit(1)
	}

	// Create and run TUI
	model := ui.NewModel(doc)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(`LUMOS - PDF Dark Mode Reader
A developer-friendly PDF reader with dark mode and vim keybindings

USAGE:
  lumos [flags] <pdf-file>

OPTIONS:
  -h, --help      Show this help message
  -v, --version   Show version information
  -k, --keys      Show keyboard shortcuts reference

EXAMPLES:
  # Open a PDF file
  lumos ~/Documents/paper.pdf

  # Show help
  lumos --help

  # Show keyboard shortcuts
  lumos --keys

PROJECT:
  LUMOS is a companion to LUMINA (markdown viewer)
  https://github.com/luxor/lumos

KEYBOARD SHORTCUTS:
  Navigation:     j/k (scroll), d/u (half page), gg/G (top/bottom)
  Page Nav:       Ctrl+N (next), Ctrl+P (previous)
  Search:         / (search), n/N (next/prev match)
  UI:             Tab (cycle panes), 1/2 (dark/light mode)
  General:        ? (help), q (quit)

For more information, see documentation at:
  /Users/manu/Documents/LUXOR/PROJECTS/LUMOS/README.md
`)
}

func printKeys() {
	fmt.Print(`LUMOS - Keyboard Shortcuts

NORMAL MODE
===========

Navigation:
  j or ↓          Scroll down one line
  k or ↑          Scroll up one line
  d               Scroll down half page
  u               Scroll up half page
  gg              Go to top of document
  G               Go to bottom of document
  Ctrl+N          Go to next page
  Ctrl+P          Go to previous page

Search:
  /               Start search
  n               Go to next match
  N               Go to previous match
  Esc             Exit search mode

UI Control:
  Tab             Cycle through panes
  1               Switch to dark mode
  2               Switch to light mode
  ?               Toggle this help screen
  :               Enter command mode

General:
  q or Ctrl+C     Quit the application

SEARCH MODE
===========
  Type to search
  Enter           Execute search
  Esc             Exit search mode
  n/N             Next/previous match

COMMAND MODE
============
  Type command name
  Enter           Execute command
  Esc             Exit command mode

TIPS
====
  • Use vim keybindings for fast navigation
  • Dark mode by default for long reading sessions
  • LRU caching keeps 5 pages in memory
  • Search is case-insensitive by default
  • Mark pages with vim marks (Phase 2+)
`)
}
