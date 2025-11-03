package main

import (
	"fmt"
	"os"
	"time"

	"github.com/luxor/lumos/pkg/pdf"
	"github.com/luxor/lumos/pkg/ui"
)

// Test with large PDFs to ensure performance is acceptable
func main() {
	fmt.Println("LUMOS Large PDF Test")
	fmt.Println("====================")
	fmt.Println()

	testFiles := []struct {
		path string
		name string
	}{
		{
			path: "../../PROJECTS/paper2agent/papers/Category Theory in Deep Learning_ A New Lens for Abstraction, Composition, and the Nature of Intelligence _ by Sethu Iyer _ Medium.pdf",
			name: "Category Theory in Deep Learning (11 MB)",
		},
		{
			path: "../../PROJECTS/paper2agent/papers/On Meta-Prompting.pdf",
			name: "On Meta-Prompting (2.9 MB)",
		},
		{
			path: "../../PROJECTS/paper2agent/papers/An Empirical Categorization of Prompting Techniques for Large Language Models_ A Practitioner's Guide.pdf",
			name: "Empirical Categorization (2.9 MB)",
		},
	}

	for i, test := range testFiles {
		fmt.Printf("[%d/%d] Testing: %s\n", i+1, len(testFiles), test.name)

		// Check if file exists
		fileInfo, err := os.Stat(test.path)
		if os.IsNotExist(err) {
			fmt.Printf("  ⚠️  SKIP - File not found\n\n")
			continue
		}

		sizeMB := float64(fileInfo.Size()) / (1024 * 1024)
		fmt.Printf("  📄 Size: %.2f MB\n", sizeMB)

		// Time the document loading
		startLoad := time.Now()
		doc, err := pdf.NewDocument(test.path, 5)
		loadTime := time.Since(startLoad)

		if err != nil {
			fmt.Printf("  ❌ FAIL - Error loading: %v\n\n", err)
			continue
		}

		pageCount := doc.GetPageCount()
		fmt.Printf("  📖 Pages: %d\n", pageCount)
		fmt.Printf("  ⏱️  Load time: %v\n", loadTime)

		// Test creating UI model
		startModel := time.Now()
		model := ui.NewModel(doc)
		modelTime := time.Since(startModel)

		if model == nil {
			fmt.Printf("  ❌ FAIL - Could not create model\n\n")
			continue
		}

		fmt.Printf("  ⏱️  Model creation: %v\n", modelTime)

		// Test loading first page
		startPage := time.Now()
		page, err := doc.GetPage(1)
		pageTime := time.Since(startPage)

		if err != nil {
			fmt.Printf("  ❌ FAIL - Could not load page 1: %v\n\n", err)
			continue
		}

		fmt.Printf("  📄 Page 1 text: %d chars\n", len(page.Text))
		fmt.Printf("  ⏱️  Page load time: %v\n", pageTime)

		// Test loading middle page
		middlePage := pageCount / 2
		startMiddle := time.Now()
		page, err = doc.GetPage(middlePage)
		middleTime := time.Since(startMiddle)

		if err != nil {
			fmt.Printf("  ⚠️  Could not load page %d: %v\n", middlePage, err)
		} else {
			fmt.Printf("  📄 Page %d text: %d chars\n", middlePage, len(page.Text))
			fmt.Printf("  ⏱️  Middle page load: %v\n", middleTime)
		}

		// Test view rendering
		startRender := time.Now()
		view := model.View()
		renderTime := time.Since(startRender)

		if view == "" {
			fmt.Printf("  ❌ FAIL - View rendering returned empty\n\n")
			continue
		}

		fmt.Printf("  🎨 View size: %d chars\n", len(view))
		fmt.Printf("  ⏱️  Render time: %v\n", renderTime)

		// Performance evaluation
		totalTime := loadTime + modelTime + pageTime + renderTime
		fmt.Printf("\n  📊 Total time: %v\n", totalTime)

		if totalTime < 1*time.Second {
			fmt.Println("  ✅ EXCELLENT - Very fast (<1s)")
		} else if totalTime < 3*time.Second {
			fmt.Println("  ✅ GOOD - Fast (<3s)")
		} else if totalTime < 5*time.Second {
			fmt.Println("  ⚠️  OK - Acceptable (<5s)")
		} else {
			fmt.Println("  ⚠️  SLOW - May need optimization (>5s)")
		}

		fmt.Println()
	}

	fmt.Println("====================")
	fmt.Println("Large PDF test completed!")
	fmt.Println()
	fmt.Println("Key Findings:")
	fmt.Println("  • LRU cache (5 pages) helps with navigation")
	fmt.Println("  • First page load time is critical for UX")
	fmt.Println("  • View rendering should be <100ms")
	fmt.Println("  • Large PDFs (>10MB) may benefit from lazy loading")
}
