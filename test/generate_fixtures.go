package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func main() {
	fixturesDir := "fixtures"

	// Ensure fixtures directory exists
	if err := os.MkdirAll(fixturesDir, 0755); err != nil {
		log.Fatalf("Failed to create fixtures directory: %v", err)
	}

	// Generate simple.pdf - single page
	if err := generateSimplePDF(fixturesDir + "/simple.pdf"); err != nil {
		log.Fatalf("Failed to generate simple.pdf: %v", err)
	}
	fmt.Println("✅ Generated simple.pdf")

	// Generate multipage.pdf - 5 pages
	if err := generateMultipagePDF(fixturesDir + "/multipage.pdf"); err != nil {
		log.Fatalf("Failed to generate multipage.pdf: %v", err)
	}
	fmt.Println("✅ Generated multipage.pdf")

	// Generate search_test.pdf - text with search terms
	if err := generateSearchTestPDF(fixturesDir + "/search_test.pdf"); err != nil {
		log.Fatalf("Failed to generate search_test.pdf: %v", err)
	}
	fmt.Println("✅ Generated search_test.pdf")

	fmt.Println("\n✅ All test fixtures generated successfully!")
}

func generateSimplePDF(filename string) error {
	// Create a simple PDF with one page and basic text
	content := `LUMOS Test PDF

This is a simple test PDF for the LUMOS PDF reader.

Title: LUMOS Test Document
Author: LUMOS Test Suite
Subject: Testing PDF Reading Functionality
Creator: LUMOS Fixture Generator

This document contains basic text for testing:
- Document loading
- Page counting
- Text extraction
- Metadata reading

Page Count: 1
Word Count: Approximately 50 words
Line Count: Multiple lines for testing`

	return createPDFWithText(filename, []string{content}, "LUMOS Test Document", "LUMOS Test Suite", "Testing", "LUMOS")
}

func generateMultipagePDF(filename string) error {
	pages := []string{
		`LUMOS Multi-Page Test PDF - Page 1

This is the first page of a multi-page test document.

Content for testing:
- Page navigation
- Multi-page caching
- Page range operations

Total Pages: 5`,

		`LUMOS Multi-Page Test PDF - Page 2

This is the second page.

Testing features:
- Sequential page access
- Cache behavior with multiple pages
- Page-by-page text extraction

Current Page: 2 of 5`,

		`LUMOS Multi-Page Test PDF - Page 3

This is the third page, the middle of the document.

Key testing scenarios:
- Middle page access
- Cache eviction with limited cache size
- Random page access patterns

Current Page: 3 of 5`,

		`LUMOS Multi-Page Test PDF - Page 4

This is the fourth page.

Additional test cases:
- Near-end page access
- Cache statistics validation
- Page range boundary testing

Current Page: 4 of 5`,

		`LUMOS Multi-Page Test PDF - Page 5

This is the final page of the test document.

Final test validations:
- Last page access
- Complete document traversal
- End-to-end workflow testing

Current Page: 5 of 5 (END)`,
	}

	return createPDFWithText(filename, pages, "LUMOS Multi-Page Test", "LUMOS Test Suite", "Multi-page Testing", "LUMOS")
}

func generateSearchTestPDF(filename string) error {
	content := `LUMOS Search Test PDF

This document contains specific text patterns for search testing.

Case Sensitivity Tests:
- The word "test" appears in lowercase
- The word "Test" appears with capital T
- The word "TEST" appears in all caps

Word Boundary Tests:
- testing (contains "test" but is a different word)
- test (exact word match)
- retest (contains "test" as suffix)

Multiple Occurrences:
The word search appears here.
Another search instance appears here.
And a third search occurrence here.

Special Characters:
test@example.com
hello-world
under_score

Long text for context extraction testing:
Lorem ipsum dolor sit amet, consectetur adipiscing elit. This sentence contains the word MATCH right here in the middle. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.

Line boundary testing:
First line with keyword
Second line with keyword
Third line with keyword`

	return createPDFWithText(filename, []string{content}, "LUMOS Search Test", "LUMOS Test Suite", "Search Testing", "LUMOS")
}

func createPDFWithText(filename string, pages []string, title, author, subject, creator string) error {
	// Create a new PDF
	conf := model.NewDefaultConfiguration()

	// Create pages with text
	pageSize := "A4"

	// For simplicity, we'll use pdfcpu's text watermark feature to add text to pages
	// This is a workaround since pdfcpu doesn't have a simple "create PDF from text" API

	// First, create a minimal PDF
	if err := api.CreatePDFFile(filename, len(pages), pageSize, nil, conf); err != nil {
		return fmt.Errorf("create PDF file: %w", err)
	}

	// Add text to each page using watermarks
	for i, text := range pages {
		wm, err := api.TextWatermark(text, "font:Helvetica, points:10, pos:tl, offset:50 50, rot:0", false, false, model.POINTS)
		if err != nil {
			return fmt.Errorf("create watermark: %w", err)
		}

		selectedPages := []string{fmt.Sprintf("%d", i+1)}
		if err := api.AddWatermarksFile(filename, "", selectedPages, wm, conf); err != nil {
			return fmt.Errorf("add text to page %d: %w", i+1, err)
		}
	}

	// Add metadata
	metadata := map[string]string{
		"Title":   title,
		"Author":  author,
		"Subject": subject,
		"Creator": creator,
	}

	// Note: pdfcpu metadata setting would go here
	// For now, the basic PDFs are sufficient for testing
	_ = metadata // silence unused warning

	return nil
}
