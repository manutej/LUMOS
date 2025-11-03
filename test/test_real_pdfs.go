package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/luxor/lumos/pkg/pdf"
	"github.com/luxor/lumos/pkg/ui"
)

// Test loading various PDFs to ensure they work with LUMOS
func main() {
	testPDFs := []string{
		"../../docs/HEKAT_OPERATIONAL_MODEL.pdf",
		"../../docs/HEKAT_PRACTICAL_PATTERNS.pdf",
		"../../docs/HEKAT_INTEGRATION_COMPLETE.pdf",
		"../../PROJECTS/paper2agent/papers/On Meta-Prompting.pdf",
		"../../PROJECTS/paper2agent/papers/consciousness-as-functor.pdf",
		"../../PROJECTS/paper2agent/cat-doc-qa/document-category-theory-paper.pdf",
	}

	fmt.Println("LUMOS PDF Loading Test")
	fmt.Println("======================")
	fmt.Println()

	passed := 0
	failed := 0

	for i, pdfPath := range testPDFs {
		fmt.Printf("[%d/%d] Testing: %s\n", i+1, len(testPDFs), filepath.Base(pdfPath))

		// Check if file exists
		if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
			fmt.Printf("  ⚠️  SKIP - File not found\n\n")
			continue
		}

		// Get file size
		fileInfo, _ := os.Stat(pdfPath)
		sizeKB := fileInfo.Size() / 1024
		fmt.Printf("  📄 Size: %d KB\n", sizeKB)

		// Try to load document
		doc, err := pdf.NewDocument(pdfPath, 5)
		if err != nil {
			fmt.Printf("  ❌ FAIL - Error loading PDF: %v\n\n", err)
			failed++
			continue
		}

		// Get basic info
		pageCount := doc.GetPageCount()
		metadata := doc.GetMetadata()
		fmt.Printf("  📖 Pages: %d\n", pageCount)
		if metadata.Title != "" {
			fmt.Printf("  📝 Title: %s\n", metadata.Title)
		}
		if metadata.Author != "" {
			fmt.Printf("  ✍️  Author: %s\n", metadata.Author)
		}

		// Try to create UI model
		model := ui.NewModel(doc)
		if model == nil {
			fmt.Printf("  ❌ FAIL - Could not create UI model\n\n")
			failed++
			continue
		}

		// Try to load first page
		page, err := doc.GetPage(1)
		if err != nil {
			fmt.Printf("  ❌ FAIL - Could not load page 1: %v\n\n", err)
			failed++
			continue
		}

		textLen := len(page.Text)
		fmt.Printf("  📄 Page 1 text length: %d characters\n", textLen)

		// Try to render (without actually displaying)
		view := model.View()
		if view == "" {
			fmt.Printf("  ❌ FAIL - View rendering returned empty string\n\n")
			failed++
			continue
		}

		fmt.Printf("  ✅ PASS - Document loaded and rendered successfully\n\n")
		passed++
	}

	fmt.Println("======================")
	fmt.Printf("Results: %d passed, %d failed\n", passed, failed)

	if failed > 0 {
		os.Exit(1)
	}
}
